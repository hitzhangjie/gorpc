package router

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/hitzhangjie/gorpc/errors"
)

type HandleWrapper func(ctx context.Context, req interface{}) (rsp interface{}, err error)

type Router struct {
	mapping map[string]HandleWrapper
	mux     *sync.RWMutex
}

func NewRouter() *Router {
	return &Router{
		mapping: make(map[string]HandleWrapper),
		mux:     &sync.RWMutex{},
	}
}

func (r *Router) RegisterService(serviceDesc *ServiceDesc, serviceImpl interface{}) error {

	// check whether serviceImpl implements serviceDesc.ServiceType
	ht := reflect.TypeOf(serviceDesc.ServiceType).Elem()
	hi := reflect.TypeOf(serviceImpl)
	if !hi.Implements(ht) {
		return fmt.Errorf("%s not implements interface %s", hi.String(), ht.String())
	}

	// check pass, now register <serviceDesc, serviceImpl>
	r.mux.Lock()
	for _, m := range serviceDesc.Method {
		rpc := serviceDesc.Name + "/" + m.Name
		// handle wrapper
		f := func(ctx context.Context, req interface{}) (rsp interface{}, err error) {
			// handle func
			return m.Method(serviceImpl, ctx, req)
		}
		r.mapping[rpc] = f
	}
	r.mux.Unlock()

	return nil
}

func (r *Router) Forward(rpcName string, handlefunc HandleWrapper) {
	r.mapping[rpcName] = handlefunc
}

func (r *Router) Route(rpc string) (HandleWrapper, error) {

	h, ok := r.mapping[rpc]
	if !ok {
		return nil, errors.ErrRouteNotFound
	}

	return h, nil
}
