package gojaflow

import (
	"fmt"
	"sync"

	"github.com/dop251/goja"
	"github.com/typomaker/flow"
	"github.com/typomaker/flow/build"
)

type Goja struct {
	pool sync.Map // map[string]*goja.Runtime
}

func New(name string) flow.Statement {
	var po sync.Pool
	var pm *goja.Program
	var mu sync.RWMutex

	return func(ctx flow.Context, target []flow.Node, next flow.Next) (err error) {
		mu.RLock()
		var p = pm
		mu.RUnlock()
		if p == nil {
			if mu.TryLock() {
				var b []byte
				if b, err = build.Goja(ctx, name); err != nil {
					return fmt.Errorf("gojaflow: %w", err)
				}
				if p, err = goja.Compile(name, string(b), true); err != nil {
					return fmt.Errorf("gojaflow: %w", err)
				}
				pm = p
			} else {
				mu.Lock()
			}
			mu.Unlock()
		}
		var rm *goja.Runtime
		if rm = po.Get().(*goja.Runtime); rm == nil {
			rm = goja.New()
			rm.SetFieldNameMapper(goja.UncapFieldNameMapper())
			if _, err = rm.RunProgram(p); err != nil {
				return fmt.Errorf("gojaflow: %w", err)
			}
			po.Put(rm)
		}

		var entry goja.Value
		if entry = rm.Get("entry").(*goja.Object).Get("default"); entry == nil {
			if entry = rm.Get("entry").(*goja.Object).Get("main"); entry == nil {
				return fmt.Errorf("gojaflow: undefined entrypoint")
			}
		}
		var main goja.Callable
		var ok bool
		if main, ok = goja.AssertFunction(entry); !ok {
			return fmt.Errorf("unexpected entrypoint, must be callable")
		}
		var jsTarget goja.Value
		if err = convert(rm, target, &jsTarget); err != nil {
			return fmt.Errorf("gojaflow: %w", err)
		}
		var this = rm.NewObject()
		if _, err = main(this, jsTarget); err != nil {
			return fmt.Errorf("gojaflow: %w", err)
		}
		if err = convert(rm, jsTarget, &target); err != nil {
			return fmt.Errorf("gojaflow: %w", err)
		}
		return nil
	}
}
