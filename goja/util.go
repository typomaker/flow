package goja

import (
	"reflect"
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

	keyCaseWhen = "when"
	keyCaseThen = "then"
)

var (
	reflectLazyNodeArray  = reflect.TypeOf((*lazyNodeArray)(nil))
	reflectLazyNodeObject = reflect.TypeOf((*lazyNodeObject)(nil))
	reflectLazyLiveObject = reflect.TypeOf((*lazyLiveObject)(nil))
	reflectLazyObject     = reflect.TypeOf((*lazyObject)(nil))
	reflectLazyArray      = reflect.TypeOf((*lazyArray)(nil))
	reflectInt64          = reflect.TypeOf(int64(0))
	reflectObject         = reflect.TypeOf((map[string]any)(nil))
	reflectArray          = reflect.TypeOf(([]any)(nil))
	reflectNull           = reflect.TypeOf(nil)
	reflectTime           = reflect.TypeOf(time.Time{})
	reflectString         = reflect.TypeOf("")
)
