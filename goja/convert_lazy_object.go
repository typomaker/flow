package goja

import (
	"github.com/dop251/goja"
)

type lazyObject struct {
	rm    *goja.Runtime
	proto map[string]any
	value map[string]goja.Value
}

var _ goja.DynamicObject = (*lazyObject)(nil)

// Delete implements goja.DynamicObject.
func (it *lazyObject) Delete(key string) bool {
	if it.value == nil {
		it.value = make(map[string]goja.Value, len(it.proto))
	}
	it.value[key] = goja.Undefined()
	return true
}

// Get implements goja.DynamicObject.
func (it *lazyObject) Get(key string) goja.Value {
	var jsAny, ok = it.value[key]
	if !ok {
		var goAny any
		if goAny, ok = it.proto[key]; !ok {
			return goja.Undefined()
		}
		var err error
		if err = convert(it.rm, goAny, &jsAny); err != nil {
			panic(it.rm.NewGoError(err))
		}
		if it.value == nil {
			it.value = make(map[string]goja.Value, len(it.proto))
		}
		it.value[key] = jsAny
	}
	return jsAny
}

// Has implements goja.DynamicObject.
func (it *lazyObject) Has(key string) (ok bool) {
	if _, ok = it.value[key]; !ok {
		_, ok = it.proto[key]
	}
	return ok
}

// Set implements goja.DynamicObject.
func (it *lazyObject) Set(key string, jsValue goja.Value) bool {
	if it.value == nil {
		it.value = make(map[string]goja.Value, len(it.proto))
	}
	it.value[key] = jsValue
	return true
}

// Keys implements goja.DynamicObject.
func (it *lazyObject) Keys() []string {
	var lng int
	if lng = len(it.proto); lng < len(it.value) {
		lng = len(it.value)
	}
	var keys = make([]string, 0, lng)
	for key, val := range it.value {
		if goja.IsUndefined(val) {
			continue
		}
		keys = append(keys, key)
	}
	for key := range it.proto {
		if _, ok := it.value[key]; ok {
			continue
		}
		keys = append(keys, key)
	}
	return keys
}
