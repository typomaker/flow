package goja

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/typomaker/flow"
)

type lazyFlowLive struct {
	rm    *goja.Runtime
	proto flow.Live
	value struct {
		Since goja.Value
		Until goja.Value
	}
}

var _ goja.DynamicObject = (*lazyFlowLive)(nil)

// Delete implements goja.DynamicObject.
func (it *lazyFlowLive) Delete(key string) bool {
	switch key {
	case keyLiveSince:
		it.value.Since = goja.Undefined()
	case keyLiveUntil:
		it.value.Until = goja.Undefined()
	}
	return false
}

// Get implements goja.DynamicObject.
func (it *lazyFlowLive) Get(key string) goja.Value {
	const jsMsPrec = 1e6
	switch key {
	case keyLiveSince:
		if it.value.Since == nil {
			switch {
			case it.proto.Since.IsZero():
				it.value.Since = goja.Undefined()
			case it.proto.Since.IsNone():
				it.value.Since = goja.Null()
			default:
				var jsSince, err = it.rm.New(
					it.rm.Get("Date").ToObject(it.rm),
					it.rm.ToValue(it.proto.Since.Get().UnixNano()/jsMsPrec),
				)
				if err != nil {
					panic(it.rm.NewGoError(fmt.Errorf("since %w", err)))
				}
				it.value.Since = it.rm.ToValue(jsSince)
			}
		}
		return it.value.Since
	case keyLiveUntil:
		if it.value.Until == nil {
			switch {
			case it.proto.Until.IsZero():
				it.value.Until = goja.Undefined()
			case it.proto.Until.IsNone():
				it.value.Until = goja.Null()
			default:
				var jsUntil, err = it.rm.New(
					it.rm.Get("Date").ToObject(it.rm),
					it.rm.ToValue(it.proto.Until.Get().UnixNano()/jsMsPrec),
				)
				if err != nil {
					panic(it.rm.NewGoError(fmt.Errorf("until %w", err)))
				}
				it.value.Until = it.rm.ToValue(jsUntil)
			}
		}
		return it.value.Until
	default:
		return nil
	}
}

// Has implements goja.DynamicObject.
func (it *lazyFlowLive) Has(key string) bool {
	switch key {
	case keyLiveSince:
		if it.value.Since != nil {
			return !goja.IsUndefined(it.value.Since)
		}
		return !it.proto.Since.IsZero()
	case keyLiveUntil:
		if it.value.Until != nil {
			return !goja.IsUndefined(it.value.Until)
		}
		return !it.proto.Until.IsZero()
	default:
		return false
	}
}

// Set implements goja.DynamicObject.
func (it *lazyFlowLive) Set(key string, val goja.Value) bool {
	switch key {
	case keyLiveSince:
		it.value.Since = val
		return true
	case keyLiveUntil:
		it.value.Until = val
		return true
	default:
		return false
	}
}

// Keys implements goja.DynamicObject.
func (it *lazyFlowLive) Keys() []string {
	var keys = make([]string, 0, 2)
	if it.value.Since != nil {
		if !goja.IsUndefined(it.value.Since) {
			keys = append(keys, keyLiveSince)
		}
	} else if !it.proto.Since.IsZero() {
		keys = append(keys, keyLiveSince)
	}

	if it.value.Until != nil {
		if !goja.IsUndefined(it.value.Until) {
			keys = append(keys, keyLiveUntil)
		}
	} else if !it.proto.Until.IsZero() {
		keys = append(keys, keyLiveUntil)
	}

	return keys
}
