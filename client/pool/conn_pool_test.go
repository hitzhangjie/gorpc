package pool

import (
	"context"
	"errors"
	"github.com/edwingeng/deque"
	"io/ioutil"
	"log"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	network     = "tcp4"
	listenAddr1 = "127.0.0.1:8888"
	listenAddr2 = "127.0.0.1:9999"
	ch          = make(chan struct{})
	buffer      = []byte("hello world")
)

func init() {
	go buildSimpleTCPServer(ch)
	<-ch
	go buildSimpleTCPServerClose(ch)
	<-ch
}

func buildSimpleTCPServer(ch chan struct{}) {
	l, err := net.Listen(network, listenAddr1)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	ch <- struct{}{}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			buffer := make([]byte, 256)
			conn.Read(buffer)
		}()
	}
}

func buildSimpleTCPServerClose(ch chan struct{}) {
	l, err := net.Listen(network, listenAddr2)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	ch <- struct{}{}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			conn.Close()
		}()
	}
}

func newConnPool(maxIdle, maxActive int) *ConnPool {
	pool := &ConnPool{
		Dial: func(context.Context) (net.Conn, error) {
			return net.Dial(network, listenAddr1)
		},
		MaxIdle:   maxIdle,
		MaxActive: maxActive,
		idle:      deque.NewDeque(),
	}
	pool.RegisterCheckFunc(time.Millisecond*50, pool.CheckAlive)

	return pool
}

func TestConnPoolGet(t *testing.T) {
	pool := newConnPool(2, 10)

	_, err := pool.Get(context.Background())
	assert.Nil(t, err)

	assert.Equal(t, pool.active, 1)
	assert.Equal(t, pool.idle.Len(), 0)
}

func TestConnPoolGetOnClose(t *testing.T) {
	pool := newConnPool(2, 10)
	assert.Nil(t, pool.Close())

	ctx := context.Background()
	_, err := pool.Get(ctx)
	assert.Equal(t, err, errPoolClosed)
}

func TestConnPoolConcurrentGet(t *testing.T) {
	maxActive := 10
	pool := newConnPool(2, maxActive)
	defer pool.Close()

	var wg sync.WaitGroup
	ctx := context.Background()
	for i := 0; i < maxActive; i++ {
		wg.Add(1)
		func() {
			_, err := pool.Get(ctx)
			assert.Nil(t, err)
			wg.Done()
		}()
	}

	wg.Wait()

	assert.Equal(t, pool.active, maxActive)
	assert.Equal(t, pool.idle.Len(), 0)
}

func TestConnPoolOverMaxActive(t *testing.T) {
	maxActive := 10
	pool := newConnPool(2, maxActive)
	defer pool.Close()

	ctx := context.Background()
	for i := 0; i < maxActive; i++ {
		_, err := pool.Get(ctx)
		assert.Nil(t, err)
	}

	assert.Equal(t, pool.active, maxActive)

	_, err := pool.Get(ctx)
	assert.Equal(t, err, errExceedPoolLimit)
}

func TestConnPoolPut(t *testing.T) {
	pool := newConnPool(5, 10)
	defer pool.Close()

	pc, err := pool.Get(context.Background())
	assert.Nil(t, err)

	assert.Equal(t, pool.active, 1)
	assert.Equal(t, pool.idle.Len(), 0)

	assert.Nil(t, pc.Close())
	assert.Equal(t, pool.active, 1)
	assert.Equal(t, pool.idle.Len(), 1)
}

func TestConnPoolMaxIdel(t *testing.T) {
	pool := newConnPool(2, 10)
	defer pool.Close()

	ctx := context.Background()

	pc1, err := pool.Get(ctx)
	assert.Nil(t, err)
	pc2, err := pool.Get(ctx)
	assert.Nil(t, err)
	pc3, err := pool.Get(ctx)
	assert.Nil(t, err)

	assert.Equal(t, pool.active, 3)
	assert.Equal(t, pool.idle.Len(), 0)

	assert.Nil(t, pc1.Close())
	assert.Nil(t, pc2.Close())
	assert.Nil(t, pc3.Close())
	assert.Equal(t, pool.active, 2)
	assert.Equal(t, pool.idle.Len(), 2)

	assert.Equal(t, 2, pool.idle.Len())
}

func TestConnPoolIdleTimeout(t *testing.T) {
	pool := newConnPool(5, 10)
	pool.IdleTimeout = time.Millisecond * 200
	defer pool.Close()

	ctx := context.Background()
	connItems := []*ConnItem{}
	for i := 0; i < 10; i++ {
		pc, err := pool.Get(ctx)
		assert.Nil(t, err)
		connItems = append(connItems, pc)
	}

	for _, it := range connItems {
		assert.Nil(t, it.Close())
	}

	assert.Equal(t, pool.idle.Len(), 5)
	assert.Equal(t, pool.active, 5)
	time.Sleep(time.Millisecond * 500)
	assert.Equal(t, pool.idle.Len(), 0)
	assert.Equal(t, pool.active, 0)
	pc, err := pool.Get(ctx)
	assert.Nil(t, err)
	assert.Equal(t, pool.idle.Len(), 0)
	assert.Equal(t, pool.active, 1)

	pc.Close()
	assert.Equal(t, pool.idle.Len(), 1)
	assert.Equal(t, pool.active, 1)
}

func TestConnPoolReUseConn(t *testing.T) {
	pool := newConnPool(5, 10)
	defer pool.Close()

	ctx := context.Background()
	for i := 0; i < 2; i++ {
		pc1, err := pool.Get(ctx)
		assert.Nil(t, err)
		pc2, err := pool.Get(ctx)
		assert.Nil(t, err)
		pc1.Close()
		pc2.Close()
	}
	assert.Equal(t, pool.idle.Len(), 2)
	assert.Equal(t, pool.active, 2)
}

func TestConnPoolMaxLifeTime(t *testing.T) {
	pool := newConnPool(5, 10)
	defer pool.Close()
	pool.MaxConnLifetime = time.Millisecond * 400
	ctx := context.Background()
	its := []*ConnItem{}
	for i := 0; i < 10; i++ {
		pc, err := pool.Get(ctx)
		assert.Nil(t, err)
		its = append(its, pc)
	}

	for _, pc := range its {
		assert.Nil(t, pc.Close())
	}

	assert.Equal(t, pool.idle.Len(), 5)
	assert.Equal(t, pool.active, 5)
	time.Sleep(time.Second * 2)
	assert.Equal(t, pool.idle.Len(), 0)
	assert.Equal(t, pool.active, 0)

	pc, err := pool.Get(ctx)
	assert.Nil(t, err)
	assert.Equal(t, pool.idle.Len(), 0)
	assert.Equal(t, pool.active, 1)

	pc.Close()
	assert.Equal(t, pool.idle.Len(), 1)
	assert.Equal(t, pool.active, 1)
}

func connGetWithoutDeadline(p *ConnPool, n int) chan error {
	ctx := context.Background()
	errs := make(chan error, n)
	for i := 0; i < cap(errs); i++ {
		go func() {
			c, err := p.Get(ctx)
			if c != nil {
				c.Close()
			}
			errs <- err
		}()
	}

	return errs
}

func TestConnPoolWait(t *testing.T) {
	pool := newConnPool(5, 10)
	pool.Wait = true
	defer pool.Close()

	ctx := context.Background()
	its := []*ConnItem{}
	for i := 0; i < 10; i++ {
		it, err := pool.Get(ctx)
		assert.Nil(t, err)
		its = append(its, it)
	}

	errs := connGetWithoutDeadline(pool, 10)
	for _, pc := range its {
		assert.Nil(t, pc.Close())
	}

	timeout := time.After(2 * time.Second)
	for i := 0; i < cap(errs); i++ {
		select {
		case err := <-errs:
			assert.Nil(t, err)
		case <-timeout:
			t.Fatalf("timeout waiting for blocked goroutine %d", i)
		}
	}
}

func TestConnPoolWaitIdleTimeout(t *testing.T) {
	pool := newConnPool(5, 10)
	pool.Wait = true
	pool.IdleTimeout = time.Millisecond * 200
	defer pool.Close()

	ctx := context.Background()
	its := []*ConnItem{}
	for i := 0; i < 10; i++ {
		it, err := pool.Get(ctx)
		assert.Nil(t, err)
		its = append(its, it)
	}

	for _, pc := range its {
		assert.Nil(t, pc.Close())
	}

	time.Sleep(time.Millisecond * 500)
	timeout := time.After(1 * time.Second)
	errs := connGetWithoutDeadline(pool, 10)
	for i := 0; i < cap(errs); i++ {
		select {
		case err := <-errs:
			assert.Nil(t, err)
		case <-timeout:
			t.Fatalf("timeout waiting for blocked goroutine %d", i)
		}
	}
}

func TestConnPoolWaitMaxLifeTime(t *testing.T) {
	pool := newConnPool(5, 10)
	pool.Wait = true
	pool.MaxConnLifetime = time.Millisecond * 200
	defer pool.Close()

	ctx := context.Background()
	its := []*ConnItem{}
	for i := 0; i < 10; i++ {
		it, err := pool.Get(ctx)
		assert.Nil(t, err)
		its = append(its, it)
	}

	for _, pc := range its {
		assert.Nil(t, pc.Close())
	}

	time.Sleep(time.Millisecond * 500)
	timeout := time.After(1 * time.Second)
	errs := connGetWithoutDeadline(pool, 10)
	for i := 0; i < cap(errs); i++ {
		select {
		case err := <-errs:
			assert.Nil(t, err)
		case <-timeout:
			t.Fatalf("timeout waiting for blocked goroutine %d", i)
		}
	}
}

func TestConnPoolDialErr(t *testing.T) {
	pool := newConnPool(5, 10)
	defer pool.Close()

	ctx := context.Background()
	_, err := pool.Get(ctx)
	assert.Nil(t, err)
	assert.Equal(t, pool.active, 1)
	assert.Equal(t, pool.idle.Len(), 0)
	pool.Dial = func(context.Context) (net.Conn, error) {
		return nil, errors.New("dial error")
	}
	_, err = pool.Get(ctx)
	assert.NotNil(t, err)
	assert.Equal(t, pool.active, 1)
	assert.Equal(t, pool.idle.Len(), 0)
}

func TestConnPoolWaitDialErr(t *testing.T) {
	pool := newConnPool(5, 10)
	pool.Wait = true
	defer pool.Close()

	ctx := context.Background()
	for i := 0; i < 9; i++ {
		_, err := pool.Get(ctx)
		assert.Nil(t, err)
	}

	assert.Equal(t, pool.active, 9)
	assert.Equal(t, pool.idle.Len(), 0)
	pool.Dial = func(context.Context) (net.Conn, error) {
		return nil, errors.New("dial error")
	}
	_, err := pool.Get(ctx)
	assert.NotNil(t, err)
	pool.Dial = func(context.Context) (net.Conn, error) {
		return net.Dial(network, listenAddr1)
	}
	timeout := time.After(1 * time.Second)
	errs := connGetWithoutDeadline(pool, 1)
	for i := 0; i < cap(errs); i++ {
		select {
		case err := <-errs:
			assert.Nil(t, err)
		case <-timeout:
			t.Fatalf("timeout waiting for blocked goroutine %d", i)
		}
	}
	assert.Equal(t, pool.active, 10)
	assert.Equal(t, pool.idle.Len(), 1)
}

func TestConnPoolConnReadClosed(t *testing.T) {
	conn1, err := net.Dial(network, listenAddr2)
	assert.Nil(t, err)
	b, err := ioutil.ReadAll(conn1)
	assert.Nil(t, err)
	assert.Equal(t, len(b), 0)
	assert.True(t, readClosed(conn1))

	conn2, err := net.Dial(network, listenAddr1)
	assert.Nil(t, err)
	assert.False(t, readClosed(conn2))
}

func TestConnPoolConnReadClosedConnBroken(t *testing.T) {
	conn1, err := net.Dial(network, listenAddr2)
	assert.Nil(t, err)
	conn1.Close()
	assert.True(t, readClosed(conn1))
}

func TestConnPoolConnRead(t *testing.T) {
	pool := newConnPool(2, 10)
	ctx := context.Background()
	pc, err := pool.Get(ctx)
	assert.Nil(t, err)
	pc.Read(nil)
	pc.closed = true
	_, err = pc.Read(nil)
	assert.Equal(t, err, errConnClosed)
}

func TestConnPoolConnWrite(t *testing.T) {
	pool := newConnPool(2, 10)
	ctx := context.Background()
	pc, err := pool.Get(ctx)
	assert.Nil(t, err)
	pc.Write(nil)
	pc.closed = true
	_, err = pc.Write(nil)
	assert.Equal(t, err, errConnClosed)
}

func TestConnPoolPrepare(t *testing.T) {
	pool := &ConnPool{
		Dial: func(context.Context) (net.Conn, error) {
			return net.Dial(network, listenAddr1)
		},
		MinIdle:   2,
		MaxIdle:   5,
		MaxActive: 10,
		idle:      deque.NewDeque(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()
	pool.Prepare(ctx)

	assert.Equal(t, pool.idle.Len(), 2)
	assert.Equal(t, pool.active, 2)
}

func TestConnPoolConnectionGet(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	pool := NewConnectionPool()
	conn, err := pool.Get(ctx, network, listenAddr1)
	assert.Nil(t, err)
	assert.Nil(t, conn.Close())
}

func BenchmarkConnPoolGet(b *testing.B) {
	ctx := context.Background()
	p := newConnPool(1, 0)
	defer p.Close()
	c, err := p.Get(ctx)
	if err != nil {
		panic(err)
	}
	if err := c.Close(); err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		c, err = p.Get(ctx)
		if err != nil {
			panic(err)
		}
		if err := c.Close(); err != nil {
			panic(err)
		}
	}
}

func BenchmarkGetConnPoolParallel(b *testing.B) {
	b.StopTimer()
	ctx := context.Background()
	pool := newConnPool(10, 0)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pc, err := pool.Get(ctx)
			if err != nil {
				panic(err)
			}
			if err := pc.Close(); err != nil {
				panic(err)
			}
		}
	})
}
