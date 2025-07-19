package flow

import (
	"log/slog"
)

type Case struct {
	When When
	Then Then
}

func (it Case) Equal(t Case) bool {
	switch {
	case !it.When.Equal(t.When):
		return false
	case !it.Then.Equal(t.Then):
		return false
	default:
		return true
	}
}
func (it Case) LogAttr() slog.Attr {
	return slog.Any("case", it.LogValue())
}
func (it Case) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Any("when", it.When),
		slog.Any("then", it.Then),
	)
}
