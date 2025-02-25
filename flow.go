package flow

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path"
	"slices"
	"strings"
	"sync"

	"github.com/dop251/goja"
	"github.com/evanw/esbuild/pkg/api"
	jsoniter "github.com/json-iterator/go"
)

type Flow struct {
	stock []Pipe
	index map[string]int
	chain [][]Pipe

	plugin []Plugin
	logger *slog.Logger

	state struct {
		mu      sync.RWMutex
		runtime map[string]*sync.Pool
		program map[string]*goja.Program
	}
}
type Option func(*Flow)

func WithFlow(v *Flow) Option {
	return func(f *Flow) {
		f.stock = append(f.stock, v.stock...)
		f.plugin = append(f.plugin, v.plugin...)
		f.logger = v.logger
	}
}
func WithPipe(v ...Pipe) Option {
	return func(f *Flow) {
		f.stock = append(f.stock, v...)
	}
}
func WithLogger(v *slog.Logger) Option {
	return func(f *Flow) {
		f.logger = v
	}
}
func WithPlugin(v Plugin) Option {
	return func(f *Flow) {
		f.plugin = append(f.plugin, v)
	}
}

func New(o ...Option) *Flow {
	var it = Flow{}
	for _, v := range o {
		v(&it)
	}
	it.stock = slices.Clip(it.stock)
	it.plugin = slices.Clip(it.plugin)

	var index = make(map[string]int, len(it.stock))
	var stateRuntime = make(map[string]*sync.Pool, len(it.stock))
	var stateProgram = make(map[string]*goja.Program, len(it.stock))
	for i := 0; i < len(it.stock); i++ {
		var pipe = it.stock[i]
		var uuid = pipe.UUID.Get().String()
		if _, ok := index[uuid]; ok {
			copy(it.stock[i:], it.stock[i+1:])
			it.stock = it.stock[:len(it.stock)-1]
			i--
			continue
		}
		index[uuid] = i
		if pipe.Name.IsSome() {
			index[pipe.Name.Get().String()] = i
		}
		stateRuntime[uuid] = &sync.Pool{}
		stateProgram[uuid] = nil
	}
	var chain = make([][]Pipe, len(it.stock))
	for i, p := range it.stock {
		var ux map[UUID]struct{}
		defer reuseMapUUIDSrtuct(&ux)()

		var pq = make([]Pipe, 0, 8)
		pq = append(pq, p)
		ux[p.UUID.Get()] = struct{}{}

		for j := 0; j < len(pq); j++ {
			for _, du := range pq[j].Next.GetOrZero() {
				if _, ok := ux[du]; ok {
					continue
				}
				ux[du] = struct{}{}

				if di, ok := index[du.String()]; ok {
					pq = append(pq, it.stock[di])
				}
			}
		}
		chain[i] = slices.Clip(pq)
	}
	it.index = index
	it.chain = chain
	it.state.runtime = stateRuntime
	it.state.program = stateProgram
	return &it
}
func (it *Flow) Work(ctx context.Context, nn []Node) (err error) {
	if len(nn) == 0 {
		return nil
	}

	// pipe to node matching
	var chain []Pipe
	defer reuseSlicePipe(&chain)()

	var group [][]Node
	defer reuseSliceSliceNode(&group)()

	var head int
	var prev Pipe
	for i := range nn {
		var self Pipe
		for _, f := range it.stock {
			if !f.When.IsSome() {
				continue
			}
			if !nn[i].When(f.When.Get()) {
				continue
			}
			self = f
			break
		}
		switch {
		case self.UUID.IsZero():
			head++
			prev = Pipe{}
		case prev.UUID != self.UUID:
			head = i
			chain = append(chain, it.stock[it.index[self.UUID.Get().String()]])
			group = append(group, nn[head:i+1])
			prev = self
		default:
			group[len(chain)-1] = nn[head : i+1]
		}
	}
	var errs []error
	for i, pipe := range chain {
		var handler = it.handler(ctx, pipe)
		if err = handler(group[i]); err != nil {
			errs = append(errs, fmt.Errorf("pipe %q %w", pipe.String(), err))
		}
	}
	if len(errs) != 0 {
		err = errors.Join(errs...)
		return fmt.Errorf("flow: %w", err)
	}
	return nil
}
func (it *Flow) Import(name string) Pipe {
	if i, ok := it.index[name]; ok {
		return it.stock[i]
	}
	return Pipe{}
}
func (it *Flow) Logger() *slog.Logger {
	if it.logger != nil {
		return it.logger
	}
	return slog.Default()
}
func (it *Flow) handler(ctx context.Context, pipe Pipe) (main func([]Node) error) {
	return func(nodes []Node) error {
		var err error
		var uuid = pipe.UUID.String()
		var pool = it.state.runtime[uuid]
		var rm, _ = pool.Get().(*goja.Runtime)
		if rm == nil {
			rm, err = it.runtime(ctx, pipe)
		}
		if err != nil {
			return err
		}
		defer pool.Put(&rm)

		var done = make(chan struct{})
		defer close(done)

		rm.ClearInterrupt()
		go func() {
			select {
			case <-ctx.Done():
				rm.Interrupt(ctx.Err())
			case <-done:
				return
			}
		}()
		var this = rm.NewObject()

		if err = it.triggerCall(ctx, Api{it, pipe, rm, this}); err != nil {
			return err
		}
		var main goja.Callable
		if err = it.exportEntry(ctx, rm, &main); err != nil {
			return err
		}
		var target goja.Value
		if err = convert(rm, nodes, &target); err != nil {
			return err
		}
		if _, err = main(this, target); err != nil {
			return err
		}
		if err = convert(rm, target, &nodes); err != nil {
			return err
		}
		if err = it.triggerQuit(ctx, Api{it, pipe, rm, this}); err != nil {
			return err
		}
		return nil
	}
}
func (it *Flow) runtime(ctx context.Context, pipe Pipe) (rm *goja.Runtime, err error) {
	ctx = context.WithoutCancel(ctx)
	rm = goja.New()
	rm.SetFieldNameMapper(goja.UncapFieldNameMapper())
	if err = it.importConsole(ctx, rm, pipe); err != nil {
		return nil, err
	}

	if err = it.triggerInit(ctx, Api{it, pipe, rm, rm.GlobalObject()}); err != nil {
		return nil, err
	}
	var pm *goja.Program
	if pm, err = it.program(ctx, pipe); err != nil {
		return nil, err
	}
	if _, err = rm.RunProgram(pm); err != nil {
		return nil, err
	}
	return rm, nil

}
func (it *Flow) program(ctx context.Context, pipe Pipe) (pm *goja.Program, err error) {
	var uuid = pipe.UUID.String()
	it.state.mu.RLock()
	pm = it.state.program[uuid]
	it.state.mu.RUnlock()
	if pm != nil {
		return pm, nil
	}
	if pm, err = it.compile(ctx, pipe); err != nil {
		return nil, err
	}
	it.state.mu.Lock()
	it.state.program[uuid] = pm
	it.state.mu.Unlock()

	return pm, nil
}
func (it *Flow) compile(ctx context.Context, pipe Pipe) (pm *goja.Program, err error) {
	var chain = it.chain[it.index[pipe.UUID.Get().String()]]
	var boot = renderLoopJS(chain)
	var bundle = api.Build(api.BuildOptions{
		Stdin: &api.StdinOptions{
			Contents:   boot,
			Sourcefile: pipe.String(),
			ResolveDir: ".",
		},
		Format:            api.FormatIIFE,
		GlobalName:        "entry",
		Bundle:            true,
		Write:             false,
		MinifyWhitespace:  false,
		MinifyIdentifiers: false,
		MinifySyntax:      false,
		TreeShaking:       api.TreeShakingTrue,
		Sourcemap:         api.SourceMapInline,
		KeepNames:         false,
		Target:            api.ES2020,
		PreserveSymlinks:  true,
		Plugins: []api.Plugin{{
			Name: "flow",
			Setup: func(pb api.PluginBuild) {
				pb.OnResolve(
					api.OnResolveOptions{Filter: "^flow:\\./[a-f0-9_-]{36}$"},
					func(ora api.OnResolveArgs) (r api.OnResolveResult, err error) {
						r.Path = ora.Path[5:]
						r.Namespace = "flow"
						return r, nil
					},
				)
				pb.OnLoad(
					api.OnLoadOptions{Filter: ".*", Namespace: "flow"},
					func(ola api.OnLoadArgs) (r api.OnLoadResult, err error) {
						var name = path.Base(ola.Path)
						var next Pipe
						if next = it.Import(name); next.IsZero() {
							return r, fmt.Errorf("flow: import %q is invalid", name)
						}
						var code string
						if next.Code.IsSome() {
							code = next.Code.Get()
						} else {
							code = renderSkipJS()
						}

						r.ResolveDir = "."
						r.Loader = api.LoaderJS
						r.Contents = &code
						return r, nil
					},
				)
			},
		}},
	})
	if len(bundle.Errors) != 0 {
		var fmsg = api.FormatMessages(bundle.Errors, api.FormatMessagesOptions{Kind: api.ErrorMessage})
		return nil, errors.New(strings.Join(fmsg, ""))
	}
	if len(bundle.Warnings) != 0 {
		var fmsg = api.FormatMessages(bundle.Warnings, api.FormatMessagesOptions{Kind: api.WarningMessage})
		return nil, errors.New(strings.Join(fmsg, ""))
	}
	if len(bundle.OutputFiles) != 1 {
		panic(fmt.Sprintf("compile unexpected number of bundle files, got %d expected 1", len(bundle.OutputFiles)))
	}

	var src = string(bundle.OutputFiles[0].Contents)
	if pm, err = goja.Compile(pipe.String(), src, true); err != nil {
		return nil, err
	}
	return pm, nil
}
func (it *Flow) exportEntry(ctx context.Context, rm *goja.Runtime, main *goja.Callable) (err error) {
	var entry goja.Value
	if entry = rm.Get("entry").(*goja.Object).Get("default"); entry == nil {
		if entry = rm.Get("entry").(*goja.Object).Get("main"); entry == nil {
			return fmt.Errorf("undefined entrypoint")
		}
	}
	var ok bool
	if *main, ok = goja.AssertFunction(entry); !ok {
		return fmt.Errorf("unexpected entrypoint, must be callable")
	}
	return nil
}
func (it *Flow) importConsole(ctx context.Context, rm *goja.Runtime, pipe Pipe) (err error) {
	var logger = it.Logger().With(slog.Any("pipe", pipe))
	var o = rm.NewObject()
	if err = o.Set("log", it.newPrinter(ctx, rm, logger.InfoContext)); err != nil {
		return err
	}
	if err = o.Set("error", it.newPrinter(ctx, rm, logger.ErrorContext)); err != nil {
		return err
	}
	if err = o.Set("warn", it.newPrinter(ctx, rm, logger.WarnContext)); err != nil {
		return err
	}
	if err = o.Set("info", it.newPrinter(ctx, rm, logger.InfoContext)); err != nil {
		return err
	}
	if err = o.Set("debug", it.newPrinter(ctx, rm, logger.DebugContext)); err != nil {
		return err
	}
	if err = rm.Set("console", o); err != nil {
		return err
	}
	return nil
}
func (it *Flow) newPrinter(ctx context.Context, rm *goja.Runtime, printer func(context.Context, string, ...any)) goja.Value {
	var err error
	return rm.ToValue(func(call goja.FunctionCall) goja.Value {
		var start = 0
		var title string
		var field = make([]any, 0, len(call.Arguments))
		if call.Argument(start).ExportType() == reflectString {
			title, _ = call.Argument(start).Export().(string)
			start++
		} else {
			title = "flow message"
		}
		if arg := call.Argument(start); arg.ExportType() == reflectObject {
			var v any
			if err = convert(rm, arg, &v); err != nil {
				field = append(field, slog.String("logError", err.Error()))
			} else if o, ok := v.(map[string]any); !ok {
				field = append(field, slog.Any("logUnknow", arg.Export()))
			} else {
				for k, v := range o {
					field = append(field, slog.Any(k, v))
				}
				start++
			}
		}
		if l := len(call.Arguments) - start; l > 0 {
			var args = make([]any, 0, l)
			for i := start; i < len(call.Arguments); i++ {
				var arg = call.Argument(i)
				var val any
				if err = convert(rm, arg, &val); err != nil {
					val = err
				}
				args = append(args, val)
			}
			if s, err := jsoniter.MarshalToString(args); err != nil {
				field = append(field, slog.String("argsError", err.Error()))
			} else {
				field = append(field, slog.Any("args", s))
			}
		}
		printer(ctx, title, field...)
		return nil
	})
}
func (it *Flow) triggerInit(ctx context.Context, t Api) (err error) {
	for _, p := range it.plugin {
		if err = p.triggerInit(ctx, t); err != nil {
			return err
		}
	}
	return nil
}
func (it *Flow) triggerCall(ctx context.Context, t Api) (err error) {
	for _, p := range it.plugin {
		if err = p.triggerCall(ctx, t); err != nil {
			return err
		}
	}
	return nil
}
func (it *Flow) triggerQuit(ctx context.Context, t Api) (err error) {
	for _, p := range it.plugin {
		if err = p.triggerQuit(ctx, t); err != nil {
			return err
		}
	}
	return nil
}
