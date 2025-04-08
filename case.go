package flow

import (
	"log/slog"

	"github.com/typomaker/option"
)

type Case struct {
	When option.Option[When]
	Then option.Option[Then]
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
