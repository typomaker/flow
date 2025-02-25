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
		WithPipe(Pipe{
			UUID: option.Some(MustUUID("ffe9a397-e091-4254-8da4-be0dcf33f481")),
			When: option.Some(When{}),
			Code: option.Some(`
				export default function main(nodes){
					nodes[0].kind="foo"
				}
			`),
		}),
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
		WithPipe(Pipe{
			UUID: option.Some(MustUUID("ffe9a397-e091-4254-8da4-be0dcf33f481")),
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
		}),
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
		WithPipe(Pipe{
			UUID: option.Some(MustUUID("ffe9a397-e091-4254-8da4-be0dcf33f481")),
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
		}),
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
		WithPipe(Pipe{
			UUID: option.Some(MustUUID("ffe9a397-e091-4254-8da4-be0dcf33f481")),
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
		}),
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
		WithPipe(Pipe{
			UUID: option.Some(MustUUID("0e4f9350-f007-474f-b051-d2510e522800")),
			When: option.Some(When{}),
			Code: option.Some(`
				export default function main(nodes, next){
					nodes[0].meta ??= {seq: []}
					nodes[0].meta.seq.push(this.FLOW_PIPE_SELF_UUID)
					next(nodes)
				}
			`),
			Next: option.Some([]UUID{
				MustUUID("1e4f9350-f007-474f-b051-d2510e522800"),
				MustUUID("2e4f9350-f007-474f-b051-d2510e522800"),
			}),
		}),
		WithPipe(Pipe{
			UUID: option.Some(MustUUID("1e4f9350-f007-474f-b051-d2510e522800")),
			Code: option.Some(`
				export default function main(nodes, next){
					nodes[0].meta ??= {seq: []}
					nodes[0].meta.seq.push(this.FLOW_PIPE_SELF_UUID)
					next(nodes)
				}
			`),
			Next: option.Some([]UUID{
				MustUUID("3e4f9350-f007-474f-b051-d2510e522800"),
				MustUUID("1e4f9350-f007-474f-b051-d2510e522800"),
			}),
		}),
		WithPipe(Pipe{
			UUID: option.Some(MustUUID("2e4f9350-f007-474f-b051-d2510e522800")),
			Code: option.Some(`
				export default function main(nodes, next){
					nodes[0].meta ??= {seq: []}
					nodes[0].meta.seq.push(this.FLOW_PIPE_SELF_UUID)
					next(nodes)
				}
			`),
			Next: option.Some([]UUID{
				MustUUID("0e4f9350-f007-474f-b051-d2510e522800"),
			}),
		}),
		WithPipe(Pipe{
			UUID: option.Some(MustUUID("3e4f9350-f007-474f-b051-d2510e522800")),
			Next: option.Some([]UUID{
				MustUUID("9e4f9350-f007-474f-b051-d2510e522800"),
			}),
		}),
		WithPipe(Pipe{
			UUID: option.Some(MustUUID("9e4f9350-f007-474f-b051-d2510e522800")),
			Code: option.Some(`
				export default function main(nodes, next){
					nodes[0].meta ??= {seq: []}
					nodes[0].meta.seq.push(this.FLOW_PIPE_SELF_UUID)
					next(nodes)
				}
			`),
		}),
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
				"0e4f9350-f007-474f-b051-d2510e522800",
				"1e4f9350-f007-474f-b051-d2510e522800",
				"2e4f9350-f007-474f-b051-d2510e522800",
				"9e4f9350-f007-474f-b051-d2510e522800",
			}}),
		},
	}
	require.Equal(t, e, a)
}
func TestPipeDedup(t *testing.T) {
	f := New(
		WithPipe(
			Pipe{
				UUID: option.Some(MustUUID("0e4f9350-f007-474f-b051-d2510e522800")),
				Name: option.Some[Name]("foo"),
			},
			Pipe{
				UUID: option.Some(MustUUID("0e4f9350-f007-474f-b051-d2510e522800")),
				Name: option.Some[Name]("bar"),
			},
			Pipe{
				UUID: option.Some(MustUUID("1e4f9350-f007-474f-b051-d2510e522800")),
			},
		),
	)
	require.Len(t, f.stock, 2)
	require.Equal(t, "0e4f9350-f007-474f-b051-d2510e522800", f.stock[0].UUID.Get().String())
	require.Equal(t, "foo", f.stock[0].Name.Get().String())
	require.Equal(t, "1e4f9350-f007-474f-b051-d2510e522800", f.stock[1].UUID.Get().String())
}
