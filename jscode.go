package flow

// func NewCode(code string, f ...Option) Code {
// 	var set Setting
// 	for _, f := range f {
// 		f(&set)
// 	}
// 	var sprm sync.Pool
// 	var appm atomic.Pointer[goja.Program]
// 	return func(scp Scope, oo []Node) (err error) {
// 		var rm, ok = sprm.Get().(*goja.Runtime)
// 		defer sprm.Put(&rm)
// 		if !ok {
// 			var scp = scp.WithContext(
// 				context.WithoutCancel(scp.Context()),
// 			)

// 			rm = goja.New()
// 			rm.SetFieldNameMapper(goja.UncapFieldNameMapper())
// 			if err = set.init(scp, rm); err != nil {
// 				return fmt.Errorf("jscode: %w", err)
// 			}
// 			var pm = appm.Load()
// 			if pm == nil {
// 				if pm, err = buildProgram(scp, code); err != nil {
// 					return fmt.Errorf("jscode: %w", err)
// 				}
// 				appm.Store(pm)
// 			}
// 			if _, err = rm.RunProgram(pm); err != nil {
// 				return fmt.Errorf("jscode: %w", err)
// 			}
// 		}
// 		// rm.Set("__flow_exec", rm.ToValue(func(call goja.FunctionCall) goja.Value {
// 		// 	var err error
// 		// 	var name string
// 		// 	if err = rm.ExportTo(call.Argument(0), &name); err != nil {
// 		// 		panic(rm.NewGoError(err))
// 		// 	}
// 		// 	var node []Node
// 		// 	if err = Convert(rm, call.Argument(1), &node); err != nil {
// 		// 		panic(rm.NewGoError(err))
// 		// 	}
// 		// 	var next Next = func(t []Node) error {

// 		// 	}

// 		// 	var pipe = scp.Import(name)
// 		// }))

// 		if err = set.call(scp, rm); err != nil {
// 			return fmt.Errorf("jscode: %w", err)
// 		}

// 		if err = set.quit(scp, rm); err != nil {
// 			return fmt.Errorf("jscode: %w", err)
// 		}
// 		return nil
// 	}
// }
