package flow

import (
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"reflect"
	"strings"
	"sync"
	"time"
)

const (
	keyUUID = "uuid"
	keyKind = "kind"
	keyMeta = "meta"
	keyHook = "hook"
	keyLive = "live"
	keyRoot = "root"

	keyLiveSince = "since"
	keyLiveUntil = "until"
)

var (
	reflectLazyList   = reflect.TypeOf((*lazyList)(nil))
	reflectLazyNode   = reflect.TypeOf((*lazyNode)(nil))
	reflectLazyLive   = reflect.TypeOf((*lazyLive)(nil))
	reflectLazyObject = reflect.TypeOf((*lazyObject)(nil))
	reflectLazyArray  = reflect.TypeOf((*lazyArray)(nil))
	reflectInt64      = reflect.TypeOf(int64(0))
	reflectObject     = reflect.TypeOf((map[string]any)(nil))
	reflectArray      = reflect.TypeOf(([]any)(nil))
	reflectNull       = reflect.TypeOf(nil)
	reflectTime       = reflect.TypeOf(time.Time{})
	reflectString     = reflect.TypeOf("")
)

var ErrUnexpected = fmt.Errorf("unexpected")

func newErrUnexpectedConvert(src, dst any) error {
	return fmt.Errorf("convert %T to %T %w", src, dst, ErrUnexpected)
}

//go:embed *.go.tpl
var tmplFS embed.FS
var tmpl = template.Must(template.New("").ParseFS(tmplFS, "*"))

func renderSkipJS() string {
	var wr strings.Builder
	defer useStringBuilder(&wr)()
	var err = tmpl.ExecuteTemplate(&wr, "skip.js.go.tpl", nil)
	if err != nil {
		panic(err)
	}
	return wr.String()
}
func renderLoopJS(next []Pipe) string {
	var wr strings.Builder
	defer useStringBuilder(&wr)()
	var err = tmpl.ExecuteTemplate(&wr, "walk.js.go.tpl", next)
	if err != nil {
		panic(err)
	}
	return wr.String()
}

var sp struct {
	SliceSliceNode sync.Pool
	SliceSlicePipe sync.Pool
	SlicePipe      sync.Pool
	MapUUIDStruct  sync.Pool
	StringsBuilder sync.Pool
	SliceSlogAttr  sync.Pool
	SliceError     sync.Pool
}

func useSliceSliceNode(v *[][]Node) (closer func()) {
	if *v != nil {
		return
	}
	var x, _ = sp.SliceSliceNode.Get().(*[][]Node)
	if x == nil {
		var a = make([][]Node, 0, 8)
		x = &a
	}
	*v = *x

	return func() {
		clear(*v)
		*v = (*v)[:0]
		sp.SliceSliceNode.Put(v)
	}
}
func useSliceSlicePipe(v *[][]Pipe) (closer func()) {
	if *v != nil {
		return
	}
	var x, _ = sp.SliceSlicePipe.Get().(*[][]Pipe)
	if x == nil {
		var a = make([][]Pipe, 0, 8)
		x = &a
	}
	*v = *x

	return func() {
		clear(*v)
		*v = (*v)[:0]
		sp.SliceSlicePipe.Put(v)
	}
}
func useSlicePipe(v *[]Pipe) (closer func()) {
	if *v != nil {
		return
	}
	var x, _ = sp.SlicePipe.Get().(*[]Pipe)
	if x == nil {
		var a = make([]Pipe, 0, 8)
		x = &a
	}
	*v = *x

	return func() {
		clear(*v)
		*v = (*v)[:0]
		sp.SlicePipe.Put(v)
	}
}
func useMapUUIDSrtuct(v *map[UUID]struct{}) (closer func()) {
	if *v != nil {
		return
	}
	var x, _ = sp.MapUUIDStruct.Get().(*map[UUID]struct{})
	if x == nil {
		var a = make(map[UUID]struct{}, 8)
		x = &a
	}
	*v = *x

	return func() {
		clear(*v)
		sp.MapUUIDStruct.Put(v)
	}
}
func useStringBuilder(v *strings.Builder) (closer func()) {
	if x, _ := sp.StringsBuilder.Get().(*strings.Builder); x != nil {
		*v = *x
	} else {
		*v = strings.Builder{}
	}
	return func() {
		v.Reset()
		sp.StringsBuilder.Put(v)
	}
}
func useSliceSlogAttr(v *[]slog.Attr) (closer func()) {
	if *v != nil {
		return
	}
	var x, _ = sp.SliceSlogAttr.Get().(*[]slog.Attr)
	if x == nil {
		var a = make([]slog.Attr, 0, 8)
		x = &a
	}
	*v = *x

	return func() {
		clear(*v)
		*v = (*v)[:0]
		sp.SliceSlogAttr.Put(v)
	}
}
func useSliceError(v *[]error) (closer func()) {
	if *v != nil {
		return
	}
	var x, _ = sp.SliceError.Get().(*[]error)
	if x == nil {
		var a = make([]error, 0, 8)
		x = &a
	}
	*v = *x

	return func() {
		clear(*v)
		*v = (*v)[:0]
		sp.SliceError.Put(v)
	}
}
