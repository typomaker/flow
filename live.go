package flow

import (
	"fmt"
	"log/slog"

	jsoniter "github.com/json-iterator/go"
	"github.com/typomaker/option"
)

type Live struct {
	Since option.Option[Time]
	Until option.Option[Time]
}

func (it Live) With(pp Live) Live {
	if !pp.Since.IsZero() {
		it.Since = pp.Since
	}
	if !pp.Until.IsZero() {
		it.Until = pp.Until
	}
	return it
}
func (it Live) LogAttr() slog.Attr {
	return slog.Any("live", it.LogValue())
}
func (it Live) LogValue() slog.Value {
	if it.IsZero() {
		return slog.GroupValue()
	}
	return slog.GroupValue(
		slog.Any("since", it.Since),
		slog.Any("until", it.Until),
	)
}

type _LiveJSON struct {
	Since jsoniter.RawMessage `json:"since,omitempty"`
	Until jsoniter.RawMessage `json:"until,omitempty"`
}

func (it Live) MarshalJSON() (b []byte, err error) {
	var js _LiveJSON
	if js.Since, err = jsoniter.Marshal(it.Since); err != nil {
		return nil, fmt.Errorf("since: %w", err)
	}
	if js.Until, err = jsoniter.Marshal(it.Until); err != nil {
		return nil, fmt.Errorf("until: %w", err)
	}
	return jsoniter.Marshal(js)
}
func (it *Live) UnmarshalJSON(b []byte) (err error) {
	var js _LiveJSON
	if err = jsoniter.Unmarshal(b, js); err != nil {
		return err
	}
	if err = jsoniter.Unmarshal(js.Since, &it.Since); err != nil {
		return fmt.Errorf("since: %w", err)
	}
	if err = jsoniter.Unmarshal(js.Until, &it.Until); err != nil {
		return fmt.Errorf("until: %w", err)
	}
	return nil
}
func (it Live) IsZero() bool {
	return it == Live{}
}
