package flow

import (
	"context"
	"io/fs"
	"log/slog"
	"slices"
	"time"

	"github.com/laher/mergefs"
	slogmulti "github.com/samber/slog-multi"
)

func New(o ...Setup) (f Flow) {
	var s Setting
	for i := range o {
		o[i].setup(&s)
	}
	f.fs = s.FS
	f.logger = s.Logger
	f.handler = s.Handler
	f.extension = slices.Clip(s.Extension)
	return f
}

type Flow struct {
	fs        fs.FS
	logger    *slog.Logger
	handler   Handler
	extension []LogAttrer
}

func (it Flow) setup(s *Setting) {
	s.FS = it.fs
	s.Logger = it.logger
	s.Handler = it.handler
}

type Setup interface {
	setup(s *Setting)
}
type Setting struct {
	FS        fs.FS
	Logger    *slog.Logger
	Handler   Handler
	Extension []LogAttrer
}

func FS(f fs.FS) Setup {
	if f == nil {
		return optionFunc(func(s *Setting) {})
	}
	return optionFunc(func(s *Setting) {
		if s.FS != nil {
			s.FS = mergefs.Merge(f, s.FS)
		} else {
			s.FS = f
		}
	})
}
func (it Flow) FS() fs.FS {
	if it.fs != nil {
		return it.fs
	}
	return noopFS{}
}
func Logger(l *slog.Logger) Setup {
	if l == nil {
		return optionFunc(func(s *Setting) {})
	}
	return optionFunc(func(s *Setting) {
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
	if it.logger != nil {
		return it.logger
	}
	return slog.Default()
}
func (it Flow) Handler() Handler {
	return it.handler
}
func Extension(l ...LogAttrer) Setup {
	if l == nil {
		return optionFunc(func(s *Setting) {})
	}
	return optionFunc(func(s *Setting) {
		s.Extension = append(s.Extension, l...)
	})
}
func (it Flow) Extension() []LogAttrer {
	return it.extension
}
func (it Flow) Run(ctx context.Context, target []Node, extension ...LogAttrer) (err error) {
	if it.handler == nil {
		return
	}
	if extension != nil {
		if it.extension != nil {
			extension = append(
				make([]LogAttrer, 0, len(it.extension)+len(extension)),
				extension...,
			)
		}
		it.extension = extension
	}
	if _, ok := ctx.Value(contextSettingKey{}).(Flow); !ok {
		ctx = ContextWith(ctx, it)
	}
	if err = it.handler(ctx, target, noopNext); err != nil {
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

func (it Handler) setup(s *Setting) {
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

type Notifier interface {
	LogAttrer
	Notify(ctx context.Context, c Case) error
}
type Modifier interface {
	LogAttrer
	Modify(ctx context.Context, m Node) error
}
type LogAttrer interface {
	LogAttr() slog.Attr
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

type optionFunc func(*Setting)

func (it optionFunc) setup(f *Setting) {
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
