package goja

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/dop251/goja"
	jsoniter "github.com/json-iterator/go"
	"github.com/typomaker/flow"
	"github.com/typomaker/flow/build"
)

func New(path string) flow.Handler {
	var po sync.Pool
	var pm *goja.Program
	var mu sync.RWMutex

	return func(ctx context.Context, target []flow.Node, next flow.Next) (err error) {
		mu.RLock()
		var ready = pm != nil
		mu.RUnlock()

		if !ready {
			if mu.TryLock() {
				var b []byte
				if b, err = build.Goja(ctx, path); err != nil {
					return fmt.Errorf("goja: %w", err)
				}
				if pm, err = goja.Compile("", string(b), true); err != nil {
					return fmt.Errorf("goja: %w", err)
				}
			} else {
				mu.Lock()
			}
			mu.Unlock()
		}
		var rm *goja.Runtime
		if rm, _ = po.Get().(*goja.Runtime); rm == nil {
			rm = goja.New()
			rm.SetFieldNameMapper(goja.UncapFieldNameMapper())
			if err = importConsole(ctx, rm, path); err != nil {
				return fmt.Errorf("goja: %w", err)
			}
			if _, err = rm.RunProgram(pm); err != nil {
				return fmt.Errorf("goja: %w", err)
			}
		}
		defer po.Put(rm)

		var jsMain goja.Callable
		if err = exportMain(ctx, rm, &jsMain); err != nil {
			return fmt.Errorf("goja: %w", err)
		}
		var jsTarget goja.Value
		if err = convert(rm, target, &jsTarget); err != nil {
			return fmt.Errorf("goja: %w", err)
		}
		var jsThis = rm.NewObject()
		if err = importModify(ctx, rm, jsThis); err != nil {
			return fmt.Errorf("goja: %w", err)
		}
		if err = importNotify(ctx, rm, jsThis); err != nil {
			return fmt.Errorf("goja: %w", err)
		}
		var jsNext goja.Value
		if err = importNext(ctx, rm, next, &jsNext); err != nil {
			return fmt.Errorf("goja: %w", err)
		}
		if _, err = jsMain(jsThis, jsTarget, jsNext); err != nil {
			return fmt.Errorf("goja: %w", err)
		}
		if err = convert(rm, jsTarget, &target); err != nil {
			return fmt.Errorf("goja: %w", err)
		}
		return nil
	}
}
func exportMain(_ context.Context, rm *goja.Runtime, main *goja.Callable) (err error) {
	var entry goja.Value
	if entry = rm.Get("entry").(*goja.Object).Get("default"); entry == nil {
		if entry = rm.Get("entry").(*goja.Object).Get("main"); entry == nil {
			return fmt.Errorf("goja: undefined entry")
		}
	}
	var ok bool
	if *main, ok = goja.AssertFunction(entry); !ok {
		return fmt.Errorf("unexpected entry, must be callable")
	}
	return nil
}
func importConsole(ctx context.Context, rm *goja.Runtime, path string) (err error) {
	var flowctx = flow.Context(ctx)

	var logger = flowctx.Logger().With(
		slog.Group("runtime",
			slog.String("name", "goja"),
			slog.String("path", path),
		),
	)
	var o = rm.NewObject()
	if err = o.Set("log", newPrinter(ctx, rm, logger.InfoContext)); err != nil {
		return err
	}
	if err = o.Set("error", newPrinter(ctx, rm, logger.ErrorContext)); err != nil {
		return err
	}
	if err = o.Set("warn", newPrinter(ctx, rm, logger.WarnContext)); err != nil {
		return err
	}
	if err = o.Set("info", newPrinter(ctx, rm, logger.InfoContext)); err != nil {
		return err
	}
	if err = o.Set("debug", newPrinter(ctx, rm, logger.DebugContext)); err != nil {
		return err
	}
	if err = rm.Set("console", o); err != nil {
		return err
	}
	return nil
}
func importModify(ctx context.Context, rm *goja.Runtime, this *goja.Object) (err error) {
	const name = "modify"
	var flowctx = flow.Context(ctx)
	var modifiers []flow.Modifier

	for _, extension := range flowctx.Extension() {
		if modifier, ok := extension.(flow.Modifier); ok {
			modifiers = append(modifiers, modifier)
		}
	}
	if len(modifiers) == 0 {
		return this.Set(name, rm.ToValue(func(c goja.FunctionCall) goja.Value {
			return goja.Undefined()
		}))
	}
	return this.Set(name, rm.ToValue(func(c goja.FunctionCall) goja.Value {
		if len(c.Arguments) == 0 {
			return goja.Undefined()
		}
		var err error
		var jsFlowNode = c.Argument(0)
		var goFlowNode flow.Node
		if err = convert(rm, jsFlowNode, &goFlowNode); err != nil {
			err = fmt.Errorf("goja: %w", err)
			panic(rm.NewGoError(err))
		}
		for _, modifier := range modifiers {
			if err = modifier.Modify(ctx, goFlowNode); err != nil {
				err = fmt.Errorf("goja: %w", err)
				panic(rm.NewGoError(err))
			}
		}

		return goja.Undefined()
	}))
}
func importNotify(ctx context.Context, rm *goja.Runtime, this *goja.Object) (err error) {
	const name = "notify"
	var flowctx = flow.Context(ctx)
	var notifiers []flow.Notifier

	for _, extension := range flowctx.Extension() {
		if notifier, ok := extension.(flow.Notifier); ok {
			notifiers = append(notifiers, notifier)
		}
	}
	if len(notifiers) == 0 {
		return this.Set(name, rm.ToValue(func(c goja.FunctionCall) goja.Value {
			return goja.Undefined()
		}))
	}
	return this.Set(name, rm.ToValue(func(c goja.FunctionCall) goja.Value {
		if len(c.Arguments) == 0 {
			return goja.Undefined()
		}
		var err error
		var jsFlowCase = c.Argument(0)
		var goFlowCase flow.Case
		if err = convert(rm, jsFlowCase, &goFlowCase); err != nil {
			err = fmt.Errorf("goja: %w", err)
			panic(rm.NewGoError(err))
		}
		for _, notifier := range notifiers {
			if err = notifier.Notify(ctx, goFlowCase); err != nil {
				err = fmt.Errorf("goja: %w", err)
				panic(rm.NewGoError(err))
			}
		}

		return goja.Undefined()
	}))
}
func importNext(_ context.Context, rm *goja.Runtime, next flow.Next, jsNext *goja.Value) (err error) {
	*jsNext = rm.ToValue(func(c goja.FunctionCall) goja.Value {
		if len(c.Arguments) == 0 {
			return goja.Undefined()
		}
		var err error
		var jsTarget = c.Arguments[0]
		var target []flow.Node
		if err = convert(rm, jsTarget, &target); err != nil {
			err = fmt.Errorf("goja: %w", err)
			panic(rm.NewGoError(err))
		}
		if err = next(target); err != nil {
			err = fmt.Errorf("goja: %w", err)
			panic(rm.NewGoError(err))
		}
		if err = convert(rm, target, &jsTarget); err != nil {
			err = fmt.Errorf("goja: %w", err)
			panic(rm.NewGoError(err))
		}
		return goja.Undefined()
	})
	return nil
}
func newPrinter(ctx context.Context, rm *goja.Runtime, printer func(context.Context, string, ...any)) goja.Value {
	var err error
	return rm.ToValue(func(call goja.FunctionCall) goja.Value {
		var offset = 0
		var message string
		var root = make([]any, 0, 1)
		var nest = make([]any, 0, len(call.Arguments))
		if call.Argument(offset).ExportType() == reflectString {
			message, _ = call.Argument(offset).Export().(string)
			offset++
		} else {
			message = "js print"
		}
		if arg := call.Argument(offset); arg.ExportType() == reflectObject {
			var jsObject, _ = arg.(*goja.Object)
			for _, key := range jsObject.Keys() {
				var jsValue = jsObject.Get(key)
				if key == "node" {
					var node flow.Node
					if err = convert(rm, jsValue, &node); err == nil {
						root = append(root, node.LogAttr())
						continue
					}
				}

				var v any
				if err = convert(rm, jsValue, &v); err != nil {
					nest = append(nest, slog.String(key+"Error", err.Error()))
				} else if s, err := jsoniter.MarshalToString(v); err != nil {
					nest = append(nest, slog.String(key+"Error", err.Error()))
				} else {
					nest = append(nest, slog.String(key, s))
				}
			}
			offset++
		}
		if l := len(call.Arguments) - offset; l > 0 {
			var args = make([]any, 0, l)
			for i := offset; i < len(call.Arguments); i++ {
				var arg = call.Argument(i)
				var val any
				if err = convert(rm, arg, &val); err != nil {
					val = err
				}
				args = append(args, val)
			}
			if s, err := jsoniter.MarshalToString(args); err != nil {
				nest = append(nest, slog.String("argsError", err.Error()))
			} else {
				nest = append(nest, slog.Any("args", s))
			}
		}

		var fields = append(root, slog.Group("js", nest...))
		printer(ctx, message, fields...)
		return nil
	})
}
