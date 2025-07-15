package flow

import (
	"fmt"
	"log/slog"

	jsoniter "github.com/json-iterator/go"
	"github.com/typomaker/option"
)

type Pipe struct {
	Name option.Option[Name]
	When option.Option[When]
	Code option.Option[Code]
	Next option.Option[[]Name]
}

func (it Pipe) IsZero() bool {
	if !it.Name.IsZero() {
		return false
	}
	if !it.When.IsZero() {
		return false
	}
	if !it.Code.IsZero() {
		return false
	}
	if !it.Next.IsZero() {
		return false
	}
	return true
}
func (it Pipe) LogAttr() slog.Attr {
	return slog.Any("pipe", it.LogValue())
}
func (it Pipe) LogValue() slog.Value {
	var attrs []slog.Attr
	defer getSliceSlogAttr(&attrs)()

	switch {
	case it.Name.IsNone():
		attrs = append(attrs, slog.Any("name", nil))
	case it.Name.IsSome():
		attrs = append(attrs, slog.Any("name", it.Name.Get()))
	}
	switch {
	case it.When.IsNone():
		attrs = append(attrs, slog.Any("when", nil))
	case it.When.IsSome():
		attrs = append(attrs, slog.Any("when", it.When.Get()))
	}
	switch {
	case it.Code.IsNone():
		attrs = append(attrs, slog.Any("code", nil))
	case it.Code.IsSome():
		attrs = append(attrs, slog.Any("code", it.Code.Get()))
	}
	switch {
	case it.Next.IsNone():
		attrs = append(attrs, slog.Any("next", nil))
	case it.Next.IsSome():
		attrs = append(attrs, slog.Any("next", it.Next.Get()))
	}
	return slog.GroupValue(attrs...)
}

type _PipeJSON struct {
	Name jsoniter.RawMessage `json:"name,omitempty"`
	When jsoniter.RawMessage `json:"when,omitempty"`
	Code jsoniter.RawMessage `json:"code,omitempty"`
	Next jsoniter.RawMessage `json:"next,omitempty"`
}

func (it Pipe) MarshalJSON() (b []byte, err error) {
	var js _PipeJSON
	if js.Name, err = jsoniter.Marshal(it.Name); err != nil {
		return nil, fmt.Errorf("name: %w", err)
	}
	if js.When, err = jsoniter.Marshal(it.When); err != nil {
		return nil, fmt.Errorf("when: %w", err)
	}
	if js.Code, err = jsoniter.Marshal(it.Code); err != nil {
		return nil, fmt.Errorf("code: %w", err)
	}
	if js.Next, err = jsoniter.Marshal(it.Next); err != nil {
		return nil, fmt.Errorf("next: %w", err)
	}
	return jsoniter.Marshal(js)
}
func (it *Pipe) UnmarshalJSON(b []byte) (err error) {
	var js _PipeJSON
	if err = jsoniter.Unmarshal(b, &js); err != nil {
		return err
	}
	if err = jsoniter.Unmarshal(js.Name, &it.Name); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	if err = jsoniter.Unmarshal(js.When, &it.When); err != nil {
		return fmt.Errorf("when: %w", err)
	}
	if err = jsoniter.Unmarshal(js.Code, &it.Code); err != nil {
		return fmt.Errorf("code: %w", err)
	}
	if err = jsoniter.Unmarshal(js.Next, &it.Next); err != nil {
		return fmt.Errorf("next: %w", err)
	}
	return nil
}
func (it Pipe) String() string {
	return it.Name.Get()
}
