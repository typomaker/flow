package flow

import (
	"context"
	"io/fs"
	"log/slog"
)

type Context struct {
	context.Context
	fs.FS
	logger *slog.Logger
}

func (it Context) Logger() *slog.Logger {
	if it.logger != nil {
		return it.logger
	}
	return slog.Default()
}
func Compose(stmts ...Statement) Statement {
	return func(ctx Context, target []Node, next Next) (err error) {
		var i = 0
		var step Next
		step = func(target []Node) error {
			var x = i
			if x < len(stmts) {
				i++
				return stmts[x](ctx, target, step)
			}
			return next(target)
		}
		return step(target)
	}
}
func If(when, then Statement) Statement {
	return func(ctx Context, target []Node, next Next) (err error) {
		var called bool
		var inter Next = func(target []Node) error {
			called = true
			return then(ctx, target, next)
		}
		if err = when(ctx, target, inter); err != nil {
			return err
		}
		if called {
			return nil
		}
		return next(target)
	}
}
func Not(stmt Statement) Statement {
	return func(ctx Context, target []Node, next Next) (err error) {
		var called bool
		var inter Next = func(target []Node) error {
			called = true
			return nil
		}
		if err = stmt(ctx, target, inter); err != nil {
			return err
		}
		if called {
			return nil
		}
		return next(target)
	}
}

type Statement func(ctx Context, target []Node, next Next) (err error)
type Next func(target []Node) error

var noop Next = func(target []Node) error { return nil }
