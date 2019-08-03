package router

import (
	"fmt"
	"github.com/hitzhangjie/go-rpc/codec"
	"reflect"
	"strings"
	"sync"
)

type Router struct {
	descMapping    map[string]*ServiceDesc
	serviceMapping map[string]interface{}
	mux            sync.RWMutex
}

func NewRouter() *Router {
	return &Router{
		descMapping:    map[string]*ServiceDesc{},
		serviceMapping: map[string]interface{}{},
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
	r.descMapping[serviceDesc.Name] = serviceDesc
	r.serviceMapping[serviceDesc.Name] = serviceImpl
	r.mux.Unlock()

	return nil
}

const (
	idxPackageName = iota
	idxServiceName
	idxMethodName
	idxUpperLimit
)

func (r *Router) Route(session codec.Session) (service interface{}, handle HandleFunc, err error) {

	// rpcName conforms to "$pkgName.$serviceName.$methodName"
	v := strings.Split(session.RPC(), ".")
	if len(v) != idxUpperLimit {
		err = routeNotFound
		return
	}
	p := v[idxPackageName]
	s := v[idxServiceName]
	m := v[idxMethodName]

	k := p + "." + s

	// find registered serviceDesc, method
	if sd, ok := r.descMapping[k]; !ok || sd == nil {
		err = routeNotFound
		return
	} else {
		if md, ok := sd.Method[m]; !ok || md == nil {
			err = routeNotFound
			return
		} else {
			handle = md.Method
		}
	}

	// find registered serviceImpl
	if s, ok := r.serviceMapping[k]; !ok || service == nil {
		err = routeNotFound
		return
	} else {
		service = s
	}

	return
}
