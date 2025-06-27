package flow

import (
	"fmt"
	"log/slog"
	"slices"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/typomaker/option"
)

type Node struct {
	origin *Node
	UUID   option.Option[UUID]
	Meta   option.Option[Meta]
	Hook   option.Option[Hook]
	Live   option.Option[Live]
}

func (it Node) Equal(t Node) bool {
	switch {
	case it.UUID != t.UUID:
		return false
	case !it.Meta.GetOrZero().Equal(t.Meta.GetOrZero()):
		return false
	case !it.Hook.GetOrZero().Equal(t.Hook.GetOrZero()):
		return false
	case !it.Live.GetOrZero().Equal(t.Live.GetOrZero()):
		return false
	default:
		return true
	}
}
func (it Node) IsZero() bool {
	if it.origin != nil {
		return false
	}
	if !it.UUID.IsZero() {
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
	if it.origin != nil {
		var cp = it.origin.Copy()
		it.origin = &cp
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
func (it *Node) Origin() Node {
	if it.origin != nil {
		return *it.origin
	}
	return Node{}
}
func (it *Node) SetOrigin(o Node) {
	if o.IsZero() {
		it.origin = nil
	} else {
		it.origin = &o
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
	var origin option.Option[Node]
	if it.origin != nil {
		origin = option.Some(*it.origin)
	}
	return slog.GroupValue(
		slog.Any("uuid", it.UUID),
		slog.Any("meta", it.Meta),
		slog.Any("hook", it.Hook),
		slog.Any("live", it.Live),
		slog.Any("origin", origin),
	)
}

type _NodeJSON struct {
	UUID   jsoniter.RawMessage `json:"uuid,omitempty"`
	Kind   jsoniter.RawMessage `json:"kind,omitempty"`
	Meta   jsoniter.RawMessage `json:"meta,omitempty"`
	Hook   jsoniter.RawMessage `json:"hook,omitempty"`
	Live   jsoniter.RawMessage `json:"live,omitempty"`
	Origin jsoniter.RawMessage `json:"origin,omitempty"`
}

func (it Node) MarshalJSON() (b []byte, err error) {
	var js _NodeJSON
	if js.UUID, err = jsoniter.Marshal(it.UUID); err != nil {
		return nil, fmt.Errorf("uuid: %w", err)
	}
	if js.Meta, err = jsoniter.Marshal(it.Meta); err != nil {
		return nil, fmt.Errorf("meta: %w", err)
	}
	if js.Hook, err = jsoniter.Marshal(it.Hook); err != nil {
		return nil, fmt.Errorf("hook: %w", err)
	}
	if js.Live, err = jsoniter.Marshal(it.Live); err != nil {
		return nil, fmt.Errorf("live: %w", err)
	}
	if js.Origin, err = jsoniter.Marshal(it.origin); err != nil {
		return nil, fmt.Errorf("origin: %w", err)
	}
	return jsoniter.Marshal(js)
}
func (it *Node) UnmarshalJSON(b []byte) (err error) {
	var js _NodeJSON
	if err = jsoniter.Unmarshal(b, &js); err != nil {
		return err
	}
	if err = jsoniter.Unmarshal(js.UUID, &it.UUID); err != nil {
		return fmt.Errorf("uuid: %w", err)
	}
	if err = jsoniter.Unmarshal(js.Meta, &it.Meta); err != nil {
		return fmt.Errorf("meta: %w", err)
	}
	if err = jsoniter.Unmarshal(js.Hook, &it.Hook); err != nil {
		return fmt.Errorf("hook: %w", err)
	}
	if err = jsoniter.Unmarshal(js.Live, &it.Live); err != nil {
		return fmt.Errorf("live: %w", err)
	}
	if err = jsoniter.Unmarshal(js.Origin, it.origin); err != nil {
		return fmt.Errorf("origin: %w", err)
	}
	return nil
}

type Code = string

type Kind = string
type Name = string

type Time = time.Time
