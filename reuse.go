package flow

import (
	"log/slog"
	"strings"
	"sync"
)

var syncpool struct {
	sliceNode       sync.Pool
	sliceSliceNode  sync.Pool
	sliceSlicePipe  sync.Pool
	slicePipe       sync.Pool
	mapStringStruct sync.Pool
	stringsBuilder  sync.Pool
	sliceSlogAttr   sync.Pool
	sliceError      sync.Pool
}

func getSliceNode(v *[]Node) (put func()) {
	if *v != nil {
		return
	}
	var x, _ = syncpool.sliceNode.Get().(*[]Node)
	if x == nil {
		var a = make([]Node, 0, 8)
		x = &a
	}
	*v = *x

	return func() {
		clear(*v)
		*v = (*v)[:0]
		syncpool.sliceNode.Put(v)
	}
}
func getSliceSliceNode(v *[][]Node) (put func()) {
	if *v != nil {
		return
	}
	var x, _ = syncpool.sliceSliceNode.Get().(*[][]Node)
	if x == nil {
		var a = make([][]Node, 0, 8)
		x = &a
	}
	*v = *x

	return func() {
		clear(*v)
		*v = (*v)[:0]
		syncpool.sliceSliceNode.Put(v)
	}
}
func getSlicePipe(v *[]Pipe) (put func()) {
	if *v != nil {
		return
	}
	var x, _ = syncpool.slicePipe.Get().(*[]Pipe)
	if x == nil {
		var a = make([]Pipe, 0, 8)
		x = &a
	}
	*v = *x

	return func() {
		clear(*v)
		*v = (*v)[:0]
		syncpool.slicePipe.Put(v)
	}
}
func getMapStringSrtuct(v *map[string]struct{}) (put func()) {
	if *v != nil {
		return
	}
	var x, _ = syncpool.mapStringStruct.Get().(*map[string]struct{})
	if x == nil {
		var a = make(map[string]struct{}, 8)
		x = &a
	}
	*v = *x

	return func() {
		clear(*v)
		syncpool.mapStringStruct.Put(v)
	}
}
func getStringBuilder(v *strings.Builder) (put func()) {
	if x, _ := syncpool.stringsBuilder.Get().(*strings.Builder); x != nil {
		*v = *x
	} else {
		*v = strings.Builder{}
	}
	return func() {
		v.Reset()
		syncpool.stringsBuilder.Put(v)
	}
}
func getSliceSlogAttr(v *[]slog.Attr) (put func()) {
	if *v != nil {
		return
	}
	var x, _ = syncpool.sliceSlogAttr.Get().(*[]slog.Attr)
	if x == nil {
		var a = make([]slog.Attr, 0, 8)
		x = &a
	}
	*v = *x

	return func() {
		clear(*v)
		*v = (*v)[:0]
		syncpool.sliceSlogAttr.Put(v)
	}
}
func getSliceError(v *[]error) (put func()) {
	if *v != nil {
		return
	}
	var x, _ = syncpool.sliceError.Get().(*[]error)
	if x == nil {
		var a = make([]error, 0, 8)
		x = &a
	}
	*v = *x

	return func() {
		clear(*v)
		*v = (*v)[:0]
		syncpool.sliceError.Put(v)
	}
}
