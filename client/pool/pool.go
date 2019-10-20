package pool

import (
	"context"
	"github.com/edwingeng/deque"
	"net"
	"sync"
	"time"
)

type PoolFactory interface {
	Get(ctx context.Context, network string, address string) (net.Conn, error)
}


// NewConnPoolFactory create a connection poolFactory manager
func NewConnPoolFactory(opt ...Option) PoolFactory {

	opts := &Options{
		MaxIdle:     5,
		IdleTimeout: 60 * time.Second,
		DialTimeout: 200 * time.Millisecond,
	}

	for _, o := range opt {
		o(opts)
	}

	return &poolFactory{
		opts:      opts,
		connPools: new(sync.Map),
	}
}

// poolFactory poolFactory manager, it maintains many <address,Pool> pairs
type poolFactory struct {
	opts      *Options
	connPools *sync.Map
}

// Get return a connection from poolFactory manager
func (p *poolFactory) Get(ctx context.Context, network string, address string) (net.Conn, error) {

	var cancel context.CancelFunc

	_, ok := ctx.Deadline()
	if !ok {
		ctx, cancel = context.WithTimeout(ctx, p.opts.DialTimeout)
		defer cancel()
	}

	key := address + "/" + network

	if v, ok := p.connPools.Load(key); ok {
		return v.(*ConnPool).Get(ctx)
	}

	newPool := &ConnPool{
		dialFunc: func(ctx context.Context) (net.Conn, error) {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}

			timeout := p.opts.DialTimeout
			t, ok := ctx.Deadline()
			if ok {
				timeout = t.Sub(time.Now())
			}
			return net.DialTimeout(network, address, timeout)
		},
		MinIdle:         p.opts.MinIdle,
		MaxIdle:         p.opts.MaxIdle,
		MaxActive:       p.opts.MaxActive,
		Wait:            p.opts.Wait,
		MaxConnLifetime: p.opts.MaxConnLifetime,
		IdleTimeout:     p.opts.IdleTimeout,
		idle:            deque.NewDeque(),
	}

	// 规避初始化连接池map并发写的问题
	v, ok := p.connPools.LoadOrStore(key, newPool)
	if !ok {
		go newPool.Prepare(ctx)
		newPool.RegisterCheckFunc(time.Second*3, newPool.CheckAlive)
		return newPool.Get(ctx)
	}
	return v.(*ConnPool).Get(ctx)
}
