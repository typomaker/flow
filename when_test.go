package flow

import (
	"bytes"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/typomaker/option"
)

func TestWhenLog(t *testing.T) {
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
		v := When{}
		l.Info("foo", slog.Any("v", v), v.LogAttr())
		require.JSONEq(t, `{"v":{}, "when":{}}`, b.String())
		b.Reset()
	})
	t.Run("none", func(t *testing.T) {
		l := slog.New(h)
		v := When{
			UUID: option.None[[]UUID](),
			Hook: option.None[[]Hook](),
			Live: option.None[[]Live](),
		}
		l.Info("foo", slog.Any("v", v), v.LogAttr())
		require.JSONEq(t, `{
				"v":{
					"uuid":null,
					"hook":null,
					"live":null
				},
				"when":{
					"uuid":null,
					"hook":null,
					"live":null
				}
			}`,
			b.String(),
		)
		b.Reset()
	})
	t.Run("some", func(t *testing.T) {
		l := slog.New(h)
		v := When{
			UUID: option.Some([]UUID{MustUUID("1f3219bb-9577-4504-90da-305772121e18")}),
			Hook: option.Some([]Hook{{"a": "1"}}),
			Live: option.Some([]Live{{Since: option.Some(time.Unix(0, 0).UTC())}}),
		}
		l.Info("foo", slog.Any("v", v), v.LogAttr())
		require.JSONEq(t,
			`{
				"v":{
					"uuid":["1f3219bb-9577-4504-90da-305772121e18"],
					"hook":[{"a":"1"}],
					"live":[{"since":"1970-01-01T00:00:00Z"}]
				},
				"when":{
					"uuid":["1f3219bb-9577-4504-90da-305772121e18"],
					"hook":[{"a":"1"}],
					"live":[{"since":"1970-01-01T00:00:00Z"}]
				}
			}`,
			b.String(),
		)
		b.Reset()
	})
}
