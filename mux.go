// Package xmux provides a type-safe, dependency injection friendly HTTP router
// with support for multiple web frameworks through adapters.
package xmux

import (
	"context"
	"sync"
)

type Api interface {
	Invoke(ctx context.Context, bind func(params any) error) (any, error)
	Params() any
	Response() any
	Function() any
}

type Controller interface {
	Handle(method string, path string, service any, api Api, options ...map[string]string)
}

type Router interface {
	Register(method string, path string, api Api, options ...map[string]string)
}

type Binder interface {
	Bind(handler Controller, bind func(service any) error) (err error)
}

type function[Params any, Response any] func(context.Context, *Params) (Response, error)

func (h function[Params, Response]) Invoke(ctx context.Context, unmarshal func(params any) error) (ret any, err error) {
	var params Params
	if err = unmarshal(&params); err != nil {
		return
	}
	return h(ctx, &params)
}

func (h function[Params, Response]) Function() any {
	return h
}

func (h function[Params, Response]) Params() any {
	var zero Params
	return zero
}

func (h function[Params, Response]) Response() any {
	var zero Response
	return zero
}

func Register[Params any, Response any](
	router Router,
	method string,
	path string,
	fn func(ctx context.Context, params *Params) (Response, error),
	options ...map[string]string,
) {
	router.Register(method, path, function[Params, Response](fn), options...)
}

func MergeOptions(options []map[string]string, desc bool) map[string]string {
	opt := make(map[string]string)
	for i := 0; i < len(options); i++ {
		o := options[i]
		if desc {
			o = options[len(options)-1-i]
		}
		for k, v := range o {
			opt[k] = v
		}
	}
	return opt
}

type RouterGroup[Service any] struct {
	register func(router Router, handler Service)
	options  []map[string]string
}

func (g RouterGroup[Service]) Bind(controller Controller, bind func(any) error) (err error) {
	var service Service
	if err = bind(&service); err != nil {
		return
	}
	g.register(registerFunc(func(method string, path string, api Api, options ...map[string]string) {
		controller.Handle(method, path, service, api, append(g.options, options...)...)
	}), service)
	return
}

type registerFunc func(method string, path string, api Api, options ...map[string]string)

func (fn registerFunc) Register(method string, path string, api Api, options ...map[string]string) {
	fn(method, path, api, options...)
}

func DefineGroup[Service any](fn func(router Router, handler Service), options ...map[string]string) Binder {
	return RouterGroup[Service]{
		options:  options,
		register: fn,
	}
}

type Groups interface {
	Binder
	Register(groups ...Binder) Groups
}

type groups struct {
	mu     sync.Mutex
	groups []Binder
}

func NewGroups(gs ...Binder) Groups {
	return &groups{groups: append(make([]Binder, 0, len(gs)), gs...)}
}

func (g *groups) Register(groups ...Binder) Groups {
	g.mu.Lock()
	g.groups = append(g.groups, groups...)
	g.mu.Unlock()
	return g
}

func (g *groups) Bind(controller Controller, bind func(service any) error) (err error) {
	g.mu.Lock()
	gs := append(make([]Binder, 0, len(g.groups)), g.groups...)
	g.mu.Unlock()
	for _, group := range gs {
		if err = group.Bind(controller, bind); err != nil {
			return
		}
	}
	return
}
