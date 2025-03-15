package flow

import (
	"log/slog"
	"slices"
	"time"

	"github.com/google/uuid"
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

type Case struct {
	When When
	Then Then
}
type When struct {
	UUID option.Option[[]UUID]
	Kind option.Option[[]Kind]
	Hook option.Option[[]Hook]
	Live option.Option[[]Live]
}

func (it When) IsZero() bool {
	if !it.UUID.IsZero() {
		return false
	}
	if !it.Kind.IsZero() {
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
func (it When) LogValue() slog.Value {
	var attrs []slog.Attr
	defer reuseSliceSlogAttr(&attrs)()

	switch {
	case it.UUID.IsNone():
		attrs = append(attrs, slog.Any("uuid", nil))
	case it.UUID.IsSome():
		attrs = append(attrs, slog.Any("uuid", it.UUID.Get()))
	}
	switch {
	case it.Kind.IsNone():
		attrs = append(attrs, slog.Any("kind", nil))
	case it.Kind.IsSome():
		attrs = append(attrs, slog.Any("kind", it.Kind.Get()))
	}
	switch {
	case it.Hook.IsNone():
		attrs = append(attrs, slog.Any("hook", nil))
	case it.Hook.IsSome():
		attrs = append(attrs, slog.Any("hook", it.Hook.Get()))
	}
	switch {
	case it.Live.IsNone():
		attrs = append(attrs, slog.Any("live", nil))
	case it.Live.IsSome():
		attrs = append(attrs, slog.Any("live", it.Live.Get()))
	}
	return slog.GroupValue(attrs...)
}

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

type Live struct {
	Since option.Option[Time]
	Until option.Option[Time]
}

func (it Live) LogValue() slog.Value {
	var attrs []slog.Attr
	defer reuseSliceSlogAttr(&attrs)()

	switch {
	case it.Since.IsNone():
		attrs = append(attrs, slog.Any("since", nil))
	case it.Since.IsSome():
		attrs = append(attrs, slog.Any("since", it.Since.Get()))
	}
	switch {
	case it.Until.IsNone():
		attrs = append(attrs, slog.Any("until", nil))
	case it.Until.IsSome():
		attrs = append(attrs, slog.Any("until", it.Until.Get()))
	}
	return slog.GroupValue(attrs...)
}

type Code = string
type UUID [16]byte

func NewUUID() UUID {
	var u = UUID(uuid.Must(uuid.NewRandom()))
	return u
}
func MustUUID(s string) UUID {
	var u, err = ParseUUID(s)
	if err != nil {
		panic(err)
	}
	return u
}
func ParseUUID(s string) (u UUID, err error) {
	var x uuid.UUID
	if x, err = uuid.Parse(s); err != nil {
		return u, err
	}
	u = UUID(x)
	return u, nil
}
func (it UUID) GoString() string {
	return "\"" + it.String() + "\""
}

type Kind = string
type Name = string
type Meta map[string]any
type Hook map[string]any
type Time = time.Time

func (it Pipe) String() string {
	return it.Name.Get()
}
func (it UUID) String() string {
	return uuid.UUID(it).String()
}
