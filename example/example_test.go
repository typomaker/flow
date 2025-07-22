package example

import (
	"context"
	"fmt"
	"testing/fstest"

	"github.com/typomaker/flow"
	"github.com/typomaker/flow/goja"
	"github.com/typomaker/option"
)

func Example_script_matching() {
	var f = flow.New(
		// you should provide then FS containing your scripts
		flow.FS(fstest.MapFS{
			// each script implements the logic for specific node kind
			"cat.js": &fstest.MapFile{
				Data: []byte(`
					export default function main(nodes) {
						for (const node of nodes) {
							node.hook = {kind:"cat"}
						}
					}
				`),
			},
			"dog.js": &fstest.MapFile{
				Data: []byte(`
					export default function main(nodes) {
						for (const node of nodes) {
							node.hook = {kind:"dog"}
						}
					}
				`),
			},
		}),
		// you should also provide the conditions for matching scripts with nodes based on their attributes
		flow.Or(
			flow.And(
				// the predefined predicate for uuid matching, you could write your own
				flow.UUID.In(flow.MustUUID("10000000-0000-0000-0000-000000000000")),
				// the goja is js runtime
				goja.New("cat.js"),
			),
			flow.And(
				flow.UUID.In(flow.MustUUID("20000000-0000-0000-0000-000000000000")),
				goja.New("dog.js"),
			),
		),
	)
	// targets for processing, the only way to exchange data between runtimes
	var target = []flow.Node{
		{UUID: option.Some(flow.MustUUID("10000000-0000-0000-0000-000000000000"))},
		{UUID: option.Some(flow.MustUUID("20000000-0000-0000-0000-000000000000"))},
	}
	var ctx = context.Background()

	// run processing
	var err = f.Run(ctx, target)
	if err != nil {
		fmt.Println("handler error", err)
	}
	if target[0].Hook.Get()["kind"] == "cat" {
		fmt.Println("target[0] kind is cat")
	}
	if target[1].Hook.Get()["kind"] == "dog" {
		fmt.Println("target[1] kind is dog")
	}
	// Output:
	// target[0] kind is cat
	// target[1] kind is dog
}

func Example_script_pipelining() {
	var f = flow.New(
		// you should provide then FS containing your scripts
		flow.FS(fstest.MapFS{
			"first.js": &fstest.MapFile{
				Data: []byte(`
					// you should call next, which will allow the following scripts to execute
					export default function main(nodes, next) {
						for (const node of nodes) {
							node.uuid="f42c13e8-520e-41a3-afee-71a8622d2e3a"
							node.meta.first=true
							node.meta.caller.push("first");
						}
						next(nodes)
					}
				`),
			},
			"second.js": &fstest.MapFile{
				Data: []byte(`
					export default function main(nodes, next) {
						for (const node of nodes) {
							node.uuid="d2fd3f06-827a-40e4-a56c-69a1083d4335"
							node.meta.second=true
							node.meta.caller.push("second");
						}
						next(nodes)
					}
				`),
			},
		}),
		// the pipe run every script one by one
		flow.Pipe(
			goja.New("first.js"),
			goja.New("second.js"),
		),
	)
	// targets for processing, the only way to exchange data between runtimes
	var target = []flow.Node{
		{Meta: option.Some(flow.Meta{"caller": []any{}})},
	}
	var ctx = context.Background()

	// run processing
	var err = f.Run(ctx, target)
	if err != nil {
		fmt.Println("handler error", err)
	} else {
		fmt.Println(target[0].Meta.Get()["caller"])
	}
	// Output:
	// [first second]
}
func Example_typescript() {
	var f = flow.New(
		// you should provide then FS containing your scripts
		flow.FS(fstest.MapFS{
			"typescript.ts": &fstest.MapFile{
				Data: []byte(`
					type Next = (_: any[])=>void;
					export default function main(nodes: any[], next: Next) {
						for (const node of nodes) {
							node.uuid="f42c13e8-520e-41a3-afee-71a8622d2e3a"
						}
						next(nodes)
					}
				`),
			},
		}),
		goja.New("typescript.ts"),
	)
	// targets for processing, the only way to exchange data between runtimes
	var target = []flow.Node{
		{},
	}
	var ctx = context.Background()

	// run processing
	var err = f.Run(ctx, target)
	if err != nil {
		fmt.Println("handler error", err)
	} else {
		fmt.Println(target[0].UUID.Get())
	}
	// Output:
	// f42c13e8-520e-41a3-afee-71a8622d2e3a
}
func Example_import_unpkg() {
	var f = flow.New(
		// you should provide then FS containing your scripts
		flow.FS(fstest.MapFS{
			"script.js": &fstest.MapFile{
				Data: []byte(`
					import isSorted from "https://unpkg.com/is-sorted@1.0.5"
					export default function main(nodes, next) {
						for (const node of nodes) {
							node.meta.sorted=isSorted([1,2,3])
							node.meta.unsorted=isSorted([1,3,2])
						}
						next(nodes)
					}
				`),
			},
		}),
		goja.New("script.js"),
	)
	// targets for processing, the only way to exchange data between runtimes
	var target = []flow.Node{
		{Meta: option.Some(flow.Meta{})},
	}
	var ctx = context.Background()

	// run processing
	var err = f.Run(ctx, target)
	if err != nil {
		fmt.Println("handler error", err)
	} else {
		fmt.Println(target[0].Meta.Get())
	}
	// Output:
	// map[sorted:true unsorted:false]
}
