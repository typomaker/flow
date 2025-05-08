package flow

import (
	"log/slog"

	jsoniter "github.com/json-iterator/go"
)

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
func (it Meta) MarshalJSON() (b []byte, err error) {
	var js = map[string]any(it)
	return jsoniter.Marshal(js)
}
func (it *Meta) UnmarshalJSON(b []byte) (err error) {
	var js map[string]any
	if err = jsoniter.Unmarshal(b, &js); err != nil {
		return err
	}
	*it = Meta(js)
	return nil
}
