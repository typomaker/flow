package flow

import (
	"context"

	"github.com/dop251/goja"
)

type Plugin struct {
	Init Callback
	Call Callback
	Quit Callback
}

func (it Plugin) triggerInit(ctx context.Context, t Api) (err error) {
	if it.Init == nil {
		return nil
	}
	return it.Init(ctx, t)
}
func (it Plugin) triggerCall(ctx context.Context, t Api) (err error) {
	if it.Call == nil {
		return nil
	}
	return it.Call(ctx, t)
}
func (it Plugin) triggerQuit(ctx context.Context, t Api) (err error) {
	if it.Quit == nil {
		return nil
	}
	return it.Quit(ctx, t)
}

type Callback = func(context.Context, Api) error

type Api struct {
	flow    *Flow
	pipe    Pipe
	runtime *goja.Runtime
	this    *goja.Object
}

func (it Api) Flow() *Flow {
	return it.flow
}
func (it Api) Pipe() Pipe {
	return it.pipe
}
func (it Api) Runtime() *goja.Runtime {
	return it.runtime
}
func (it Api) This() *goja.Object {
	return it.this
}
