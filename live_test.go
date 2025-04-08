package flow

import (
	"bytes"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/typomaker/option"
)

func TestLiveLog(t *testing.T) {
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
		v := Live{}
		l.Info("foo", slog.Any("v", v), v.LogAttr())
		require.JSONEq(t, `{}`, b.String())
		b.Reset()
	})
	t.Run("none", func(t *testing.T) {
		l := slog.New(h)
		v := Live{
			Since: option.None[time.Time](),
			Until: option.None[time.Time](),
		}
		l.Info("foo", slog.Any("v", v), v.LogAttr())
		require.JSONEq(t, `{"v":{"since":null,"until":null},"live":{"since":null,"until":null}}`, b.String())
		b.Reset()
	})

	t.Run("some", func(t *testing.T) {
		l := slog.New(h)
		v := Live{
			Since: option.Some(time.Unix(0, 0).UTC()),
			Until: option.Some(time.Unix(86400, 0).UTC()),
		}
		l.Info("foo", slog.Any("v", v), v.LogAttr())
		require.JSONEq(t,
			`{
				"v":{
					"since": "1970-01-01T00:00:00Z",
					"until": "1970-01-02T00:00:00Z"
				},
				"live":{
					"since": "1970-01-01T00:00:00Z",
					"until": "1970-01-02T00:00:00Z"
				}
			}`,
			b.String(),
		)
		b.Reset()
	})

}
