package flow

import "log/slog"

type Meta map[string]any

func (it Meta) With(pp Meta) Meta {
	if len(it) == 0 {
		return pp
	}
	if len(pp) == 0 {
		return it
	}
	for k := range pp {
		it[k] = deepWith(it[k], pp[k], k[0] == '$')
	}
	return it
}
func (it Meta) LogAttr() slog.Attr {
	return slog.Any("meta", it.LogValue())
}
func (it Meta) LogValue() slog.Value {
	return slog.AnyValue(map[string]any(it))
}
