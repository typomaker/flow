package flow

import (
	"context"
	"errors"
	"io/fs"
	"log/slog"
	"time"

	"github.com/laher/mergefs"
	slogmulti "github.com/samber/slog-multi"
)

func New(o ...Option) Flow {
	var it Flow
	it.config = it.config.With(o...)
	return it
}
func (it Flow) With(o ...Option) Flow {
	it.config = it.config.With(o...)
	return it
}

type Flow struct {
	config Config
}

func (it Flow) apply(s *Config) {
	if v := it.config.FS; v != nil {
		FS(v).apply(s)
	}
	if v := it.config.Logger; v != nil {
		Logger(v).apply(s)
	}
	if v := it.config.Handler; v != nil {
		Handler(v).apply(s)
	}
	if v := it.config.Modifier; v != nil {
		Modifier(v).apply(s)
	}
	if v := it.config.Notifier; v != nil {
		Notifier(v).apply(s)
	}
}

type Option interface {
	apply(s *Config)
}

func FS(f fs.FS) Option {
	if f == nil {
		return optionFunc(func(s *Config) {})
	}
	return optionFunc(func(s *Config) {
		if s.FS != nil {
			s.FS = mergefs.Merge(f, s.FS)
		} else {
			s.FS = f
		}
	})
}
func (it Flow) FS() fs.FS {
	if it.config.FS != nil {
		return it.config.FS
	}
	return noopFS{}
}
func Logger(l *slog.Logger) Option {
	if l == nil {
		return optionFunc(func(s *Config) {})
	}
	return optionFunc(func(s *Config) {
		if s.Logger != nil {
			s.Logger = slog.New(
				slogmulti.Fanout(
					s.Logger.Handler(),
					l.Handler(),
				),
			)
		} else {
			s.Logger = l
		}
	})
}
func (it Flow) Logger() *slog.Logger {
	if it.config.Logger != nil {
		return it.config.Logger
	}
	return slog.Default()
}
func (it Flow) Handler() Handler {
	return it.config.Handler
}

type Modifier func(ctx context.Context, c Node) error

func (it Modifier) apply(c *Config) {
	if c.Modifier != nil {
		var m = c.Modifier
		c.Modifier = func(ctx context.Context, n Node) error {
			return errors.Join(
				m(ctx, n),
				it(ctx, n),
			)
		}
	} else {
		c.Modifier = it
	}
}
func (it Flow) Modifier() Modifier {
	if it.config.Modifier == nil {
		return func(ctx context.Context, n Node) error { return nil }
	}
	return it.config.Modifier
}

type Notifier func(ctx context.Context, c Case) error

func (it Notifier) apply(c *Config) {
	if c.Notifier != nil {
		var m = c.Notifier
		c.Notifier = func(ctx context.Context, n Case) error {
			return errors.Join(
				m(ctx, n),
				it(ctx, n),
			)
		}
	} else {
		c.Notifier = it
	}
}
func (it Flow) Notifier() Notifier {
	if it.config.Notifier == nil {
		return func(ctx context.Context, c Case) error { return nil }
	}
	return it.config.Notifier
}
func (it Flow) Run(ctx context.Context, target []Node) (err error) {
	if it.config.Handler == nil {
		return
	}
	if _, ok := ctx.Value(contextSettingKey{}).(Flow); !ok {
		ctx = ContextWith(ctx, it)
	}
	if err = it.config.Handler(ctx, target, noopNext); err != nil {
		return err
	}
	return nil
}
func Pipe(hs ...Handler) Handler {
	return func(ctx context.Context, target []Node, next Next) (err error) {
		var i = 0
		var step Next
		step = func(target []Node) error {
			var x = i
			if x < len(hs) {
				i++
				return hs[x](ctx, target, step)
			}
			return next(target)
		}
		return step(target)
	}
}
func And(hs ...Handler) Handler {
	if len(hs) == 0 {
		return func(ctx context.Context, target []Node, next Next) (err error) {
			return next(target)
		}
	}
	return func(ctx context.Context, target []Node, next Next) (err error) {
		var step Next
		step = func(target []Node) error {
			if len(hs) == 0 {
				return next(target)
			}
			var h = hs[0]
			hs = hs[1:]
			return h(ctx, target, step)
		}
		if err = step(target); err != nil {
			return err
		}
		if len(hs) != 0 {
			return next(target)
		}
		return
	}
}
func Or(hs ...Handler) Handler {
	if len(hs) == 0 {
		return func(ctx context.Context, target []Node, next Next) (err error) {
			return next(target)
		}
	}
	return func(ctx context.Context, target []Node, next Next) (err error) {
		var ok bool
		for _, h := range hs {
			var fn Next = func(target []Node) error {
				ok = true
				return next(target)
			}
			if err = h(ctx, target, fn); err != nil {
				return err
			}
			if ok {
				return nil
			}
		}
		return next(target)
	}
}
func Not(h Handler) Handler {
	return func(ctx context.Context, target []Node, next Next) (err error) {
		var ok bool
		var fn Next = func(target []Node) error {
			ok = true
			return nil
		}
		if err = h(ctx, target, fn); err != nil {
			return err
		}
		if ok {
			return nil
		}
		return next(target)
	}
}
func Always(h Handler) Handler {
	return func(ctx context.Context, target []Node, next Next) (err error) {
		var ok bool
		var fn Next = func(target []Node) error {
			ok = true
			return next(target)
		}
		if err = h(ctx, target, fn); err != nil {
			return err
		}
		if !ok {
			return next(target)
		}
		return nil
	}
}
func Never(h Handler) Handler {
	return func(ctx context.Context, target []Node, _ Next) (err error) {
		return h(ctx, target, noopNext)
	}
}

type Handler func(ctx context.Context, target []Node, next Next) (err error)

func (it Handler) apply(s *Config) {
	if s.Handler != nil {
		s.Handler = Pipe(s.Handler, it)
	} else {
		s.Handler = it
	}
}

type Next func(target []Node) error

func noopNext(n []Node) error {
	return nil
}

type contextSettingKey struct{}

func Context(ctx context.Context) Flow {
	if v, ok := ctx.Value(contextSettingKey{}).(Flow); ok {
		return v
	}
	return Flow{}
}
func ContextWith(ctx context.Context, s Flow) context.Context {
	return context.WithValue(ctx, contextSettingKey{}, s)
}

type optionFunc func(*Config)

func (it optionFunc) apply(f *Config) {
	it(f)
}

type noopFS struct{}

func (n noopFS) Open(name string) (fs.File, error) {
	return nil, fs.ErrNotExist
}

type noopFile struct{}

func (n noopFile) Stat() (fs.FileInfo, error) {
	return nil, fs.ErrNotExist
}
func (n noopFile) Read(b []byte) (int, error) {
	return 0, nil
}
func (n noopFile) Close() error {
	return nil
}

type noopFileInfo struct{}

func (n noopFileInfo) Name() string       { return "" }
func (n noopFileInfo) Size() int64        { return 0 }
func (n noopFileInfo) Mode() fs.FileMode  { return 0 }
func (n noopFileInfo) ModTime() time.Time { return time.Time{} }
func (n noopFileInfo) IsDir() bool        { return false }
func (n noopFileInfo) Sys() interface{}   { return nil }
