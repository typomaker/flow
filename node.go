package flow

import (
	"log/slog"
	"slices"
	"time"

	"github.com/typomaker/option"
)

type Node struct {
	root *Node
	UUID option.Option[UUID]
	Kind option.Option[Kind]
	Meta option.Option[Meta]
	Hook option.Option[Hook]
	Live option.Option[Live]
}

func (it Node) IsZero() bool {
	if it.root != nil {
		return false
	}
	if !it.UUID.IsZero() {
		return false
	}
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
func (it Node) Copy() Node {
	if it.Meta.IsSome() {
		it.Meta = option.Some(it.Meta.Get().Copy())
	}
	if it.Hook.IsSome() {
		it.Hook = option.Some(it.Hook.Get().Copy())
	}
	if it.root != nil {
		var cp = it.root.Copy()
		it.root = &cp
	}
	return it
}
func (it Meta) Copy() Meta {
	if it == nil {
		return nil
	}
	var cp = make(Meta, len(it))
	for k := range it {
		cp[k] = deepCopy(it[k])
	}
	return cp
}
func (it Hook) Copy() Hook {
	if it == nil {
		return nil
	}
	var cp = make(Hook, len(it))
	for k := range it {
		cp[k] = deepCopy(it[k])
	}
	return cp
}
func (it *Node) Root() Node {
	if it.root != nil {
		return *it.root
	}
	return Node{}
}
func (it *Node) SetRoot(o Node) {
	if o.IsZero() {
		it.root = nil
	} else {
		it.root = &o
	}
}
func (it Node) When(w When) bool {
	switch {
	case
		w.UUID.IsNone() && !it.UUID.IsNone(),
		w.UUID.IsSome() && !it.UUID.IsSome(),
		w.UUID.IsSome() && !slices.Contains(w.UUID.Get(), it.UUID.Get()):
		return false
	case
		w.Kind.IsNone() && !it.Kind.IsNone(),
		w.Kind.IsSome() && !it.Kind.IsSome(),
		w.Kind.IsSome() && !slices.Contains(w.Kind.Get(), it.Kind.Get()):
		return false
	case
		w.Hook.IsNone() && !it.Hook.IsNone(),
		w.Hook.IsSome() && !it.Hook.IsSome(),
		w.Hook.IsSome() &&
			!slices.ContainsFunc(w.Hook.Get(), func(h Hook) bool {
				return deepHave(map[string]any(it.Hook.Get()), map[string]any(h))
			}):
		return false
	case
		w.Live.IsNone() && !it.Live.IsNone(),
		w.Live.IsSome() && !it.Live.IsSome(),
		w.Live.IsSome() &&
			!slices.ContainsFunc(w.Live.Get(), func(l Live) bool {
				switch {
				case l.Since.IsNone() && !it.Live.Get().Since.IsNone():
				case l.Since.IsSome() && !it.Live.Get().Since.IsSome():
				case l.Since.IsSome() && l.Since.Get().After(it.Live.Get().Since.Get()):
					return false
				}
				switch {
				case l.Until.IsNone() && !it.Live.Get().Until.IsNone():
				case l.Until.IsSome() && !it.Live.Get().Until.IsSome():
				case l.Until.IsSome() && l.Until.Get().Before(it.Live.Get().Until.Get()):
					return false
				}
				return true
			}):
		return false
	}
	return true
}
func (it Node) With(pp Node) Node {
	if !pp.UUID.IsZero() {
		it.UUID = pp.UUID
	}
	if !pp.Kind.IsZero() {
		it.Kind = pp.Kind
	}
	switch {
	case pp.Meta.IsSome():
		var with = it.Meta.GetOrZero().With(pp.Meta.GetOrZero())
		it.Meta = option.Some(with)
	case pp.Meta.IsNone():
		it.Meta = pp.Meta
	}
	switch {
	case pp.Hook.IsSome():
		var with = it.Hook.GetOrZero().With(pp.Hook.GetOrZero())
		it.Hook = option.Some(with)
	case pp.Hook.IsNone():
		it.Hook = pp.Hook
	}
	switch {
	case pp.Live.IsSome():
		var with = it.Live.GetOrZero().With(pp.Live.GetOrZero())
		it.Live = option.Some(with)
	case pp.Live.IsNone():
		it.Live = pp.Live
	}
	return it
}
func (it Node) LogAttr() slog.Attr {
	return slog.Any("node", it.LogValue())
}
func (it Node) LogValue() slog.Value {
	var root option.Option[Node]
	if it.root != nil {
		root = option.Some(*it.root)
	}
	return slog.GroupValue(
		slog.Any("uuid", it.UUID),
		slog.Any("kind", it.Kind),
		slog.Any("meta", it.Meta),
		slog.Any("hook", it.Hook),
		slog.Any("live", it.Live),
		slog.Any("root", root),
	)
}

type Code = string

type Kind = string
type Name = string

type Time = time.Time
