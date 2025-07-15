package gojaflow

import (
	"slices"

	"github.com/dop251/goja"
	"github.com/typomaker/flow"
)

type lazyFlowList struct {
	rm    *goja.Runtime
	proto []flow.Node
	value []goja.Value
}

var _ goja.DynamicArray = (*lazyFlowList)(nil)

// Get implements goja.DynamicArray.

// Get implements goja.DynamicArray.
func (it *lazyFlowList) Get(idx int) goja.Value {
	var jsVal goja.Value
	if it.value != nil {
		if idx >= len(it.value) {
			return goja.Undefined()
		}
		if jsVal = it.value[idx]; jsVal != nil {
			return jsVal
		}
	}
	if idx >= len(it.proto) {
		return goja.Undefined()
	}
	var err error
	if err = convert(it.rm, it.proto[idx], &jsVal); err != nil {
		panic(it.rm.NewGoError(err))
	}
	if it.value == nil {
		it.value = make([]goja.Value, len(it.proto))
	}
	it.value[idx] = jsVal
	return jsVal
}

// Len implements goja.DynamicArray.
func (it *lazyFlowList) Len() int {
	if it.value != nil {
		return len(it.value)
	}
	return len(it.proto)
}

// SetLen implements goja.DynamicArray.
func (it *lazyFlowList) SetLen(size int) bool {
	if it.value == nil {
		it.value = make([]goja.Value, len(it.proto))
	}
	if size > len(it.value) {
		it.value = slices.Grow(it.value, size)
	}
	var valueLen = len(it.value)
	switch {
	case valueLen < size:
		it.value = it.value[:size]
		for i := valueLen; i < size; i++ {
			it.value[i] = goja.Undefined()
		}
	case size < valueLen:
		for i := size; i < valueLen; i++ {
			it.value[i] = goja.Undefined()
		}
		it.value = it.value[:size]
	}
	return true
}

// Set implements goja.DynamicArray.
func (it *lazyFlowList) Set(idx int, val goja.Value) bool {
	if it.value == nil {
		it.value = make([]goja.Value, len(it.proto))
	}
	var valueLen = idx + 1
	if valueLen > len(it.value) {
		it.value = slices.Grow(it.value, valueLen)[:valueLen]
	}
	it.value[idx] = val
	return true
}
