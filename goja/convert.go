//nolint:stylecheck,errcheck // todo: отрефакторить
package goja

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"time"

	"github.com/dop251/goja"
	"github.com/typomaker/flow"
	"github.com/typomaker/option"
)

func convert(rm *goja.Runtime, src, dst any) (err error) {
	switch src := src.(type) {
	case []flow.Node:
		err = convert_FlowNodeArray(rm, src, dst)
	case *lazyNodeArray:
		err = convert_LazyNodeArray(rm, src, dst)
	case flow.Node:
		err = convert_FlowNodeObject(rm, src, dst)
	case *lazyNodeObject:
		err = convert_LazyNodeObject(rm, src, dst)
	case flow.Meta:
		err = convert_FlowMetaObject(rm, src, dst)
	case flow.Hook:
		err = convert_FlowHookObject(rm, src, dst)
	case flow.Live:
		err = convert_FlowLiveObject(rm, src, dst)
	case *lazyLiveObject:
		err = convert_LazyLiveObject(rm, src, dst)
	case map[string]any:
		err = convert_Object(rm, src, dst)
	case *lazyObject:
		err = convert_LazyObject(rm, src, dst)
	case []any:
		err = convert_Array(rm, src, dst)
	case *lazyArray:
		err = convert_LazyArray(rm, src, dst)
	case goja.Value:
		err = convert_GojaValue(rm, src, dst)
	default:
		err = convert_Any(rm, src, dst)
	}
	return err
}
func convert_GojaValue(rm *goja.Runtime, src goja.Value, dst any) (err error) {
	switch dst := dst.(type) {
	case *[]flow.Node:
		err = convert_GojaValue_FlowList(rm, src, dst)
	case *flow.Node:
		err = convert_GojaValue_FlowNode(rm, src, dst)
	case *flow.Meta:
		err = convert_GojaValue_FlowMeta(rm, src, dst)
	case *flow.Hook:
		err = convert_GojaValue_FlowHook(rm, src, dst)
	case *flow.Live:
		err = convert_GojaValue_FlowLive(rm, src, dst)
	case *flow.Case:
		err = convert_GojaValue_FlowCase(rm, src, dst)
	case *flow.When:
		err = convert_GojaValue_FlowWhen(rm, src, dst)
	case *flow.Then:
		err = convert_GojaValue_FlowThen(rm, src, dst)
	case *any:
		err = convert_GojaValue_Any(rm, src, dst)
	default:
		err = newErrUnexpectedType(src, dst)
	}
	return err
}
func convert_GojaValue_Any(rm *goja.Runtime, src goja.Value, dst *any) (err error) {
	switch src.ExportType() {
	case reflectNull:
		*dst = nil
	case reflectLazyLiveObject:
		var s = src.Export().(*lazyLiveObject)
		err = convert_LazyLiveObject(rm, s, dst)
	case reflectLazyNodeObject:
		var s = src.Export().(*lazyNodeObject)
		err = convert_LazyNodeObject(rm, s, dst)
	case reflectLazyNodeArray:
		var s = src.Export().(*lazyNodeArray)
		err = convert_LazyNodeArray(rm, s, dst)
	case reflectLazyArray:
		var s = src.Export().(*lazyArray)
		err = convert_LazyArray(rm, s, dst)
	case reflectLazyObject:
		var s = src.Export().(*lazyObject)
		err = convert_LazyObject(rm, s, dst)
	case reflectArray:
		var d = src.Export().([]any)
		err = convert_GojaValue_Array(rm, src, &d)
		*dst = d
	case reflectObject:
		var d = src.Export().(map[string]any)
		err = convert_GojaValue_Object(rm, src, &d)
		*dst = d
	case reflectInt64:
		var d float64
		err = rm.ExportTo(src, &d)
		*dst = d
	case reflectTime:
		if d, ok := src.Export().(time.Time); ok {
			// время экспортируется как строка для совместимости с последующим маршалингом в structpb
			*dst = d.UTC().Format(time.RFC3339)
		} else {
			// обработка InvalidDate
			*dst = nil
		}
	default:
		err = rm.ExportTo(src, dst)
	}
	return err
}
func convert_GojaValue_FlowList(rm *goja.Runtime, src goja.Value, dst *[]flow.Node) (err error) {
	switch src.ExportType() {
	case reflectLazyNodeArray:
		var src = src.Export().(*lazyNodeArray)
		err = convert_LazyNodeArray_FlowNodeArray(rm, src, dst)
	case reflectArray:
		var src = src.(*goja.Object)
		var lng = int(src.Get("length").ToInteger())
		var laz = &lazyNodeArray{rm: rm, value: make([]goja.Value, lng)}
		for i := range lng {
			laz.Set(i, src.Get(strconv.Itoa(i)))
		}
		err = convert_LazyNodeArray_FlowNodeArray(rm, laz, dst)
	default:
		err = newErrUnexpectedType(src, dst)
	}
	return err
}
func convert_GojaValue_FlowNode(rm *goja.Runtime, src goja.Value, dst *flow.Node) (err error) {
	switch src.ExportType() {
	case reflectLazyNodeObject:
		var laz = src.Export().(*lazyNodeObject)
		err = convert_LazyItem_Item(rm, laz, dst)
	case reflectObject:
		var src = src.(*goja.Object)
		var keys = src.Keys()
		var laz = &lazyNodeObject{rm: rm}
		for _, key := range keys {
			laz.Set(key, src.Get(key))
		}
		err = convert_LazyItem_Item(rm, laz, dst)
	case reflectNull:
		*dst = flow.Node{}
	default:
		err = newErrUnexpectedType(src, dst)
	}
	return err
}
func convert_GojaValue_FlowMeta(rm *goja.Runtime, src goja.Value, dst *flow.Meta) (err error) {
	switch src.ExportType() {
	case reflectLazyObject:
		var src = src.Export().(*lazyObject)
		err = convert_LazyObject_Meta(rm, src, dst)
	case reflectObject:
		var src = src.(*goja.Object)
		var keys = src.Keys()
		var laz = &lazyObject{rm: rm, value: make(map[string]goja.Value, len(keys))}
		for _, key := range keys {
			laz.Set(key, src.Get(key))
		}
		err = convert_LazyObject_Meta(rm, laz, dst)
	default:
		err = newErrUnexpectedType(src, dst)
	}
	return err
}
func convert_GojaValue_FlowHook(rm *goja.Runtime, src goja.Value, dst *flow.Hook) (err error) {
	switch src.ExportType() {
	case reflectLazyObject:
		var laz = src.Export().(*lazyObject)
		err = convert_LazyObject_Object(rm, laz, &laz.proto)
		*dst = laz.proto
	case reflectObject:
		var src = src.(*goja.Object)
		var keys = src.Keys()
		var laz = &lazyObject{rm: rm, value: make(map[string]goja.Value, len(keys))}
		for _, key := range keys {
			laz.Set(key, src.Get(key))
		}
		err = convert_LazyObject_Object(rm, laz, &laz.proto)
		*dst = laz.proto
	default:
		err = newErrUnexpectedType(src, dst)
	}
	return err
}
func convert_LazyObject_Hook(rm *goja.Runtime, src *lazyObject, dst *flow.Hook) (err error) {
	if err = convert_LazyObject_Object(rm, src, &src.proto); err != nil {
		return err
	}
	*dst = src.proto
	return nil
}
func convert_GojaValue_FlowLive(rm *goja.Runtime, src goja.Value, dst *flow.Live) (err error) {
	switch src.ExportType() {
	case reflectLazyLiveObject:
		var laz = src.Export().(*lazyLiveObject)
		err = convert_LazyLiveObject_FlowLiveObject(rm, laz, dst)
	case reflectObject:
		var src = src.(*goja.Object)
		var keys = src.Keys()
		var laz = &lazyLiveObject{rm: rm}
		for _, key := range keys {
			laz.Set(key, src.Get(key))
		}
		err = convert_LazyLiveObject_FlowLiveObject(rm, laz, dst)
	default:
		err = newErrUnexpectedType(src, dst)
	}
	return err
}
func convert_GojaValue_FlowCase(rm *goja.Runtime, src goja.Value, dst *flow.Case) (err error) {
	if src == nil {
		return nil
	}
	if goja.IsUndefined(src) || goja.IsNull(src) {
		*dst = flow.Case{}
		return nil
	}
	var jsCase, ok = src.(*goja.Object)
	if !ok {
		return fmt.Errorf(`must be "Object"`)
	}
	if jsWhen := jsCase.Get(keyCaseWhen); jsWhen != nil {
		switch {
		case goja.IsUndefined(jsWhen):
			dst.When = flow.When{}
		case goja.IsNull(jsWhen):
			dst.When = flow.When{}
		default:
			if err = convert_GojaValue_FlowWhen(rm, jsWhen, &dst.When); err != nil {
				return fmt.Errorf(`when %w`, err)
			}
		}
	}
	if jsThen := jsCase.Get(keyCaseThen); jsThen != nil {
		switch {
		case goja.IsUndefined(jsThen):
			dst.When = flow.When{}
		case goja.IsNull(jsThen):
			dst.When = flow.When{}
		default:
			if err = convert_GojaValue_FlowThen(rm, jsThen, &dst.Then); err != nil {
				return fmt.Errorf(`then %w`, err)
			}
		}
	}
	return err
}

//nolint:funlen,gocognit // отрефакторить
func convert_GojaValue_FlowWhen(rm *goja.Runtime, src goja.Value, dst *flow.When) (err error) {
	if src == nil {
		return nil
	}
	if goja.IsUndefined(src) || goja.IsNull(src) {
		*dst = flow.When{}
		return nil
	}
	var jsWhen, ok = src.(*goja.Object)
	if !ok {
		return fmt.Errorf(`must be "Object"`)
	}
	if jsUUID := jsWhen.Get(keyUUID); jsUUID != nil {
		var goUUID []flow.UUID
		switch {
		case goja.IsUndefined(jsUUID):
			dst.UUID = option.Option[[]flow.UUID]{}
		case goja.IsNull(jsUUID):
			dst.UUID = option.None[[]flow.UUID]()
		default:
			var obj, ok = jsUUID.(*goja.Object)
			if !ok {
				return fmt.Errorf("uuid must be array")
			}
			var length = int(obj.Get("length").ToInteger())
			goUUID = make([]flow.UUID, 0, length)
			for i := 0; i < length; i++ {
				var jsUUIDItem = obj.Get(strconv.Itoa(i))
				var goUUIDItem flow.UUID
				switch {
				case goja.IsUndefined(jsUUIDItem):
					continue
				case goja.IsNull(jsUUIDItem):
					continue
				default:
					var goUUIDItemString string
					if err = rm.ExportTo(jsUUIDItem, &goUUIDItemString); err != nil {
						return fmt.Errorf("uuid %d %w", i, err)
					}
					if goUUIDItem, err = flow.ParseUUID(goUUIDItemString); err != nil {
						return fmt.Errorf(`uuid %d %w`, i, err)
					}
					goUUID = append(goUUID, goUUIDItem)
				}
			}
			dst.UUID = option.Some(goUUID)
		}
	}
	if jsHook := jsWhen.Get(keyHook); jsHook != nil {
		var goHook []flow.Hook
		switch {
		case goja.IsUndefined(jsHook):
			dst.Hook = option.Option[[]flow.Hook]{}
		case goja.IsNull(jsHook):
			dst.Hook = option.None[[]flow.Hook]()
		default:
			var obj, ok = jsHook.(*goja.Object)
			if !ok {
				return fmt.Errorf("hook must be array")
			}
			var length = int(obj.Get("length").ToInteger())
			goHook = make([]flow.Hook, 0, length)
			for i := 0; i < length; i++ {
				var jsHookItem = obj.Get(strconv.Itoa(i))
				var goHookItem = flow.Hook{}
				switch {
				case goja.IsUndefined(jsHookItem):
					continue
				case goja.IsNull(jsHookItem):
					continue
				default:
					if err = convert(rm, jsHookItem, &goHookItem); err != nil {
						return fmt.Errorf(`hook %d %w`, i, err)
					}
					goHook = append(goHook, goHookItem)
				}
			}
			dst.Hook = option.Some(goHook)
		}
	}
	return err
}

//nolint:funlen,gocognit // отрефакторить
func convert_GojaValue_FlowThen(rm *goja.Runtime, src goja.Value, dst *flow.Then) (err error) {
	if src == nil {
		return nil
	}

	if goja.IsUndefined(src) || goja.IsNull(src) {
		*dst = flow.Then{}
		return nil
	}
	var jsThen, ok = src.(*goja.Object)
	if !ok {
		return fmt.Errorf(`must be "Object"`)
	}
	if jsMeta := jsThen.Get(keyMeta); jsMeta != nil {
		var goMeta flow.Meta
		switch {
		case goja.IsUndefined(jsMeta):
			dst.Meta = option.Option[flow.Meta]{}
		case goja.IsNull(jsMeta):
			dst.Meta = option.None[flow.Meta]()
		default:
			if err = convert(rm, jsMeta, &goMeta); err != nil {
				return fmt.Errorf(`meta %w`, err)
			}
			dst.Meta = option.Some(goMeta)
		}
	}
	if jsHook := jsThen.Get(keyHook); jsHook != nil {
		var goHook flow.Hook
		switch {
		case goja.IsUndefined(jsHook):
			dst.Hook = option.Option[flow.Hook]{}
		case goja.IsNull(jsHook):
			dst.Hook = option.None[flow.Hook]()
		default:
			if err = convert(rm, jsHook, &goHook); err != nil {
				return fmt.Errorf(`hook %w`, err)
			}
			dst.Hook = option.Some(goHook)
		}
	}
	if jsLive := jsThen.Get(keyLive); jsLive != nil {
		var goLive flow.Live
		switch {
		case goja.IsUndefined(jsLive):
			dst.Live = option.Option[flow.Live]{}
		case goja.IsNull(jsLive):
			dst.Live = option.None[flow.Live]()
		default:
			if err = convert(rm, jsLive, &goLive); err != nil {
				return fmt.Errorf(`live %w`, err)
			}
			dst.Live = option.Some(goLive)
		}
	}
	return err
}
func convert_GojaValue_Object(rm *goja.Runtime, src goja.Value, dst *map[string]any) (err error) {
	var obj = src.(*goja.Object)
	var keys = obj.Keys()
	var d = make(map[string]any, len(keys))
	for _, key := range keys {
		var val = obj.Get(key)
		if goja.IsUndefined(val) {
			continue
		}
		var goVal any
		if err = convert_GojaValue(rm, val, &goVal); err != nil {
			return fmt.Errorf("%s %w", key, err)
		}
		d[key] = goVal
	}
	*dst = d
	return nil
}
func convert_GojaValue_Array(rm *goja.Runtime, src goja.Value, dst *[]any) (err error) {
	var obj = src.(*goja.Object)
	var length = int(obj.Get("length").ToInteger())
	var d = make([]any, 0, length)
	for i := 0; i < length; i++ {
		var key = strconv.Itoa(i)
		var val = obj.Get(key)
		if goja.IsUndefined(val) {
			continue
		}
		var goVal any
		if err = convert_GojaValue(rm, val, &goVal); err != nil {
			return fmt.Errorf("%d %w", i, err)
		}
		d = append(d, goVal)
	}
	*dst = d
	return nil
}
func convert_FlowNodeArray(rm *goja.Runtime, src []flow.Node, dst any) (err error) {
	switch dst := dst.(type) {
	case *goja.Value:
		err = convert_FlowNodeArray_GojaValue(rm, src, dst)
	case *any:
		err = convert_FlowNodeArray_Array(rm, src, dst)
	default:
		err = newErrUnexpectedType(src, dst)
	}
	return err
}
func convert_FlowNodeArray_GojaValue(rm *goja.Runtime, src []flow.Node, dst *goja.Value) (err error) {
	var d = *dst
	var t reflect.Type
	if d != nil {
		t = d.ExportType()
	}
	switch t {
	case reflectLazyNodeArray:
		var s = d.Export().(*lazyNodeArray)
		err = convert_FlowNodeArray_LazyNodeArray(rm, src, &s)
	default: // todo nil?
		*dst = rm.NewDynamicArray(&lazyNodeArray{rm: rm, proto: src})
	}
	return err
}
func convert_FlowNodeArray_Array(rm *goja.Runtime, src []flow.Node, dst *any) (err error) {
	var d, ok = (*dst).([]any)
	if !ok {
		d = make([]any, 0, len(src))
	} else {
		// todo: replace to [::]
		d = slices.Grow(d[:0], len(src))[:0]
	}
	for _, v := range src {
		var src = &lazyNodeObject{rm: rm, proto: v}
		var val any
		if err = convert_LazyNodeObject_Object(rm, src, &val); err != nil {
			return err
		}
		d = append(d, val)
	}
	*dst = d
	return err
}
func convert_FlowNodeArray_LazyNodeArray(rm *goja.Runtime, src []flow.Node, dst **lazyNodeArray) (err error) {
	var d = *dst
	d.value = d.value[0:len(src):len(src)]
	for i, v := range src {
		if err = convert_FlowNodeObject(rm, v, &d.value[i]); err != nil {
			return err
		}
	}
	*dst = d
	return err
}
func convert_FlowNodeObject(rm *goja.Runtime, src flow.Node, dst any) (err error) {
	switch dst := dst.(type) {
	case *goja.Value:
		err = convert_FlowNodeObject_GojaValue(rm, src, dst)
	case *any:
		err = convert_LazyNodeObject_Object(rm, &lazyNodeObject{proto: src}, dst)
	default:
		err = newErrUnexpectedType(src, dst)
	}
	return err
}
func convert_FlowNodeObject_GojaValue(rm *goja.Runtime, src flow.Node, dst *goja.Value) (err error) {
	var d = *dst
	var t reflect.Type
	if d != nil {
		t = d.ExportType()
	}
	switch t {
	case reflectLazyNodeObject:
		var s = d.Export().(*lazyNodeObject)
		err = convert_FlowNodeObject_LazyNodeObject(rm, src, &s)
	default:
		*dst = rm.NewDynamicObject(&lazyNodeObject{rm: rm, proto: src})
	}
	return err
}
func convert_FlowNodeObject_LazyNodeObject(rm *goja.Runtime, src flow.Node, dst **lazyNodeObject) (err error) {
	var d = *dst
	switch {
	case src.UUID.IsZero():
		d.value.UUID = goja.Undefined()
	case src.UUID.IsNone():
		d.value.UUID = goja.Null()
	default:
		d.value.UUID = rm.ToValue(src.UUID.Get())
	}
	switch {
	case src.Meta.IsZero():
		d.value.Meta = goja.Undefined()
	case src.Meta.IsNone():
		d.value.Meta = goja.Null()
	default:
		err = convert_FlowMetaObject_GojaValue(rm, src.Meta.Get(), &d.value.Meta)
	}
	switch {
	case src.Hook.IsZero():
		d.value.Hook = goja.Undefined()
	case src.Hook.IsNone():
		d.value.Hook = goja.Null()
	default:
		err = convert_FlowHookObject_GojaValue(rm, src.Hook.Get(), &d.value.Hook)
	}
	switch {
	case src.Live.IsZero():
		d.value.Live = goja.Undefined()
	case src.Live.IsNone():
		d.value.Live = goja.Null()
	default:
		err = convert_FlowLiveObject_GojaValue(rm, src.Live.Get(), &d.value.Live)
	}
	return err
}
func convert_FlowMetaObject(rm *goja.Runtime, src flow.Meta, dst any) (err error) {
	switch dst := dst.(type) {
	case *goja.Value:
		err = convert_FlowMetaObject_GojaValue(rm, src, dst)
	case *any:
		err = convert_FlowMetaObject_Any(rm, src, dst)
	default:
		err = newErrUnexpectedType(src, dst)
	}
	return err
}
func convert_FlowMetaObject_Any(rm *goja.Runtime, src flow.Meta, dst *any) (err error) {
	return convert_Object_Any(rm, src, dst)
}
func convert_FlowMetaObject_GojaValue(rm *goja.Runtime, src flow.Meta, dst *goja.Value) (err error) {
	//todo: add convert to passed lazyObject insted of creating new
	*dst = rm.NewDynamicObject(&lazyObject{rm: rm, proto: src})
	return err
}
func convert_FlowHookObject(rm *goja.Runtime, src flow.Hook, dst any) (err error) {
	switch dst := dst.(type) {
	case *goja.Value:
		err = convert_FlowHookObject_GojaValue(rm, src, dst)
	case *any:
		err = convert_FlowHookObject_Any(rm, src, dst)
	default:
		err = newErrUnexpectedType(src, dst)
	}
	return err
}
func convert_FlowHookObject_GojaValue(rm *goja.Runtime, src flow.Hook, dst *goja.Value) (err error) {
	*dst = rm.NewDynamicObject(&lazyObject{rm: rm, proto: src})
	return err
}
func convert_FlowHookObject_Any(rm *goja.Runtime, src flow.Hook, dst *any) (err error) {
	err = convert_Object_Any(rm, src, dst)
	return err
}
func convert_FlowLiveObject(rm *goja.Runtime, src flow.Live, dst any) (err error) {
	switch dst := dst.(type) {
	case *goja.Value:
		err = convert_FlowLiveObject_GojaValue(rm, src, dst)
	case *any:
		err = convert_FlowLiveObject_Any(rm, src, dst)
	default:
		err = newErrUnexpectedType(src, dst)
	}
	return err
}
func convert_FlowLiveObject_Any(rm *goja.Runtime, src flow.Live, dst *any) (err error) {
	var d, ok = (*dst).(map[string]any)
	if !ok {
		d = make(map[string]any, 2)
	} else {
		clear(d)
	}
	switch {
	case src.Since.IsNone():
		d["since"] = nil
	case src.Since.IsSome():
		d["since"] = src.Since.Get().Format(time.RFC3339)
	}
	switch {
	case src.Until.IsNone():
		d["until"] = nil
	case src.Until.IsSome():
		d["until"] = src.Until.Get().Format(time.RFC3339)
	}
	*dst = d
	return nil
}
func convert_FlowLiveObject_GojaValue(rm *goja.Runtime, src flow.Live, dst *goja.Value) (err error) {
	*dst = rm.NewDynamicObject(&lazyLiveObject{rm: rm, proto: src})
	return err
}
func convert_Any(rm *goja.Runtime, src any, dst any) (err error) {
	switch dst := dst.(type) {
	case *goja.Value:
		err = convert_Any_GojaValue(rm, src, dst)
	case *any:
		err = convert_Any_Any(rm, src, dst)
	default:
		err = newErrUnexpectedType(src, dst)
	}
	return err
}
func convert_Any_GojaValue(rm *goja.Runtime, src any, dst *goja.Value) (err error) {
	var goAny any
	if err = convert_Any_Any(rm, src, &goAny); err != nil {
		return err
	}
	*dst = rm.ToValue(goAny)
	return nil
}
func convert_Any_Any(_ *goja.Runtime, src any, dst *any) (err error) {
	switch src := src.(type) {
	case int:
		*dst = float64(src)
	case int8:
		*dst = float64(src)
	case int16:
		*dst = float64(src)
	case int32:
		*dst = float64(src)
	case int64:
		*dst = float64(src)
	case uint:
		*dst = float64(src)
	case uint8:
		*dst = float64(src)
	case uint16:
		*dst = float64(src)
	case uint32:
		*dst = float64(src)
	case uint64:
		*dst = float64(src)
	case float32:
		*dst = float64(src)
	case time.Time:
		*dst = src.Format(time.RFC3339)
	default:
		*dst = src
	}
	return err
}
func convert_LazyNodeArray(rm *goja.Runtime, src *lazyNodeArray, dst any) (err error) {
	switch dst := dst.(type) {
	case *[]flow.Node:
		err = convert_LazyNodeArray_FlowNodeArray(rm, src, dst)
	case *any:
		err = convert_LazyNodeArray_Array(rm, src, dst)
	default:
		err = newErrUnexpectedType(src, dst)
	}
	return err
}
func convert_LazyNodeArray_Array(rm *goja.Runtime, src *lazyNodeArray, dst *any) (err error) {
	if src.value == nil {
		return convert_FlowNodeArray_Array(rm, src.proto, dst)
	}
	var d, ok = (*dst).([]any)
	if !ok {
		d = make([]any, 0, src.Len())
	} else {
		// todo: replace to [::]
		d = slices.Grow(d[:0], src.Len())[:0]
	}
	for i := 0; i < len(src.value); i++ {
		if goja.IsUndefined(src.value[i]) {
			continue
		}
		var jsVal, goVal any
		var idx = len(d)
		d = d[:idx+1]
		if jsVal = src.value[i]; jsVal == nil {
			if i >= len(src.proto) {
				d[idx] = nil
				continue
			}
			jsVal = &lazyNodeObject{rm: rm, proto: src.proto[i]}
		}
		if err = convert(rm, jsVal, &goVal); err != nil {
			return fmt.Errorf("%d %w", i, err)
		}
		d[idx] = goVal
	}
	*dst = d
	return nil
}
func convert_LazyNodeArray_FlowNodeArray(rm *goja.Runtime, src *lazyNodeArray, dst *[]flow.Node) (err error) {
	var s = src.value
	var d = src.proto
	if s != nil && len(d) > len(s) {
		d = d[:len(s)]
	}
	for i := 0; i < len(s); i++ {
		var v = s[i]
		// без изменений
		if v == nil {
			continue
		}
		// удалено
		if goja.IsUndefined(v) {
			// удалить если было определено
			if i < len(d) {
				copy(d[i:], d[i+1:])
				d = d[:len(d)-1]
				copy(s[i:], s[i+1:])
				s = s[:len(s)-1]
				i--
			}
			continue
		}
		var p flow.Node
		if err = convert(rm, v, &p); err != nil {
			return fmt.Errorf("%d %w", i, err)
		}
		if i < len(d) {
			// переопределить
			d[i] = p
		} else {
			// добавить
			d = append(d, p)
		}
	}
	*dst = d
	return nil
}
func convert_LazyNodeObject(rm *goja.Runtime, src *lazyNodeObject, dst any) (err error) {
	switch dst := dst.(type) {
	case *flow.Node:
		err = convert_LazyItem_Item(rm, src, dst)
	case *any:
		err = convert_LazyNodeObject_Object(rm, src, dst)
	default:
		err = newErrUnexpectedType(src, dst)
	}
	return err
}

//nolint:funlen,gocyclo,gocognit,cyclop // отрефакторить если есть желание
func convert_LazyNodeObject_Object(rm *goja.Runtime, src *lazyNodeObject, dst *any) (err error) {
	var d, ok = (*dst).(map[string]any)
	if !ok {
		d = make(map[string]any, 6)
	}

	switch {
	case goja.IsUndefined(src.value.UUID):
		delete(d, "uuid")
	case goja.IsNull(src.value.UUID):
		d["uuid"] = nil
	case src.value.UUID != nil:
		var goUUID string
		if err = rm.ExportTo(src.value.UUID, &goUUID); err != nil {
			return fmt.Errorf("uuid %w", err)
		}
		d["uuid"] = goUUID
	case src.proto.UUID.IsZero():
		delete(d, "uuid")
	case src.proto.UUID.IsNone():
		d["uuid"] = nil
	case src.proto.UUID.IsSome():
		d["uuid"] = src.proto.UUID.Get().String()
	}

	switch {
	case goja.IsUndefined(src.value.Meta):
		delete(d, "meta")
	case goja.IsNull(src.value.Meta):
		d["meta"] = nil
	case src.value.Meta != nil:
		var goMeta any
		if err = convert(rm, src.value.Meta, &goMeta); err != nil {
			return fmt.Errorf("meta %w", err)
		}
		d["meta"] = goMeta
	case src.proto.Meta.IsZero():
		delete(d, "meta")
	case src.proto.Meta.IsNone():
		d["meta"] = nil
	case src.proto.Meta.IsSome():
		var goMeta any
		if err = convert(rm, src.proto.Meta.Get(), &goMeta); err != nil {
			return err
		}
		d["meta"] = goMeta
	}

	switch {
	case goja.IsUndefined(src.value.Hook):
		delete(d, "hook")
	case goja.IsNull(src.value.Hook):
		d["hook"] = nil
	case src.value.Hook != nil:
		var goHook any
		if err = convert(rm, src.value.Hook, &goHook); err != nil {
			return fmt.Errorf("hook %w", err)
		}
		d["hook"] = goHook
	case src.proto.Hook.IsZero():
		delete(d, "hook")
	case src.proto.Hook.IsNone():
		d["hook"] = nil
	case src.proto.Hook.IsSome():
		var goHook any
		if err = convert(rm, src.proto.Hook.Get(), &goHook); err != nil {
			return err
		}
		d["hook"] = goHook
	}

	switch {
	case goja.IsUndefined(src.value.Live):
		delete(d, "live")
	case goja.IsNull(src.value.Live):
		d["live"] = nil
	case src.value.Live != nil:
		var goLive any
		if err = convert(rm, src.value.Live, &goLive); err != nil {
			return fmt.Errorf("live %w", err)
		}
		d["live"] = goLive
	case src.proto.Live.IsZero():
		delete(d, "live")
	case src.proto.Live.IsNone():
		d["live"] = nil
	case src.proto.Live.IsSome():
		var goLive any
		if err = convert(rm, src.proto.Live.Get(), &goLive); err != nil {
			return err
		}
		d["live"] = goLive
	}

	switch {
	case goja.IsUndefined(src.value.Origin):
		delete(d, "origin")
	case goja.IsNull(src.value.Origin):
		d["origin"] = nil
	case src.value.Origin != nil:
		var goOrigin any
		if err = convert(rm, src.value.Origin, &goOrigin); err != nil {
			return fmt.Errorf("origin %w", err)
		}
		d["origin"] = goOrigin
	case src.proto.Origin().IsZero():
		delete(d, "origin")
	default:
		var goOrigin any
		if err = convert(rm, src.proto.Origin(), &goOrigin); err != nil {
			return err
		}
		d["origin"] = goOrigin
	}

	*dst = d
	return nil
}

//nolint:cyclop,funlen,gocognit // todo: отрефакторить
func convert_LazyItem_Item(rm *goja.Runtime, src *lazyNodeObject, dst *flow.Node) (err error) {
	*dst = src.proto
	if jsUUID := src.value.UUID; jsUUID != nil {
		switch {
		case goja.IsUndefined(jsUUID):
			dst.UUID = option.Option[flow.UUID]{}
		case goja.IsNull(jsUUID):
			dst.UUID = option.None[flow.UUID]()
		default:
			var goString string
			if err = rm.ExportTo(jsUUID, &goString); err != nil {
				return fmt.Errorf("uuid %w", err)
			}
			var goUUID flow.UUID
			if goUUID, err = flow.ParseUUID(goString); err != nil {
				return fmt.Errorf("uuid %w", err)
			}
			dst.UUID = option.Some(goUUID)
		}
	}
	if jsMeta := src.value.Meta; jsMeta != nil {
		switch {
		case goja.IsUndefined(jsMeta):
			dst.Meta = option.Option[flow.Meta]{}
		case goja.IsNull(jsMeta):
			dst.Meta = option.None[flow.Meta]()
		default:
			var goMeta flow.Meta
			if err = convert(rm, jsMeta, &goMeta); err != nil {
				return fmt.Errorf("meta %w", err)
			}
			dst.Meta = option.Some(goMeta)
		}
	}
	if jsHook := src.value.Hook; jsHook != nil {
		switch {
		case goja.IsUndefined(jsHook):
			dst.Hook = option.Option[flow.Hook]{}
		case goja.IsNull(jsHook):
			dst.Hook = option.None[flow.Hook]()
		default:
			var goHook flow.Hook
			if err = convert(rm, jsHook, &goHook); err != nil {
				return fmt.Errorf("hook %w", err)
			}
			dst.Hook = option.Some(goHook)
		}
	}
	if jsLive := src.value.Live; jsLive != nil {
		switch {
		case goja.IsUndefined(jsLive):
			dst.Live = option.Option[flow.Live]{}
		case goja.IsNull(jsLive):
			dst.Live = option.None[flow.Live]()
		default:
			var goLive flow.Live
			if err = convert(rm, jsLive, &goLive); err != nil {
				return fmt.Errorf("live %w", err)
			}
			dst.Live = option.Some(goLive)
		}
	}
	if jsOrigin := src.value.Origin; jsOrigin != nil {
		switch {
		case goja.IsUndefined(jsOrigin), goja.IsNull(jsOrigin):
			dst.SetOrigin(flow.Node{})
		default:
			var goOrigin flow.Node
			if err = convert(rm, jsOrigin, &goOrigin); err != nil {
				return fmt.Errorf("rate %w", err)
			}
			dst.SetOrigin(goOrigin)
		}
	}
	return nil
}
func convert_LazyLiveObject(rm *goja.Runtime, src *lazyLiveObject, dst any) (err error) {
	switch dst := dst.(type) {
	case *flow.Live:
		err = convert_LazyLiveObject_FlowLiveObject(rm, src, dst)
	case *any:
		err = convert_LazyLiveObject_Object(rm, src, dst)
	default:
		err = newErrUnexpectedType(src, dst)
	}
	return err
}
func convert_LazyLiveObject_FlowLiveObject(rm *goja.Runtime, src *lazyLiveObject, dst *flow.Live) (err error) {
	*dst = src.proto
	if src.value.Since != nil {
		switch {
		case goja.IsUndefined(src.value.Since):
			dst.Since = option.Option[time.Time]{}
		case goja.IsNull(src.value.Since):
			dst.Since = option.None[time.Time]()
		default:
			var goSince time.Time
			if err = rm.ExportTo(src.value.Since, &goSince); err != nil {
				return fmt.Errorf("since %w", err)
			}
			dst.Since = option.Some(goSince.In(time.UTC))
		}
	}
	if src.value.Until != nil {
		switch {
		case goja.IsUndefined(src.value.Until):
			dst.Until = option.Option[time.Time]{}
		case goja.IsNull(src.value.Until):
			dst.Until = option.None[time.Time]()
		default:
			var goUntil time.Time
			if err = rm.ExportTo(src.value.Until, &goUntil); err != nil {
				return fmt.Errorf("until %w", err)
			}
			dst.Until = option.Some(goUntil.In(time.UTC))
		}
	}
	return nil
}
func convert_LazyLiveObject_Object(rm *goja.Runtime, src *lazyLiveObject, dst *any) (err error) {
	var d, ok = (*dst).(map[string]any)
	if !ok {
		d = make(map[string]any, 2)
	}

	switch {
	case goja.IsUndefined(src.value.Since):
		delete(d, "since")
	case goja.IsNull(src.value.Since):
		d["since"] = nil
	case src.value.Since != nil:
		var goSince time.Time
		if err = rm.ExportTo(src.value.Since, &goSince); err != nil {
			return fmt.Errorf("since %w", err)
		}
		d["since"] = goSince.In(time.UTC).Format(time.RFC3339)
	case src.proto.Since.IsZero():
		delete(d, "since")
	case src.proto.Since.IsNone():
		d["since"] = nil
	case src.proto.Since.IsSome():
		d["since"] = src.proto.Since.Get().Format(time.RFC3339)
	}

	switch {
	case goja.IsUndefined(src.value.Until):
		delete(d, "until")
	case goja.IsNull(src.value.Until):
		d["until"] = nil
	case src.value.Until != nil:
		var goUntil time.Time
		if err = rm.ExportTo(src.value.Until, &goUntil); err != nil {
			return fmt.Errorf("until %w", err)
		}
		d["until"] = goUntil.In(time.UTC).Format(time.RFC3339)
	case src.proto.Until.IsZero():
		delete(d, "until")
	case src.proto.Until.IsNone():
		d["until"] = nil
	case src.proto.Until.IsSome():
		d["until"] = src.proto.Until.Get().Format(time.RFC3339)
	}
	*dst = d
	return nil
}
func convert_LazyObject_Meta(rm *goja.Runtime, src *lazyObject, dst *flow.Meta) (err error) {
	if err = convert_LazyObject_Object(rm, src, &src.proto); err != nil {
		return err
	}
	*dst = src.proto
	return nil
}
func convert_LazyObject(rm *goja.Runtime, src *lazyObject, dst any) (err error) {
	switch dst := dst.(type) {
	case *map[string]any:
		err = convert_LazyObject_Object(rm, src, dst)
	case *any:
		err = convert_LazyObject_Object(rm, src, &src.proto)
		*dst = src.proto
	default:
		err = newErrUnexpectedType(src, dst)
	}
	return err
}
func convert_LazyObject_Object(rm *goja.Runtime, src *lazyObject, dst *map[string]any) (err error) {
	*dst = src.proto
	var d = *dst
	if d == nil {
		d = make(map[string]any, len(src.value))
	}
	for key, val := range src.value {
		if val == nil {
			continue
		}
		if goja.IsUndefined(val) {
			delete(d, key)
			continue
		}
		var protoVal any
		if err = convert(rm, val, &protoVal); err != nil {
			return fmt.Errorf("%s %w", key, err)
		}
		d[key] = protoVal
	}
	if dst != &src.proto {
		for key, protoVal := range src.proto {
			if _, ok := src.value[key]; ok {
				continue
			}
			d[key] = protoVal
		}
	}
	*dst = d
	return nil
}
func convert_LazyArray(rm *goja.Runtime, src *lazyArray, dst any) (err error) {
	switch dst := dst.(type) {
	case *[]any:
		err = convert_LazyArray_Array(rm, src, dst)
		*dst = src.proto
	case *any:
		err = convert_LazyArray_Array(rm, src, &src.proto)
		*dst = src.proto
	default:
		err = newErrUnexpectedType(src, dst)
	}
	return err
}
func convert_LazyArray_Array(rm *goja.Runtime, src *lazyArray, dst *[]any) (err error) {
	var s = src.value
	var d = src.proto
	if s != nil && len(d) > len(s) {
		d = d[:len(s)]
	}
	for i := 0; i < len(s); i++ {
		var v = s[i]
		// без изменений
		if v == nil {
			continue
		}
		// удалено
		if goja.IsUndefined(v) {
			// удалить если было определено
			if i < len(d) {
				copy(d[i:], d[i+1:])
				d = d[:len(d)-1]
				copy(s[i:], s[i+1:])
				s = s[:len(s)-1]
				i--
			}
			continue
		}
		var p any
		if err = convert(rm, v, &p); err != nil {
			return fmt.Errorf("%d %w", i, err)
		}
		if i < len(d) {
			// переопределить
			d[i] = p
		} else {
			// добавить
			d = append(d, p)
		}
	}
	*dst = d
	return nil
}
func convert_Object(rm *goja.Runtime, src map[string]any, dst any) (err error) {
	switch dst := dst.(type) {
	case *goja.Value:
		*dst = rm.NewDynamicObject(&lazyObject{rm: rm, proto: src})
	case *any:
		err = convert_Object_Any(rm, src, dst)
	default:
		err = newErrUnexpectedType(src, dst)
	}
	return err
}
func convert_Object_Any(rm *goja.Runtime, src map[string]any, dst *any) (err error) {
	var d, ok = (*dst).(map[string]any)
	if !ok {
		d = make(map[string]any, len(src))
	}
	for key, val := range src {
		var goVal any
		if err = convert(rm, val, &goVal); err != nil {
			return err
		}
		d[key] = goVal
	}
	*dst = d
	return nil
}
func convert_Array(rm *goja.Runtime, src []any, dst any) (err error) {
	switch dst := dst.(type) {
	case *goja.Value:
		*dst = rm.NewDynamicArray(&lazyArray{rm: rm, proto: src})
	case *any:
		err = convert_Array_Any(rm, src, dst)
	default:
		err = newErrUnexpectedType(src, dst)
	}
	return err
}
func convert_Array_Any(rm *goja.Runtime, src []any, dst *any) (err error) {
	var d, ok = (*dst).([]any)
	if !ok {
		d = make([]any, len(src))
	}
	for idx, val := range src {
		var goVal any
		if err = convert(rm, val, &goVal); err != nil {
			return err
		}
		d[idx] = goVal
	}
	*dst = d
	return nil
}

var ErrUnexpected = fmt.Errorf("unexpected")

func newErrUnexpectedType(src, dst any) error {
	return fmt.Errorf("convert %T to %T is %w", src, dst, ErrUnexpected)
}
