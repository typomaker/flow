package flow

import (
	"context"
	"log/slog"
	"slices"

	jsoniter "github.com/json-iterator/go"
)

type Hook map[string]any

func (it Hook) Equal(t Hook) bool {
	return deepEqual(it, t)
}
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
func (it Hook) In(s ...Hook) Handler {
	s = append(s, it)
	var predicat = func(n Node) bool {
		if !n.Hook.IsSome() {
			return false
		}
		return slices.ContainsFunc(s, func(h Hook) bool {
			return deepContains(map[string]any(n.Hook.Get()), map[string]any(h))
		})
	}
	return func(ctx context.Context, target []Node, next Next) (err error) {
		return nextIf(target, next, predicat)
	}
}
func (it Hook) LogAttr() slog.Attr {
	return slog.Any("hook", it.LogValue())
}
func (it Hook) LogValue() slog.Value {
	var s = make([]slog.Attr, 0, len(it))
	for k, v := range it {
		if t, err := jsoniter.MarshalToString(v); err != nil {
			s = append(s, slog.String(k+"Error", err.Error()))
		} else {
			s = append(s, slog.String(k, t))
		}
	}
	return slog.GroupValue(s...)
}
func (it Hook) MarshalJSON() (b []byte, err error) {
	var js = map[string]any(it)
	return jsoniter.Marshal(js)
}
func (it *Hook) UnmarshalJSON(b []byte) (err error) {
	var js map[string]any
	if err = jsoniter.Unmarshal(b, &js); err != nil {
		return err
	}
	*it = Hook(js)
	return nil
}
