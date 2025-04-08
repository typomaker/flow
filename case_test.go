package flow

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typomaker/option"
)

func TestCaseLog(t *testing.T) {
	b := bytes.Buffer{}
	h := slog.NewJSONHandler(&b, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case "time", "level", "msg":
				return slog.Attr{}
			}
			return a
		},
	})
	t.Run("zero", func(t *testing.T) {
		l := slog.New(h)
		v := Case{}
		l.Info("foo", slog.Any("v", v), v.LogAttr())
		require.JSONEq(t, `{"v":{}, "case":{}}`, b.String())
		b.Reset()
	})
	t.Run("none", func(t *testing.T) {
		l := slog.New(h)
		v := Case{
			When: option.None[When](),
			Then: option.None[Then](),
		}
		l.Info("foo", slog.Any("v", v), v.LogAttr())
		require.JSONEq(t, `{
				"v":{
					"when":null,
					"then":null
				},
				"case":{
					"when":null,
					"then":null
				}
			}`,
			b.String(),
		)
		b.Reset()
	})
	t.Run("some", func(t *testing.T) {
		l := slog.New(h)
		v := Case{
			When: option.Some(When{Kind: option.Some([]Kind{"foo"})}),
			Then: option.Some(Then{Kind: option.Some("foo")}),
		}
		l.Info("foo", slog.Any("v", v), v.LogAttr())
		require.JSONEq(t,
			`{
				"v":{
					"when":{"kind":["foo"]},
					"then":{"kind":"foo"}
				},
				"case":{
					"when":{"kind":["foo"]},
					"then":{"kind":"foo"}
				}
			}`,
			b.String(),
		)
		b.Reset()
	})
}
