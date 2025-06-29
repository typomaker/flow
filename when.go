package flow

import (
	"fmt"
	"log/slog"
	"slices"

	jsoniter "github.com/json-iterator/go"
	"github.com/typomaker/option"
)

type When struct {
	UUID option.Option[[]UUID] `json:"uuid,omitempty"`
	Hook option.Option[[]Hook] `json:"hook,omitempty"`
	Live option.Option[[]Live] `json:"live,omitempty"`
}

func (it When) Equal(t When) bool {
	switch {
	case !slices.Equal(it.UUID.GetOrZero(), t.UUID.GetOrZero()):
		return false
	case !slices.EqualFunc(it.Hook.GetOrZero(), t.Hook.GetOrZero(), Hook.Equal):
		return false
	case !slices.EqualFunc(it.Live.GetOrZero(), t.Live.GetOrZero(), Live.Equal):
		return false
	default:
		return true
	}
}
func (it When) IsZero() bool {
	if !it.UUID.IsZero() {
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
		slog.Any("hook", it.Hook),
		slog.Any("live", it.Live),
	)
}

type _WhenJSON struct {
	UUID jsoniter.RawMessage `json:"uuid,omitempty"`
	Hook jsoniter.RawMessage `json:"hook,omitempty"`
	Live jsoniter.RawMessage `json:"live,omitempty"`
}

func (it When) MarshalJSON() (b []byte, err error) {
	var js _WhenJSON
	if js.UUID, err = jsoniter.Marshal(it.UUID); err != nil {
		return nil, fmt.Errorf("uuid: %w", err)
	}
	if js.Hook, err = jsoniter.Marshal(it.Hook); err != nil {
		return nil, fmt.Errorf("hook: %w", err)
	}
	if js.Live, err = jsoniter.Marshal(it.Live); err != nil {
		return nil, fmt.Errorf("live: %w", err)
	}
	return jsoniter.Marshal(js)
}
func (it *When) UnmarshalJSON(b []byte) (err error) {
	var js _WhenJSON
	if err = jsoniter.Unmarshal(b, &js); err != nil {
		return err
	}
	if err = jsoniter.Unmarshal(js.UUID, &it.UUID); err != nil {
		return fmt.Errorf("uuid: %w", err)
	}
	if err = jsoniter.Unmarshal(js.Hook, &it.Hook); err != nil {
		return fmt.Errorf("hook: %w", err)
	}
	if err = jsoniter.Unmarshal(js.Live, &it.Live); err != nil {
		return fmt.Errorf("live: %w", err)
	}
	return nil
}
