package flow

import (
	"fmt"
	"log/slog"

	jsoniter "github.com/json-iterator/go"
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

type _ThenJSON struct {
	Kind jsoniter.RawMessage `json:"kind,omitempty"`
	Hook jsoniter.RawMessage `json:"hook,omitempty"`
	Live jsoniter.RawMessage `json:"live,omitempty"`
}

func (it Then) MarshalJSON() (b []byte, err error) {
	var js _ThenJSON
	if js.Kind, err = jsoniter.Marshal(it.Kind); err != nil {
		return nil, fmt.Errorf("kind: %w", err)
	}
	if js.Hook, err = jsoniter.Marshal(it.Hook); err != nil {
		return nil, fmt.Errorf("hook: %w", err)
	}
	if js.Live, err = jsoniter.Marshal(it.Live); err != nil {
		return nil, fmt.Errorf("live: %w", err)
	}
	return jsoniter.Marshal(js)
}
func (it *Then) UnmarshalJSON(b []byte) (err error) {
	var js _ThenJSON
	if err = jsoniter.Unmarshal(b, &js); err != nil {
		return err
	}
	if err = jsoniter.Unmarshal(js.Kind, &it.Kind); err != nil {
		return fmt.Errorf("kind: %w", err)
	}
	if err = jsoniter.Unmarshal(js.Hook, &it.Hook); err != nil {
		return fmt.Errorf("hook: %w", err)
	}
	if err = jsoniter.Unmarshal(js.Live, &it.Live); err != nil {
		return fmt.Errorf("live: %w", err)
	}
	return nil
}
