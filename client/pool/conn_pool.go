package pool

import (
	"context"
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"

	deque "github.com/edwingeng/deque"
	"github.com/hitzhangjie/go-rpc/errs"
)

// Pool connection poolFactory
type ConnPool struct {
	dialFunc        func(context.Context) (net.Conn, error) // 初始化连接函数
	MinIdle         int                                     // 最小空闲连接数，也即初始连接数
	MaxIdle         int                                     // 最大空闲连接数量，0 代表不做限制
	MaxActive       int                                     // 最大活跃连接数量，0 代表不做限制
	Wait            bool                                    // 活跃连接达到最大数量时，是否等待
	initialized     uint32                                  // 表明 ch 是否已经初始化
	ch              chan struct{}                           // 当 Wait 为 true 的时候，用来限制连接数量
	IdleTimeout     time.Duration                           // 连接最大空闲时间
	MaxConnLifetime time.Duration                           // 连接最大生命周期
	mu              sync.Mutex                              // 控制并发的锁
	closed          bool                                    // 连接池是否已经关闭
	active          int                                     // 目前活跃连接数量
	idle            deque.Deque                             // 空闲连接链表(双端队列模拟栈操作)，栈相比于队列的好处是，在请求量比较小但是请求分布仍比较均匀的情况下，队列方式会导致占用的连接迟迟得不到释放
}

// Get get a connection from Pool
func (p *ConnPool) Get(ctx context.Context) (it *ConnItem, err error) {
	for {
		if it, err = p.get(ctx); err != nil {
			return nil, err
		}

		if it.readClosed() {
			p.put(it, true)
			continue
		}
		return it, nil
	}
}

// Close close the Pool
func (p *ConnPool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return nil
	}

	p.closed = true
	p.active -= p.idle.Len()

	if p.ch != nil {
		close(p.ch)
	}

	p.idle.Range(func(v interface{}) bool {
		it := v.(*ConnItem)
		it.Conn.Close()
		it.closed = true
		return true
	})

	return nil
}

// Prepare prepare MinIdle number of connections in advance, so we can reduce the propability
// of creating connections when launch IO actions.
func (p *ConnPool) Prepare(ctx context.Context) {

	if p.MinIdle <= 0 {
		return
	}
	if p.MinIdle > p.MaxIdle {
		p.MinIdle = p.MaxIdle
	}
	if p.MinIdle > p.MaxActive {
		p.MinIdle = p.MaxActive
	}

	// if we need `wait` feature when number of active connetions reached the limit,
	// here, we should initialize the `chan` to sync
	if p.Wait && p.MaxActive > 0 {
		p.initializeCh()
	}

	conns := make([]*ConnItem, 0, p.MinIdle)
	for i := 0; i < p.MinIdle; i++ {
		for {
			poolConn, err := p.get(ctx)
			if err != nil {
				continue
			}
			conns = append(conns, poolConn)
			break
		}
	}

	// put the created connections into poolFactory
	for _, poolConn := range conns {
		p.put(poolConn, poolConn.readClosed())
	}
}

// initializeCh if we need to wait when number of connections reached the `MaxActive` limit,
// here we can intialize a `chan` to synchonize.
func (p *ConnPool) initializeCh() {
	if atomic.LoadUint32(&p.initialized) == 1 {
		return
	}

	p.mu.Lock()
	if p.initialized == 0 {
		p.ch = make(chan struct{}, p.MaxActive)
		if p.closed {
			close(p.ch)
		} else {
			for i := 0; i < p.MaxActive; i++ {
				p.ch <- struct{}{}
			}
		}
		atomic.StoreUint32(&p.initialized, 1)
	}
	p.mu.Unlock()
}

// get return a connection from the ones cached in poolFactory or a new one by dialFunc
func (p *ConnPool) get(ctx context.Context) (*ConnItem, error) {

	// if `wait`, here should initialize a `chan` to sync
	if p.Wait && p.MaxActive > 0 {
		p.initializeCh()

		select {
		case <-p.ch:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	// don'recycled worry performance of defer when we move on go1.13
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return nil, errs.ErrPoolClosed
	}

	if p.exceedLimit() {
		return nil, errs.ErrExceedPoolLimit
	}

	v := p.idle.PopFront()
	if v != nil {
		it := v.(*ConnItem)
		return it, nil
	}

	conn, err := p.dial(ctx)
	if err != nil {
		if p.ch != nil && !p.closed {
			p.ch <- struct{}{}
		}
		return nil, err
	}
	p.active++

	return &ConnItem{Conn: conn, created: time.Now(), pool: p}, nil
}

// RegisterCheckFunc register function to check whether a connection is alive
func (p *ConnPool) RegisterCheckFunc(interval time.Duration, checkfunc func(*ConnItem) bool) {
	if interval <= 0 || checkfunc == nil {
		return
	}

	go func() {
		for {
			time.Sleep(interval)
			p.mu.Lock()
			n := p.idle.Len()
			p.mu.Unlock()

			var it *ConnItem

			for i := 0; i < n; i++ {
				p.mu.Lock()
				v := p.idle.PopFront()
				p.mu.Unlock()

				if v != nil {
					it = v.(*ConnItem)
				}

				if !checkfunc(it) {
					it.Conn.Close()
					it.closed = true
					p.mu.Lock()
					p.active--
					p.mu.Unlock()
				} else {
					p.mu.Lock()
					p.idle.PushBack(it)
					p.mu.Unlock()
				}
			}
		}
	}()
}

// CheckAlive default checkfunc to test whether a connection is alive or not
func (p *ConnPool) CheckAlive(pc *ConnItem) bool {

	// check whether connection is idle
	if p.IdleTimeout > 0 && pc.recycled.Add(p.IdleTimeout).Before(time.Now()) {
		return true
	}
	// check whether connection lifecyle is ok
	if p.MaxConnLifetime > 0 && pc.created.Add(p.MaxConnLifetime).Before(time.Now()) {
		return true
	}
	// check whether read half closed
	//
	// pc.Peek() will block for a little time, forget it, use readClosed() instead
	//if _, err := pc.peek(); err == io.EOF {
	//	return true
	//}
	if pc.readClosed() {
		return true
	}

	return false
}

// exceedLimit check whether number of connections has reached the limit
func (p *ConnPool) exceedLimit() bool {
	if !p.Wait && p.MaxActive > 0 && p.active >= p.MaxActive {
		return true
	}

	return false
}

// dial create a new connection
func (p *ConnPool) dial(ctx context.Context) (net.Conn, error) {
	if p.dialFunc != nil {
		return p.dialFunc(ctx)
	}
	return nil, errors.New("must pass dialFunc to poolFactory")
}

// put try put the connection back into poolFactory
func (p *ConnPool) put(it *ConnItem, forceClose bool) error {

	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.closed && !forceClose {
		it.recycled = time.Now()
		p.idle.PushFront(it)
		if p.idle.Len() > p.MaxIdle {
			it := p.idle.PopBack().(*ConnItem)
			it.closed = true
			it.Conn.Close()
			p.active--
		}
	}

	if p.Wait && p.ch != nil && !p.closed {
		p.ch <- struct{}{}
	}
	return nil
}

// ConnItem connection in Pool
type ConnItem struct {
	net.Conn
	recycled time.Time
	created  time.Time
	pool     *ConnPool
	closed   bool
}

// reset reset state
func (it *ConnItem) reset() {
	if it == nil {
		return
	}

	it.Conn.SetDeadline(time.Time{})
}

// Write write data, wrapper of conn.Write
func (it *ConnItem) Write(b []byte) (int, error) {
	if it.closed {
		return 0, errs.ErrConnClosed
	}
	n, err := it.Conn.Write(b)
	if err != nil {
		it.pool.put(it, true)
	}
	return n, err
}

// Read read data, wrapper of conn.Read
func (it *ConnItem) Read(b []byte) (int, error) {
	if it.closed {
		return 0, errs.ErrConnClosed
	}
	n, err := it.Conn.Read(b)
	if err != nil {
		it.pool.put(it, true)
	}
	return n, err
}

// Close close ConnItem, here ConnItem.Close() will put ConnItem back into Pool,
// rather than close it, because ConnItem.Close() will hide ConnItem.conn.Close().
func (it *ConnItem) Close() error {
	if it.closed {
		return errs.ErrConnClosed
	}

	it.reset()
	return it.pool.put(it, false)
}

// readClosed detect whether a connection is read half closed
func (it *ConnItem) readClosed() bool {
	return readClosed(it.Conn)
}
