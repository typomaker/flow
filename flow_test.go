package flow

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typomaker/option"
)

func TestPipeBeforeNext(t *testing.T) {
	c := []int{}
	f := Pipe(
		func(ctx context.Context, target []Node, next Next) (err error) {
			c = append(c, 0)
			next(target)
			return nil
		},
		func(ctx context.Context, target []Node, next Next) (err error) {
			c = append(c, 1)
			next(target)
			return nil
		},
		func(ctx context.Context, target []Node, next Next) (err error) {
			c = append(c, 2)
			next(target)
			return nil
		},
	)
	next := Next(func(target []Node) error {
		c = append(c, 3)
		return nil
	})
	ctx := context.Background()
	target := []Node{}
	err := f(ctx, target, next)
	require.NoError(t, err)
	require.Equal(t, []int{0, 1, 2, 3}, c)
}
func TestPipeAfterNext(t *testing.T) {
	c := []int{}
	f := Pipe(
		func(ctx context.Context, target []Node, next Next) (err error) {
			next(target)
			c = append(c, 3)
			return nil
		},
		func(ctx context.Context, target []Node, next Next) (err error) {
			next(target)
			c = append(c, 2)
			return nil
		},
		func(ctx context.Context, target []Node, next Next) (err error) {
			next(target)
			c = append(c, 1)
			return nil
		},
	)
	next := Next(func(target []Node) error {
		c = append(c, 0)
		return nil
	})
	ctx := context.Background()
	target := []Node{}
	err := f(ctx, target, next)
	require.NoError(t, err)
	require.Equal(t, []int{0, 1, 2, 3}, c)
}
func TestAndCallNextLater(t *testing.T) {
	c := []int{}
	f := And(
		func(ctx context.Context, target []Node, next Next) (err error) {
			c = append(c, 0)
			next(target)
			return nil
		},
		func(ctx context.Context, target []Node, next Next) (err error) {
			c = append(c, 1)
			next(target)
			return nil
		},
		func(ctx context.Context, target []Node, next Next) (err error) {
			c = append(c, 2)
			next(target)
			return nil
		},
	)
	next := Next(func(target []Node) error {
		c = append(c, 3)
		return nil
	})
	ctx := context.Background()
	target := []Node{
		{UUID: option.Some(MustUUID("e42e04d7-d016-4175-bf3a-e0201ff9f6a7"))},
	}
	err := f(ctx, target, next)
	require.NoError(t, err)
	require.Equal(t, []int{0, 1, 2, 3}, c)
}
func TestAndCallNextEarly(t *testing.T) {
	c := []int{}
	f := And(
		func(ctx context.Context, target []Node, next Next) (err error) {
			next(target)
			c = append(c, 3)
			return nil
		},
		func(ctx context.Context, target []Node, next Next) (err error) {
			next(target)
			c = append(c, 2)
			return nil
		},
		func(ctx context.Context, target []Node, next Next) (err error) {
			next(target)
			c = append(c, 1)
			return nil
		},
	)
	next := Next(func(target []Node) error {
		c = append(c, 0)
		return nil
	})
	ctx := context.Background()
	target := []Node{
		{UUID: option.Some(MustUUID("e42e04d7-d016-4175-bf3a-e0201ff9f6a7"))},
	}
	err := f(ctx, target, next)
	require.NoError(t, err)
	require.Equal(t, []int{0, 1, 2, 3}, c)
}
func TestAndCallNextAfterBreak(t *testing.T) {
	c := []int{}
	f := And(
		func(ctx context.Context, target []Node, next Next) (err error) {
			next(target)
			c = append(c, 3)
			return nil
		},
		func(ctx context.Context, target []Node, next Next) (err error) {
			c = append(c, 2)
			return nil
		},
		func(ctx context.Context, target []Node, next Next) (err error) {
			next(target)
			c = append(c, 1)
			return nil
		},
	)
	next := Next(func(target []Node) error {
		c = append(c, 0)
		return nil
	})
	ctx := context.Background()
	target := []Node{
		{UUID: option.Some(MustUUID("e42e04d7-d016-4175-bf3a-e0201ff9f6a7"))},
	}
	err := f(ctx, target, next)
	require.NoError(t, err)
	require.Equal(t, []int{2, 3, 0}, c)
}
func TestAndCallNextOnEmptyInput(t *testing.T) {
	c := []int{}
	f := And()
	next := Next(func(target []Node) error {
		c = append(c, 0)
		return nil
	})
	ctx := context.Background()
	target := []Node{
		{UUID: option.Some(MustUUID("e42e04d7-d016-4175-bf3a-e0201ff9f6a7"))},
	}
	err := f(ctx, target, next)
	require.NoError(t, err)
	require.Equal(t, []int{0}, c)
}

// func TestNil(t *testing.T) {
// 	ctx := context.Background()
// 	f := New()
// 	err := f.Work(ctx, nil)
// 	require.NoError(t, err)
// }

// func TestNoop(t *testing.T) {
// 	ctx := context.Background()
// 	f := New()
// 	a := []Node{
// 		{UUID: option.Some(MustUUID("aee6576f-19f8-419c-b3f8-41b770006332"))},
// 	}
// 	err := f.Work(ctx, a)
// 	require.NoError(t, err)
// 	e := []Node{
// 		{UUID: option.Some(MustUUID("aee6576f-19f8-419c-b3f8-41b770006332"))},
// 	}
// 	require.Equal(t, e, a)
// }
// func TestPipeWhenZero(t *testing.T) {
// 	ctx := context.Background()
// 	f := New(
// 		Pipe{
// 			Name: option.Some("foo"),
// 			When: option.Some(When{}),
// 			Code: option.Some(`
// 				export default function main(nodes){
// 					nodes[0].meta={val: "foo"}
// 				}
// 			`),
// 		},
// 	)
// 	a := []Node{
// 		{UUID: option.Some(MustUUID("bee6576f-19f8-419c-b3f8-41b770006332"))},
// 	}
// 	err := f.Work(ctx, a)
// 	require.NoError(t, err)
// 	e := []Node{
// 		{UUID: option.Some(MustUUID("bee6576f-19f8-419c-b3f8-41b770006332")), Meta: option.Some(Meta{"val": "foo"})},
// 	}
// 	require.Equal(t, e, a)
// }
// func TestPipeWhenUUIDSome(t *testing.T) {
// 	ctx := context.Background()
// 	f := New(
// 		Pipe{
// 			Name: option.Some("foo"),
// 			When: option.Some(When{
// 				UUID: option.Some([]UUID{
// 					MustUUID("bee6576f-19f8-419c-b3f8-41b770006332"),
// 				}),
// 			}),
// 			Code: option.Some(`
// 				export default function main(nodes){
// 					nodes[0].meta = {val:"foo"}
// 				}
// 			`),
// 		},
// 	)
// 	a := []Node{
// 		{UUID: option.Some(MustUUID("aee6576f-19f8-419c-b3f8-41b770006332"))},
// 		{UUID: option.Some(MustUUID("bee6576f-19f8-419c-b3f8-41b770006332"))},
// 	}
// 	err := f.Work(ctx, a)
// 	require.NoError(t, err)
// 	e := []Node{
// 		{UUID: option.Some(MustUUID("aee6576f-19f8-419c-b3f8-41b770006332"))},
// 		{UUID: option.Some(MustUUID("bee6576f-19f8-419c-b3f8-41b770006332")), Meta: option.Some(Meta{"val": "foo"})},
// 	}
// 	require.Equal(t, e, a)
// }
// func TestPipeWhenHookSome(t *testing.T) {
// 	ctx := context.Background()
// 	f := New(
// 		Pipe{
// 			Name: option.Some("foo"),
// 			When: option.Some(When{
// 				Hook: option.Some([]Hook{
// 					{"foo": "bar"},
// 				}),
// 			}),
// 			Code: option.Some(`
// 				export default function main(nodes){
// 					nodes[0].meta={val:"foo"}
// 				}
// 			`),
// 		},
// 	)
// 	a := []Node{
// 		{UUID: option.Some(MustUUID("aee6576f-19f8-419c-b3f8-41b770006332")), Hook: option.Some(Hook{"foo": "buz"})},
// 		{UUID: option.Some(MustUUID("bee6576f-19f8-419c-b3f8-41b770006332")), Hook: option.Some(Hook{"foo": "bar"})},
// 	}
// 	err := f.Work(ctx, a)
// 	require.NoError(t, err)
// 	e := []Node{
// 		{UUID: option.Some(MustUUID("aee6576f-19f8-419c-b3f8-41b770006332")), Hook: option.Some(Hook{"foo": "buz"})},
// 		{UUID: option.Some(MustUUID("bee6576f-19f8-419c-b3f8-41b770006332")), Meta: option.Some(Meta{"val": "foo"}), Hook: option.Some(Hook{"foo": "bar"})},
// 	}
// 	require.Equal(t, e, a)
// }
// func TestPipeNext(t *testing.T) {
// 	ctx := context.Background()
// 	f := New(
// 		Pipe{
// 			Name: option.Some("f1"),
// 			When: option.Some(When{}),
// 			Code: option.Some(`
// 				export default function main(nodes, next){
// 					nodes[0].meta ??= {seq: []}
// 					nodes[0].meta.seq.push(this.FLOW_PIPE_NAME)
// 					next(nodes)
// 				}
// 			`),
// 			Next: option.Some([]Name{"f2", "f3"}),
// 		},
// 		Pipe{
// 			Name: option.Some("f2"),
// 			Code: option.Some(`
// 				export default function main(nodes, next){
// 					nodes[0].meta ??= {seq: []}
// 					nodes[0].meta.seq.push(this.FLOW_PIPE_NAME)
// 					next(nodes)
// 				}
// 			`),
// 			Next: option.Some([]Name{"f4", "f2"}),
// 		},
// 		Pipe{
// 			Name: option.Some("f3"),
// 			Code: option.Some(`
// 				export default function main(nodes, next){
// 					nodes[0].meta ??= {seq: []}
// 					nodes[0].meta.seq.push(this.FLOW_PIPE_NAME)
// 					next(nodes)
// 				}
// 			`),
// 			Next: option.Some([]Name{"f1"}),
// 		},
// 		Pipe{
// 			Name: option.Some("f4"),
// 			Next: option.Some([]Name{"f5"}),
// 		},
// 		Pipe{
// 			Name: option.Some("f5"),
// 			Code: option.Some(`
// 				export default function main(nodes, next){
// 					nodes[0].meta ??= {seq: []}
// 					nodes[0].meta.seq.push(this.FLOW_PIPE_NAME)
// 					next(nodes)
// 				}
// 			`),
// 		},
// 	)
// 	a := []Node{
// 		{UUID: option.Some(MustUUID("aee6576f-19f8-419c-b3f8-41b770006332"))},
// 	}
// 	err := f.Work(ctx, a)
// 	require.NoError(t, err)
// 	e := []Node{
// 		{
// 			UUID: option.Some(MustUUID("aee6576f-19f8-419c-b3f8-41b770006332")),
// 			Meta: option.Some(Meta{"seq": []any{
// 				"f1",
// 				"f2",
// 				"f3",
// 				"f5",
// 			}}),
// 		},
// 	}
// 	require.Equal(t, e, a)
// }
// func TestPipeDedup(t *testing.T) {
// 	f := New(
// 		Pipe{
// 			Name: option.Some[Name]("foo"),
// 			Code: option.Some("foo"),
// 		},
// 		Pipe{
// 			Name: option.Some[Name]("foo"),
// 			Code: option.Some("bar"),
// 		},
// 		Pipe{
// 			Name: option.Some("buz"),
// 			Code: option.Some("buz"),
// 		},
// 	)
// 	require.Len(t, f.stock, 2)
// 	require.Equal(t, "foo", f.stock[0].Name.Get())
// 	require.Equal(t, "foo", f.stock[0].Code.Get())
// 	require.Equal(t, "buz", f.stock[1].Name.Get())
// 	require.Equal(t, "buz", f.stock[1].Code.Get())
// }
// func TestPriority(t *testing.T) {
// 	f := New(
// 		Pipe{
// 			Name: option.Some("f1"),
// 			Code: option.Some(`export default function main(nodes) {}`),
// 		},
// 		Pipe{
// 			Name: option.Some("f2"),
// 			Code: option.Some(`export default function main(nodes) {}`),
// 			Next: option.Some([]string{"f1", "f3"}),
// 		},
// 		Pipe{
// 			Name: option.Some("f3"),
// 			Code: option.Some(`export default function main(nodes) {}`),
// 			Next: option.Some([]string{"f1", "f4"}),
// 		},
// 		Pipe{
// 			Name: option.Some("f4"),
// 			Code: option.Some(`export default function main(nodes) {}`),
// 			When: option.Some(When{}),
// 		},
// 		Pipe{
// 			Name: option.Some("f5"),
// 			Code: option.Some(`export default function main(nodes) {}`),
// 			When: option.Some(When{Hook: option.Some([]Hook{})}),
// 		},
// 		Pipe{
// 			Name: option.Some("f8"),
// 			Code: option.Some(`export default function main(nodes) {}`),
// 			When: option.Some(When{UUID: option.Some([]UUID{})}),
// 		},
// 		Pipe{
// 			Name: option.Some("f10"),
// 			Code: option.Some(`export default function main(nodes) {}`),
// 			When: option.Some(When{UUID: option.Some([]UUID{}), Hook: option.Some([]Hook{})}),
// 		},
// 	)
// 	require.Equal(t, "f10", f.stock[0].Name.Get())
// 	require.Equal(t, "f8", f.stock[1].Name.Get())
// 	require.Equal(t, "f5", f.stock[2].Name.Get())
// 	require.Equal(t, "f4", f.stock[3].Name.Get())
// 	require.Equal(t, "f2", f.stock[4].Name.Get())
// 	require.Equal(t, "f3", f.stock[5].Name.Get())
// 	require.Equal(t, "f1", f.stock[6].Name.Get())
// }
