package goja

import (
	"context"
	"errors"

	"github.com/dop251/goja"
	"github.com/typomaker/flow"
)

var OnInit = flow.Plugin[Extension]{Name: "goja.oninit"}
var OnCall = flow.Plugin[Extension]{Name: "goja.oncall"}

type Extension func(ctx context.Context, rm *goja.Runtime, this *goja.Object) error

func (it Extension) With(v Extension) Extension {
	if it == nil {
		return v
	}
	return func(ctx context.Context, rm *goja.Runtime, this *goja.Object) error {
		return errors.Join(
			it(ctx, rm, this),
			v(ctx, rm, this),
		)
	}
}
