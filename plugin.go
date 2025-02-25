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

func (it Plugin) triggerInit(ctx context.Context, t Target) (err error) {
	if it.Init == nil {
		return nil
	}
	return it.Init(ctx, t)
}
func (it Plugin) triggerCall(ctx context.Context, t Target) (err error) {
	if it.Call == nil {
		return nil
	}
	return it.Call(ctx, t)
}
func (it Plugin) triggerQuit(ctx context.Context, t Target) (err error) {
	if it.Quit == nil {
		return nil
	}
	return it.Quit(ctx, t)
}

type Callback = func(context.Context, Target) error

type Target struct {
	flow *Flow
	pipe Pipe
	goja *goja.Runtime
	this *goja.Object
}

func (it Target) Flow() *Flow {
	return it.flow
}
func (it Target) Pipe() Pipe {
	return it.pipe
}
func (it Target) Goja() *goja.Runtime {
	return it.goja
}
func (it Target) This() *goja.Object {
	return it.this
}
