package flow

import (
	"bytes"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/typomaker/option"
)

func TestNodeLog(t *testing.T) {
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
		v := Node{}
		l.Info("foo", slog.Any("v", v), v.LogAttr())
		require.JSONEq(t, `{"v":{}, "node":{}}`, b.String())
		b.Reset()
	})
	t.Run("none", func(t *testing.T) {
		l := slog.New(h)
		v := Node{
			UUID: option.None[UUID](),
			Kind: option.None[Kind](),
			Meta: option.None[Meta](),
			Hook: option.None[Hook](),
			Live: option.None[Live](),
		}
		v.SetOrigin(Node{})
		l.Info("foo", slog.Any("v", v), v.LogAttr())
		require.JSONEq(t, `{
				"v":{
					"uuid":null,
					"kind":null,
					"meta":null,
					"hook":null,
					"live":null
				},
				"node":{
					"uuid":null,
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
		v := Node{
			UUID: option.Some(MustUUID("85432856-ba6c-46d3-9fcf-05650bfd5814")),
			Kind: option.Some("foo"),
			Meta: option.Some(Meta{"b": "2"}),
			Hook: option.Some(Hook{"a": "1"}),
			Live: option.Some(Live{Since: option.Some(time.Unix(0, 0).UTC())}),
		}
		v.SetOrigin(v)

		l.Info("foo", slog.Any("v", v), v.LogAttr())
		require.JSONEq(t,
			`{
				"v":{
					"uuid":"85432856-ba6c-46d3-9fcf-05650bfd5814",
					"kind":"foo",
					"meta":{"b":"2"},
					"hook":{"a":"1"},
					"live":{"since":"1970-01-01T00:00:00Z"},
					"origin":{
						"uuid":"85432856-ba6c-46d3-9fcf-05650bfd5814",
						"kind":"foo",
						"meta":{"b":"2"},
						"hook":{"a":"1"},
						"live":{"since":"1970-01-01T00:00:00Z"}
					}
				},
				"node":{
					"uuid":"85432856-ba6c-46d3-9fcf-05650bfd5814",
					"kind":"foo",
					"meta":{"b":"2"},
					"hook":{"a":"1"},
					"live":{"since":"1970-01-01T00:00:00Z"},
					"origin":{
						"uuid":"85432856-ba6c-46d3-9fcf-05650bfd5814",
						"kind":"foo",
						"meta":{"b":"2"},
						"hook":{"a":"1"},
						"live":{"since":"1970-01-01T00:00:00Z"}
					}
				}
			}`,
			b.String(),
		)
		b.Reset()
	})
}
