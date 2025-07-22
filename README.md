# Flow
A tool that allows you to group and execute your business logic into separate JS scripts.

The approach is similar to the one used in the gaming engines, when the behavior of game objects is described in some kind of scripting language. This reduces the delivery time for updates. It allows you to easily extend the business logic and also simply undo changes.

**This library is under development.**

## Examples

### Process a node with a specific pipe.
```go
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
```
