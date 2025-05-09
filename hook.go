package flow

import (
	"log/slog"

	jsoniter "github.com/json-iterator/go"
)

type Hook map[string]any

func (it Hook) Equal(t Hook) bool {
	return equal(it, t)
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
func (it Hook) LogAttr() slog.Attr {
	return slog.Any("hook", it.LogValue())
}
func (it Hook) LogValue() slog.Value {
	return slog.AnyValue(map[string]any(it))
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
