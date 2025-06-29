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
			Name: option.Some("foo"),
			When: option.Some(When{}),
			Code: option.Some(`
				export default function main(nodes){
					nodes[0].meta = {val: initPlugin1()+initPlugin2()}
				}
			`),
		},
		Plugin{
			Init: func(ctx context.Context, x Api) error {
				require.NotNil(t, x.Flow())
				require.NotNil(t, x.Runtime())
				require.NotZero(t, x.Pipe())
				require.NotNil(t, x.This())

				x.This().Set("initPlugin1", x.Runtime().ToValue(func(call goja.ConstructorCall) goja.Value {
					return x.Runtime().ToValue("foo")
				}))
				return nil
			},
		},
		Plugin{
			Init: func(ctx context.Context, x Api) error {
				require.NotNil(t, x.Flow())
				require.NotNil(t, x.Runtime())
				require.NotZero(t, x.Pipe())
				require.NotNil(t, x.This())

				x.This().Set("initPlugin2", x.Runtime().ToValue(func(call goja.ConstructorCall) goja.Value {
					return x.Runtime().ToValue("bar")
				}))
				return nil
			},
		},
	)
	a := []Node{{}}
	err := f.Work(ctx, a)
	require.NoError(t, err)
	require.Equal(t, a[0].Meta.Get()["val"], "foobar")
}
func TestCall(t *testing.T) {
	ctx := context.Background()
	f := New(
		Pipe{
			Name: option.Some("foo"),
			When: option.Some(When{}),
			Code: option.Some(`
				export default function main(nodes){
					nodes[0].meta={val:this.callPlugin1()+this.callPlugin2()}
				}
			`),
		},
		Plugin{
			Call: func(ctx context.Context, x Api) error {
				require.NotNil(t, x.Flow())
				require.NotNil(t, x.Runtime())
				require.NotZero(t, x.Pipe())
				require.NotNil(t, x.This())
				x.This().Set("callPlugin1", x.Runtime().ToValue(func(call goja.ConstructorCall) goja.Value {
					return x.Runtime().ToValue("foo")
				}))
				return nil
			},
		},
		Plugin{
			Call: func(ctx context.Context, x Api) error {
				require.NotNil(t, x.Flow())
				require.NotNil(t, x.Runtime())
				require.NotZero(t, x.Pipe())
				require.NotNil(t, x.This())
				x.This().Set("callPlugin2", x.Runtime().ToValue(func(call goja.ConstructorCall) goja.Value {
					return x.Runtime().ToValue("bar")
				}))
				return nil
			},
		},
	)
	a := []Node{{}}
	err := f.Work(ctx, a)
	require.NoError(t, err)
	require.Equal(t, a[0].Meta.Get()["val"], "foobar")
}
func TestQuit(t *testing.T) {
	ctx := context.Background()
	var quitPlugin1, quitPlugin2 string
	f := New(
		Pipe{
			Name: option.Some("foo"),
			When: option.Some(When{}),
			Code: option.Some(`
				export default function main(nodes){
					this.quitPlugin1="foo"
					this.quitPlugin2="bar"
				}
			`),
		},
		Plugin{
			Quit: func(ctx context.Context, x Api) error {
				require.NotNil(t, x.Flow())
				require.NotNil(t, x.Runtime())
				require.NotZero(t, x.Pipe())
				require.NotNil(t, x.This())
				err := x.Runtime().ExportTo(x.This().Get("quitPlugin1"), &quitPlugin1)
				require.NoError(t, err)
				return nil
			},
		},
		Plugin{
			Quit: func(ctx context.Context, x Api) error {
				require.NotNil(t, x.Flow())
				require.NotNil(t, x.Runtime())
				require.NotZero(t, x.Pipe())
				require.NotNil(t, x.This())
				err := x.Runtime().ExportTo(x.This().Get("quitPlugin2"), &quitPlugin2)
				require.NoError(t, err)
				return nil
			},
		},
	)
	a := []Node{{}}
	err := f.Work(ctx, a)
	require.NoError(t, err)
	require.Equal(t, "foo", quitPlugin1)
	require.Equal(t, "bar", quitPlugin2)
}
