package flow

import (
	"fmt"

	"github.com/dop251/goja"
)

type lazyNode struct {
	rm    *goja.Runtime
	proto Node
	value struct {
		UUID goja.Value
		Kind goja.Value
		Meta goja.Value
		Hook goja.Value
		Live goja.Value
		Root goja.Value
	}
}

var _ goja.DynamicObject = (*lazyNode)(nil)

// Delete implements goja.DynamicObject.
func (it *lazyNode) Delete(key string) bool {
	switch key {
	case keyUUID:
		it.value.UUID = goja.Undefined()
		return true
	case keyKind:
		it.value.Kind = goja.Undefined()
		return true
	case keyMeta:
		it.value.Meta = goja.Undefined()
		return true
	case keyHook:
		it.value.Hook = goja.Undefined()
		return true
	case keyLive:
		it.value.Live = goja.Undefined()
		return true
	default:
		return false
	}
}

//nolint:gocognit,funlen // todo: отрефакторить
func (it *lazyNode) Get(key string) (val goja.Value) {
	var err error
	switch key {
	case keyUUID:
		if it.value.UUID == nil {
			switch {
			case it.proto.UUID.IsNone():
				it.value.UUID = goja.Null()
			case it.proto.UUID.IsSome():
				it.value.UUID = it.rm.ToValue(it.proto.UUID.Get().String())
			}
		}
		val = it.value.UUID
	case keyKind:
		if it.value.Kind == nil {
			switch {
			case it.proto.Kind.IsNone():
				it.value.Kind = goja.Null()
			case it.proto.Kind.IsSome():
				it.value.Kind = it.rm.ToValue(it.proto.Kind.Get())
			}
		}
		val = it.value.Kind
	case keyMeta:
		if it.value.Meta == nil {
			switch {
			case it.proto.Meta.IsNone():
				it.value.Meta = goja.Null()
			case it.proto.Meta.IsSome():
				err = Convert(it.rm, it.proto.Meta.Get(), &it.value.Meta)
			}
		}
		val = it.value.Meta
	case keyHook:
		if it.value.Hook == nil {
			switch {
			case it.proto.Hook.IsNone():
				it.value.Hook = goja.Null()
			case it.proto.Hook.IsSome():
				err = Convert(it.rm, it.proto.Hook.Get(), &it.value.Hook)
			}
		}
		val = it.value.Hook
	case keyLive:
		if it.value.Live == nil {
			switch {
			case it.proto.Live.IsNone():
				it.value.Live = goja.Null()
			case it.proto.Live.IsSome():
				err = Convert(it.rm, it.proto.Live.Get(), &it.value.Live)
			}
		}
		val = it.value.Live
	case keyRoot:
		if it.value.Root == nil {
			if !it.proto.Root().IsZero() {
				err = Convert(it.rm, it.proto.Root(), &it.value.Root)
			}
		}
		val = it.value.Root
	}
	if err != nil {
		err = fmt.Errorf("flowItem.get: %w", err)
		panic(it.rm.NewGoError(err))
	}
	if val == nil {
		return goja.Undefined()
	}
	return val
}

// Has implements goja.DynamicObject.
func (it *lazyNode) Has(key string) bool {
	switch key {
	case keyUUID:
		if it.value.UUID != nil {
			return !goja.IsUndefined(it.value.UUID)
		}
		return !it.proto.UUID.IsZero()
	case keyKind:
		if it.value.Kind != nil {
			return !goja.IsUndefined(it.value.Kind)
		}
		return !it.proto.Kind.IsZero()
	case keyMeta:
		if it.value.Meta != nil {
			return !goja.IsUndefined(it.value.Meta)
		}
		return !it.proto.Meta.IsZero()
	case keyHook:
		if it.value.Hook != nil {
			return !goja.IsUndefined(it.value.Hook)
		}
		return !it.proto.Hook.IsZero()
	case keyLive:
		if it.value.Live != nil {
			return !goja.IsUndefined(it.value.Live)
		}
		return !it.proto.Live.IsZero()
	case keyRoot:
		if it.value.Root != nil {
			return !goja.IsUndefined(it.value.Root)
		}
		return !it.proto.Root().IsZero()
	default:
		return false
	}
}

// Set implements goja.DynamicObject.
func (it *lazyNode) Set(key string, val goja.Value) bool {
	switch key {
	case keyUUID:
		it.value.UUID = val
		return true
	case keyKind:
		it.value.Kind = val
		return true
	case keyMeta:
		it.value.Meta = val
		return true
	case keyHook:
		it.value.Hook = val
		return true
	case keyLive:
		it.value.Live = val
		return true
	case keyRoot:
		it.value.Root = val
		return true
	default:
		return false
	}
}

//nolint:gocognit // todo: отрефакторить
func (it *lazyNode) Keys() []string {
	var keys = make([]string, 0, 7)

	const uuidKey = keyUUID
	if it.value.UUID != nil {
		if !goja.IsUndefined(it.value.UUID) {
			keys = append(keys, uuidKey)
		}
	} else if !it.proto.UUID.IsZero() {
		keys = append(keys, uuidKey)
	}

	if it.value.Kind != nil {
		if !goja.IsUndefined(it.value.Kind) {
			keys = append(keys, keyKind)
		}
	} else if !it.proto.Kind.IsZero() {
		keys = append(keys, keyKind)
	}

	if it.value.Meta != nil {
		if !goja.IsUndefined(it.value.Meta) {
			keys = append(keys, keyMeta)
		}
	} else if !it.proto.Meta.IsZero() {
		keys = append(keys, keyMeta)
	}

	if it.value.Hook != nil {
		if !goja.IsUndefined(it.value.Hook) {
			keys = append(keys, keyHook)
		}
	} else if !it.proto.Hook.IsZero() {
		keys = append(keys, keyHook)
	}

	if it.value.Live != nil {
		if !goja.IsUndefined(it.value.Live) {
			keys = append(keys, keyLive)
		}
	} else if !it.proto.Live.IsZero() {
		keys = append(keys, keyLive)
	}

	if it.value.Root != nil {
		if !goja.IsUndefined(it.value.Root) {
			keys = append(keys, keyRoot)
		}
	} else if !it.proto.Root().IsZero() {
		keys = append(keys, keyRoot)
	}
	return keys
}
