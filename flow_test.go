package flow

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typomaker/option"
)

func TestNil(t *testing.T) {
	ctx := context.Background()
	f := New()
	err := f.Work(ctx, nil)
	require.NoError(t, err)
}

func TestNoop(t *testing.T) {
	ctx := context.Background()
	f := New()
	a := []Node{
		{UUID: option.Some(MustUUID("aee6576f-19f8-419c-b3f8-41b770006332"))},
	}
	err := f.Work(ctx, a)
	require.NoError(t, err)
	e := []Node{
		{UUID: option.Some(MustUUID("aee6576f-19f8-419c-b3f8-41b770006332"))},
	}
	require.Equal(t, e, a)
}
func TestPipeWhenZero(t *testing.T) {
	ctx := context.Background()
	f := New(
		Pipe{
			Name: option.Some("foo"),
			When: option.Some(When{}),
			Code: option.Some(`
				export default function main(nodes){
					nodes[0].kind="foo"
				}
			`),
		},
	)
	a := []Node{
		{UUID: option.Some(MustUUID("bee6576f-19f8-419c-b3f8-41b770006332"))},
	}
	err := f.Work(ctx, a)
	require.NoError(t, err)
	e := []Node{
		{UUID: option.Some(MustUUID("bee6576f-19f8-419c-b3f8-41b770006332")), Kind: option.Some("foo")},
	}
	require.Equal(t, e, a)
}
func TestPipeWhenUUIDSome(t *testing.T) {
	ctx := context.Background()
	f := New(
		Pipe{
			Name: option.Some("foo"),
			When: option.Some(When{
				UUID: option.Some([]UUID{
					MustUUID("bee6576f-19f8-419c-b3f8-41b770006332"),
				}),
			}),
			Code: option.Some(`
				export default function main(nodes){
					nodes[0].kind="foo"
				}
			`),
		},
	)
	a := []Node{
		{UUID: option.Some(MustUUID("aee6576f-19f8-419c-b3f8-41b770006332"))},
		{UUID: option.Some(MustUUID("bee6576f-19f8-419c-b3f8-41b770006332"))},
	}
	err := f.Work(ctx, a)
	require.NoError(t, err)
	e := []Node{
		{UUID: option.Some(MustUUID("aee6576f-19f8-419c-b3f8-41b770006332"))},
		{UUID: option.Some(MustUUID("bee6576f-19f8-419c-b3f8-41b770006332")), Kind: option.Some("foo")},
	}
	require.Equal(t, e, a)
}
func TestPipeWhenKindSome(t *testing.T) {
	ctx := context.Background()
	f := New(
		Pipe{
			Name: option.Some("foo"),
			When: option.Some(When{
				Kind: option.Some([]Kind{
					"bar",
				}),
			}),
			Code: option.Some(`
				export default function main(nodes){
					nodes[0].kind="foo"
				}
			`),
		},
	)
	a := []Node{
		{UUID: option.Some(MustUUID("aee6576f-19f8-419c-b3f8-41b770006332")), Kind: option.Some("buz")},
		{UUID: option.Some(MustUUID("bee6576f-19f8-419c-b3f8-41b770006332")), Kind: option.Some("bar")},
	}
	err := f.Work(ctx, a)
	require.NoError(t, err)
	e := []Node{
		{UUID: option.Some(MustUUID("aee6576f-19f8-419c-b3f8-41b770006332")), Kind: option.Some("buz")},
		{UUID: option.Some(MustUUID("bee6576f-19f8-419c-b3f8-41b770006332")), Kind: option.Some("foo")},
	}
	require.Equal(t, e, a)
}
func TestPipeWhenHookSome(t *testing.T) {
	ctx := context.Background()
	f := New(
		Pipe{
			Name: option.Some("foo"),
			When: option.Some(When{
				Hook: option.Some([]Hook{
					{"foo": "bar"},
				}),
			}),
			Code: option.Some(`
				export default function main(nodes){
					nodes[0].kind="foo"
				}
			`),
		},
	)
	a := []Node{
		{UUID: option.Some(MustUUID("aee6576f-19f8-419c-b3f8-41b770006332")), Hook: option.Some(Hook{"foo": "buz"})},
		{UUID: option.Some(MustUUID("bee6576f-19f8-419c-b3f8-41b770006332")), Hook: option.Some(Hook{"foo": "bar"})},
	}
	err := f.Work(ctx, a)
	require.NoError(t, err)
	e := []Node{
		{UUID: option.Some(MustUUID("aee6576f-19f8-419c-b3f8-41b770006332")), Hook: option.Some(Hook{"foo": "buz"})},
		{UUID: option.Some(MustUUID("bee6576f-19f8-419c-b3f8-41b770006332")), Kind: option.Some("foo"), Hook: option.Some(Hook{"foo": "bar"})},
	}
	require.Equal(t, e, a)
}
func TestPipeNext(t *testing.T) {
	ctx := context.Background()
	f := New(
		Pipe{
			Name: option.Some("f1"),
			When: option.Some(When{}),
			Code: option.Some(`
				export default function main(nodes, next){
					nodes[0].meta ??= {seq: []}
					nodes[0].meta.seq.push(this.FLOW_PIPE_NAME)
					next(nodes)
				}
			`),
			Next: option.Some([]Name{"f2", "f3"}),
		},
		Pipe{
			Name: option.Some("f2"),
			Code: option.Some(`
				export default function main(nodes, next){
					nodes[0].meta ??= {seq: []}
					nodes[0].meta.seq.push(this.FLOW_PIPE_NAME)
					next(nodes)
				}
			`),
			Next: option.Some([]Name{"f4", "f2"}),
		},
		Pipe{
			Name: option.Some("f3"),
			Code: option.Some(`
				export default function main(nodes, next){
					nodes[0].meta ??= {seq: []}
					nodes[0].meta.seq.push(this.FLOW_PIPE_NAME)
					next(nodes)
				}
			`),
			Next: option.Some([]Name{"f1"}),
		},
		Pipe{
			Name: option.Some("f4"),
			Next: option.Some([]Name{"f5"}),
		},
		Pipe{
			Name: option.Some("f5"),
			Code: option.Some(`
				export default function main(nodes, next){
					nodes[0].meta ??= {seq: []}
					nodes[0].meta.seq.push(this.FLOW_PIPE_NAME)
					next(nodes)
				}
			`),
		},
	)
	a := []Node{
		{UUID: option.Some(MustUUID("aee6576f-19f8-419c-b3f8-41b770006332"))},
	}
	err := f.Work(ctx, a)
	require.NoError(t, err)
	e := []Node{
		{
			UUID: option.Some(MustUUID("aee6576f-19f8-419c-b3f8-41b770006332")),
			Meta: option.Some(Meta{"seq": []any{
				"f1",
				"f2",
				"f3",
				"f5",
			}}),
		},
	}
	require.Equal(t, e, a)
}
func TestPipeDedup(t *testing.T) {
	f := New(
		Pipe{
			Name: option.Some[Name]("foo"),
			Code: option.Some("foo"),
		},
		Pipe{
			Name: option.Some[Name]("foo"),
			Code: option.Some("bar"),
		},
		Pipe{
			Name: option.Some("buz"),
			Code: option.Some("buz"),
		},
	)
	require.Len(t, f.stock, 2)
	require.Equal(t, "foo", f.stock[0].Name.Get())
	require.Equal(t, "foo", f.stock[0].Code.Get())
	require.Equal(t, "buz", f.stock[1].Name.Get())
	require.Equal(t, "buz", f.stock[1].Code.Get())
}
