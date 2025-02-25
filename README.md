# Flow
A tool that allows you to group and execute your business logic into separate JS scripts.

The approach is similar to the one used in the gaming industry, when the behavior of game objects is described in some kind of scripting language. This reduces the delivery time for updates. It allows you to easily extend the business logic and also simply undo changes.

**This library is under development.**

## Examples

### Process a node with a specific pipe.
```go
fl := New(
    Pipe{
        UUID: option.Some(MustUUID("7ae0a3ea-e17a-44c6-ba30-641df6fdd26d")),
        Name: option.Some[Name]("flow_dog.js"),
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
        UUID: option.Some(MustUUID("eef97892-b1b5-422a-8d4c-1fba9a2950ce")),
        Name: option.Some[Name]("flow_cat.js"),
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
```

### Implement a chain of responsibility pattern.
```go
fl := New(
    Pipe{
        UUID: option.Some(MustUUID("7ae0a3ea-e17a-44c6-ba30-641df6fdd26d")),
        Name: option.Some[Name]("first.js"),
        // this pipe will process all nodes due to an empty condition.
        When: option.Some(When{}),
        Next: option.Some([]UUID{MustUUID("eef97892-b1b5-422a-8d4c-1fba9a2950ce")}),
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
        UUID: option.Some(MustUUID("eef97892-b1b5-422a-8d4c-1fba9a2950ce")),
        Name: option.Some[Name]("second.js"),
        // this pipe will not process nodes according to an unreachable condition, but it will be executed after the above pipe, since it refers to this through the Next property.
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
```

### Extend the js runtime
```go
fl := New(
    Plugin{
        // will be called once after compile and create a goja.Runtime
        Init: func(ctx context.Context, t Api) error {
            t.This().Set("formatBar", t.Goja().ToValue(func(goja.FunctionCall) goja.Value {
                return t.Goja().ToValue("BAR")
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
        UUID: option.Some(MustUUID("7ae0a3ea-e17a-44c6-ba30-641df6fdd26d")),
        Name: option.Some[Name]("first.js"),
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
```