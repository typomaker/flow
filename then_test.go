package flow

import (
	"bytes"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/typomaker/option"
)

func TestThenLog(t *testing.T) {
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
		v := Then{}
		l.Info("foo", slog.Any("v", v), v.LogAttr())
		require.JSONEq(t, `{"v":{}, "then":{}}`, b.String())
		b.Reset()
	})
	t.Run("none", func(t *testing.T) {
		l := slog.New(h)
		v := Then{
			Kind: option.None[Kind](),
			Meta: option.None[Meta](),
			Hook: option.None[Hook](),
			Live: option.None[Live](),
		}
		l.Info("foo", slog.Any("v", v), v.LogAttr())
		require.JSONEq(t, `{
				"v":{
					"kind":null,
					"meta":null,
					"hook":null,
					"live":null
				},
				"then":{
					"kind":null,
					"meta":null,
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
		v := Then{
			Kind: option.Some("foo"),
			Meta: option.Some(Meta{"b": "2"}),
			Hook: option.Some(Hook{"a": "1"}),
			Live: option.Some(Live{Since: option.Some(time.Unix(0, 0).UTC())}),
		}
		l.Info("foo", slog.Any("v", v), v.LogAttr())
		require.JSONEq(t,
			`{
				"v":{
					"kind":"foo",
					"meta":{"b":"2"},
					"hook":{"a":"1"},
					"live":{"since":"1970-01-01T00:00:00Z"}
				},
				"then":{
					"kind":"foo",
					"meta":{"b":"2"},
					"hook":{"a":"1"},
					"live":{"since":"1970-01-01T00:00:00Z"}
				}
			}`,
			b.String(),
		)
		b.Reset()
	})
}
