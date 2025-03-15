package flow

import (
	"context"
	"fmt"

	"github.com/dop251/goja"
	"github.com/typomaker/option"
)

// ExamplePipe_When
// How to process a node with a specific pipe
func ExamplePipe_When() {
	fl := New(
		Pipe{
			Name: option.Some("flow_dog.js"),
			// only nodes that meet this condition will be processed by this pipe.
			// in this case, all nodes whose "kind" is equal to "dog" will be processed.
			When: option.Some(When{Kind: option.Some([]Kind{"dog"})}),
			Code: option.Some(`
				export default function main(nodes, next) {
					for (const node of nodes) {
						node.meta = {sayword: "woof"}
					}
				}
			`),
		},
		Pipe{
			Name: option.Some("flow_cat.js"),
			// in this case, all nodes whose "kind" is equal to "cat" will be processed.
			When: option.Some(When{Kind: option.Some([]Kind{"cat"})}),
			Code: option.Some(`
				export default function main(nodes, next) {
					for (const node of nodes) {
						node.meta = {sayword: "meow"}
					}
				}
			`),
		},
	)

	target := []Node{
		{
			UUID: option.Some(MustUUID("5f883c2d-52a6-4d55-84b5-6763b2717f86")),
			Kind: option.Some("dog"),
		},
		{
			UUID: option.Some(MustUUID("5719e68a-963b-4c88-9018-7ee216454350")),
			Kind: option.Some("cat"),
		},
	}
	err := fl.Work(context.Background(), target)
	if err != nil {
		panic(err)
	}

	fmt.Println("dog sayword", target[0].Meta.Get()["sayword"])
	fmt.Println("cat sayword", target[1].Meta.Get()["sayword"])
	// Output:
	// dog sayword woof
	// cat sayword meow
}

// ExamplePipe_Next
// How to implement a chain of responsibility pattern
func ExamplePipe_Next() {
	fl := New(
		Pipe{
			Name: option.Some("first.js"),
			// this pipe will process all nodes due to an empty condition.
			When: option.Some(When{}),
			Next: option.Some([]Name{"second.js"}),
			Code: option.Some(`
				export default function main(nodes, next) {
					for (const node of nodes) {
						node.meta ??= {value: []}
						node.meta.value.push("first")
					}
					next(nodes)
				}
			`),
		},
		Pipe{
			Name: option.Some("second.js"),
			// this pipe will not process nodes, but will be called after the upstream pipe because it is referenced in the Next property.
			When: option.None[When](),
			Code: option.Some(`
				export default function main(nodes) {
					for (const node of nodes) {
						node.meta ??= {value: []}
						node.meta.value.push("second")
					}
				}
			`),
		},
	)

	target := []Node{
		{
			UUID: option.Some(MustUUID("1bcc8106-fec1-4b21-9014-b53a2f9ed5a2")),
		},
	}
	err := fl.Work(context.Background(), target)
	if err != nil {
		panic(err)
	}

	fmt.Println(target[0].Meta.Get()["value"])
	// Output:
	// [first second]
}

// ExamplePlugin
// How to extend the js runtime
func ExamplePlugin() {
	fl := New(
		Plugin{
			// will be called once after compile and create a goja.Runtime
			Init: func(ctx context.Context, t Api) error {
				t.Runtime().Set("formatBar", t.Runtime().ToValue(func(goja.FunctionCall) goja.Value {
					return t.Runtime().ToValue("BAR")
				}))
				return nil
			},
			// will be called every time before executing a script
			// it may be useful to set values in this context
			Call: func(ctx context.Context, t Api) error {
				_ = t.This().Set("foo", "FOO")
				return nil
			},
			// will be called every time after executing a script
			// it might be useful to clear all the values that were defined in the Call section.
			Quit: func(ctx context.Context, t Api) error {
				_ = t.This().Delete("foo")
				return nil
			},
		},
		Pipe{
			Name: option.Some("first.js"),
			// this pipe will process all nodes due to an empty condition.
			When: option.Some(When{}),
			Code: option.Some(`
				export default function main(nodes, next) {
					for (const node of nodes) {
						node.meta = {value: formatBar()+this.foo}
					}
				}
			`),
		},
	)

	target := []Node{
		{
			UUID: option.Some(MustUUID("1bcc8106-fec1-4b21-9014-b53a2f9ed5a2")),
		},
	}
	err := fl.Work(context.Background(), target)
	if err != nil {
		panic(err)
	}

	fmt.Println(target[0].Meta.Get()["value"])
	// Output:
	// BARFOO
}
