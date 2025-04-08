package flow

import "log/slog"

type Hook map[string]any

func (it Hook) With(pp Hook) Hook {
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
func (it Hook) LogAttr() slog.Attr {
	return slog.Any("hook", it.LogValue())
}
func (it Hook) LogValue() slog.Value {
	return slog.AnyValue(map[string]any(it))
}
