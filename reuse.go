package flow

import (
	"log/slog"
	"strings"
	"sync"
)

var reuse struct {
	SliceSliceNode sync.Pool
	SliceSlicePipe sync.Pool
	SlicePipe      sync.Pool
	MapUUIDStruct  sync.Pool
	StringsBuilder sync.Pool
	SliceSlogAttr  sync.Pool
	SliceError     sync.Pool
}

func reuseSliceSliceNode(v *[][]Node) (closer func()) {
	if *v != nil {
		return
	}
	var x, _ = reuse.SliceSliceNode.Get().(*[][]Node)
	if x == nil {
		var a = make([][]Node, 0, 8)
		x = &a
	}
	*v = *x

	return func() {
		clear(*v)
		*v = (*v)[:0]
		reuse.SliceSliceNode.Put(v)
	}
}
func reuseSliceSlicePipe(v *[][]Pipe) (closer func()) {
	if *v != nil {
		return
	}
	var x, _ = reuse.SliceSlicePipe.Get().(*[][]Pipe)
	if x == nil {
		var a = make([][]Pipe, 0, 8)
		x = &a
	}
	*v = *x

	return func() {
		clear(*v)
		*v = (*v)[:0]
		reuse.SliceSlicePipe.Put(v)
	}
}
func reuseSlicePipe(v *[]Pipe) (closer func()) {
	if *v != nil {
		return
	}
	var x, _ = reuse.SlicePipe.Get().(*[]Pipe)
	if x == nil {
		var a = make([]Pipe, 0, 8)
		x = &a
	}
	*v = *x

	return func() {
		clear(*v)
		*v = (*v)[:0]
		reuse.SlicePipe.Put(v)
	}
}
func reuseMapUUIDSrtuct(v *map[UUID]struct{}) (closer func()) {
	if *v != nil {
		return
	}
	var x, _ = reuse.MapUUIDStruct.Get().(*map[UUID]struct{})
	if x == nil {
		var a = make(map[UUID]struct{}, 8)
		x = &a
	}
	*v = *x

	return func() {
		clear(*v)
		reuse.MapUUIDStruct.Put(v)
	}
}
func reuseStringBuilder(v *strings.Builder) (closer func()) {
	if x, _ := reuse.StringsBuilder.Get().(*strings.Builder); x != nil {
		*v = *x
	} else {
		*v = strings.Builder{}
	}
	return func() {
		v.Reset()
		reuse.StringsBuilder.Put(v)
	}
}
func reuseSliceSlogAttr(v *[]slog.Attr) (closer func()) {
	if *v != nil {
		return
	}
	var x, _ = reuse.SliceSlogAttr.Get().(*[]slog.Attr)
	if x == nil {
		var a = make([]slog.Attr, 0, 8)
		x = &a
	}
	*v = *x

	return func() {
		clear(*v)
		*v = (*v)[:0]
		reuse.SliceSlogAttr.Put(v)
	}
}
func reuseSliceError(v *[]error) (closer func()) {
	if *v != nil {
		return
	}
	var x, _ = reuse.SliceError.Get().(*[]error)
	if x == nil {
		var a = make([]error, 0, 8)
		x = &a
	}
	*v = *x

	return func() {
		clear(*v)
		*v = (*v)[:0]
		reuse.SliceError.Put(v)
	}
}
