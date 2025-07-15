package test

import (
	"github.com/typomaker/flow"
	"github.com/typomaker/flow/gojaflow"
)

func init() {
	var stmt = flow.Compose(
		flow.If(
			flow.UUID.In(flow.MustUUID("6d54d24c-f132-44c7-82c9-f52df81cbfb6")),
			gojaflow.New("ffff"),
		),
		flow.If(
			flow.And(
				flow.Not(
					flow.UUID.In(flow.MustUUID("6d54d24c-f132-44c7-82c9-f52df81cbfb6")),
				),
				flow.Hook.In(
					flow.Hook{"a": true},
					flow.Hook{"b": true},
				),
			),
			gojaflow.New("ffff"),
		),
	)

	var ctx = flow.Context{}
	var target = flow.Node{}
	stmt(ctx, &target, flow.Next(func(n *flow.Node) error { return nil }))
}
