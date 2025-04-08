package flow

import (
	"log/slog"

	"github.com/typomaker/option"
)

type Then struct {
	Kind option.Option[Kind]
	Meta option.Option[Meta]
	Hook option.Option[Hook]
	Live option.Option[Live]
}

func (it Then) IsZero() bool {
	if !it.Kind.IsZero() {
		return false
	}
	if !it.Meta.IsZero() {
		return false
	}
	if !it.Hook.IsZero() {
		return false
	}
	if !it.Live.IsZero() {
		return false
	}
	return true
}

func (it Then) LogAttr() slog.Attr {
	return slog.Any("then", it.LogValue())
}
func (it Then) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Any("kind", it.Kind),
		slog.Any("meta", it.Meta),
		slog.Any("hook", it.Hook),
		slog.Any("live", it.Live),
	)
}
