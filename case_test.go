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
			When: option.Some(When{UUID: option.Some([]UUID{MustUUID("eaeb5a25-21e5-47ff-a142-7a9987f2e3f0")})}),
			Then: option.Some(Then{Kind: option.Some("foo")}),
		}
		l.Info("foo", slog.Any("v", v), v.LogAttr())
		require.JSONEq(t,
			`{
				"v":{
					"when":{"uuid":["eaeb5a25-21e5-47ff-a142-7a9987f2e3f0"]},
					"then":{"kind":"foo"}
				},
				"case":{
					"when":{"uuid":["eaeb5a25-21e5-47ff-a142-7a9987f2e3f0"]},
					"then":{"kind":"foo"}
				}
			}`,
			b.String(),
		)
		b.Reset()
	})
}
