package flow

import (
	"context"
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/require"
	"github.com/typomaker/option"
)

func TestInit(t *testing.T) {
	ctx := context.Background()
	f := New(
		Pipe{
			UUID: option.Some(MustUUID("afe9a397-e091-4254-8da4-be0dcf33f481")),
			When: option.Some(When{}),
			Code: option.Some(`
				export default function main(nodes){
					nodes[0].kind=initPlugin1()+initPlugin2()
				}
			`),
		},
		Plugin{
			Init: func(ctx context.Context, x Api) error {
				require.NotNil(t, x.Flow())
				require.NotNil(t, x.Goja())
				require.NotZero(t, x.Pipe())
				require.NotNil(t, x.This())

				x.This().Set("initPlugin1", x.Goja().ToValue(func(call goja.ConstructorCall) goja.Value {
					return x.Goja().ToValue("foo")
				}))
				return nil
			},
		},
		Plugin{
			Init: func(ctx context.Context, x Api) error {
				require.NotNil(t, x.Flow())
				require.NotNil(t, x.Goja())
				require.NotZero(t, x.Pipe())
				require.NotNil(t, x.This())

				x.This().Set("initPlugin2", x.Goja().ToValue(func(call goja.ConstructorCall) goja.Value {
					return x.Goja().ToValue("bar")
				}))
				return nil
			},
		},
	)
	a := []Node{{}}
	err := f.Work(ctx, a)
	require.NoError(t, err)
	require.Equal(t, a[0].Kind.Get(), "foobar")
}
func TestCall(t *testing.T) {
	ctx := context.Background()
	f := New(
		Pipe{
			UUID: option.Some(MustUUID("afe9a397-e091-4254-8da4-be0dcf33f481")),
			When: option.Some(When{}),
			Code: option.Some(`
				export default function main(nodes){
					nodes[0].kind=this.callPlugin1()+this.callPlugin2()
				}
			`),
		},
		Plugin{
			Call: func(ctx context.Context, x Api) error {
				require.NotNil(t, x.Flow())
				require.NotNil(t, x.Goja())
				require.NotZero(t, x.Pipe())
				require.NotNil(t, x.This())
				x.This().Set("callPlugin1", x.Goja().ToValue(func(call goja.ConstructorCall) goja.Value {
					return x.Goja().ToValue("foo")
				}))
				return nil
			},
		},
		Plugin{
			Call: func(ctx context.Context, x Api) error {
				require.NotNil(t, x.Flow())
				require.NotNil(t, x.Goja())
				require.NotZero(t, x.Pipe())
				require.NotNil(t, x.This())
				x.This().Set("callPlugin2", x.Goja().ToValue(func(call goja.ConstructorCall) goja.Value {
					return x.Goja().ToValue("bar")
				}))
				return nil
			},
		},
	)
	a := []Node{{}}
	err := f.Work(ctx, a)
	require.NoError(t, err)
	require.Equal(t, a[0].Kind.Get(), "foobar")
}
func TestQuit(t *testing.T) {
	ctx := context.Background()
	var quitPlugin1, quitPlugin2 string
	f := New(
		Pipe{
			UUID: option.Some(MustUUID("afe9a397-e091-4254-8da4-be0dcf33f481")),
			When: option.Some(When{}),
			Code: option.Some(`
				export default function main(nodes){
					nodes[0].kind="foobar"
					this.quitPlugin1="foo"
					this.quitPlugin2="bar"
				}
			`),
		},
		Plugin{
			Quit: func(ctx context.Context, x Api) error {
				require.NotNil(t, x.Flow())
				require.NotNil(t, x.Goja())
				require.NotZero(t, x.Pipe())
				require.NotNil(t, x.This())
				err := x.Goja().ExportTo(x.This().Get("quitPlugin1"), &quitPlugin1)
				require.NoError(t, err)
				return nil
			},
		},
		Plugin{
			Quit: func(ctx context.Context, x Api) error {
				require.NotNil(t, x.Flow())
				require.NotNil(t, x.Goja())
				require.NotZero(t, x.Pipe())
				require.NotNil(t, x.This())
				err := x.Goja().ExportTo(x.This().Get("quitPlugin2"), &quitPlugin2)
				require.NoError(t, err)
				return nil
			},
		},
	)
	a := []Node{{}}
	err := f.Work(ctx, a)
	require.NoError(t, err)
	require.Equal(t, a[0].Kind.Get(), "foobar")
	require.Equal(t, "foo", quitPlugin1)
	require.Equal(t, "bar", quitPlugin2)
}
