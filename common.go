package flow

import (
	"embed"
	"fmt"
	"html/template"
	"reflect"
	"strings"
	"time"
)

const (
	keyUUID   = "uuid"
	keyMeta   = "meta"
	keyHook   = "hook"
	keyLive   = "live"
	keyOrigin = "origin"

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
	defer getStringBuilder(&wr)()
	var err = tmpl.ExecuteTemplate(&wr, "skip.js.go.tpl", nil)
	if err != nil {
		panic(err)
	}
	return wr.String()
}
func renderWalkJS(next []Pipe) string {
	var wr strings.Builder
	defer getStringBuilder(&wr)()
	var err = tmpl.ExecuteTemplate(&wr, "walk.js.go.tpl", next)
	if err != nil {
		panic(err)
	}
	return wr.String()
}
