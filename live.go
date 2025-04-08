package flow

import (
	"encoding/json"
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
func (it Live) MarshalJSON() (b []byte, err error) {
	var j struct {
		Since json.RawMessage `json:"since,omitempty"`
		Until json.RawMessage `json:"until,omitempty"`
	}
	if j.Since, err = it.Since.MarshalJSON(); err != nil {
		return nil, fmt.Errorf("since: %w", err)
	}
	if j.Until, err = it.Until.MarshalJSON(); err != nil {
		return nil, fmt.Errorf("until: %w", err)
	}
	return jsoniter.Marshal(j)
}
func (it Live) IsZero() bool {
	return it == Live{}
}
