package goja

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/typomaker/flow"
)

type lazyFlowNode struct {
	rm    *goja.Runtime
	proto flow.Node
	value struct {
		UUID   goja.Value
		Kind   goja.Value
		Meta   goja.Value
		Hook   goja.Value
		Live   goja.Value
		Origin goja.Value
	}
}

var _ goja.DynamicObject = (*lazyFlowNode)(nil)

// Delete implements goja.DynamicObject.
func (it *lazyFlowNode) Delete(key string) bool {
	switch key {
	case keyUUID:
		it.value.UUID = goja.Undefined()
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
func (it *lazyFlowNode) Get(key string) (val goja.Value) {
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
	case keyMeta:
		if it.value.Meta == nil {
			switch {
			case it.proto.Meta.IsNone():
				it.value.Meta = goja.Null()
			case it.proto.Meta.IsSome():
				err = convert(it.rm, it.proto.Meta.Get(), &it.value.Meta)
			}
		}
		val = it.value.Meta
	case keyHook:
		if it.value.Hook == nil {
			switch {
			case it.proto.Hook.IsNone():
				it.value.Hook = goja.Null()
			case it.proto.Hook.IsSome():
				err = convert(it.rm, it.proto.Hook.Get(), &it.value.Hook)
			}
		}
		val = it.value.Hook
	case keyLive:
		if it.value.Live == nil {
			switch {
			case it.proto.Live.IsNone():
				it.value.Live = goja.Null()
			case it.proto.Live.IsSome():
				err = convert(it.rm, it.proto.Live.Get(), &it.value.Live)
			}
		}
		val = it.value.Live
	case keyOrigin:
		if it.value.Origin == nil {
			if !it.proto.Origin().IsZero() {
				err = convert(it.rm, it.proto.Origin(), &it.value.Origin)
			}
		}
		val = it.value.Origin
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
func (it *lazyFlowNode) Has(key string) bool {
	switch key {
	case keyUUID:
		if it.value.UUID != nil {
			return !goja.IsUndefined(it.value.UUID)
		}
		return !it.proto.UUID.IsZero()
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
	case keyOrigin:
		if it.value.Origin != nil {
			return !goja.IsUndefined(it.value.Origin)
		}
		return !it.proto.Origin().IsZero()
	default:
		return false
	}
}

// Set implements goja.DynamicObject.
func (it *lazyFlowNode) Set(key string, val goja.Value) bool {
	switch key {
	case keyUUID:
		it.value.UUID = val
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
	case keyOrigin:
		it.value.Origin = val
		return true
	default:
		return false
	}
}

//nolint:gocognit // todo: отрефакторить
func (it *lazyFlowNode) Keys() []string {
	var keys = make([]string, 0, 7)

	const uuidKey = keyUUID
	if it.value.UUID != nil {
		if !goja.IsUndefined(it.value.UUID) {
			keys = append(keys, uuidKey)
		}
	} else if !it.proto.UUID.IsZero() {
		keys = append(keys, uuidKey)
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

	if it.value.Origin != nil {
		if !goja.IsUndefined(it.value.Origin) {
			keys = append(keys, keyOrigin)
		}
	} else if !it.proto.Origin().IsZero() {
		keys = append(keys, keyOrigin)
	}
	return keys
}
