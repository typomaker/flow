package flow

import (
	"log/slog"

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
	defer reuseSliceSlogAttr(&attrs)()

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
func (it Pipe) String() string {
	return it.Name.Get()
}
