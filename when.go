package flow

import (
	"log/slog"

	"github.com/typomaker/option"
)

type When struct {
	UUID option.Option[[]UUID]
	Kind option.Option[[]Kind]
	Hook option.Option[[]Hook]
	Live option.Option[[]Live]
}

func (it When) IsZero() bool {
	if !it.UUID.IsZero() {
		return false
	}
	if !it.Kind.IsZero() {
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
func (it When) LogAttr() slog.Attr {
	return slog.Any("when", it.LogValue())
}
func (it When) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Any("uuid", it.UUID),
		slog.Any("kind", it.Kind),
		slog.Any("hook", it.Hook),
		slog.Any("live", it.Live),
	)
}
