package flow

import (
	"bytes"
	"fmt"
	"time"
)

func deepCopy(v any) any {
	switch v := v.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		bool, string, float32, float64, nil:
		return v
	case []any:
		var cp = make([]any, len(v))
		for i := range v {
			cp[i] = deepCopy(v[i])
		}
		return cp
	case map[string]any:
		var cp = make(map[string]any, len(v))
		for i := range v {
			cp[i] = deepCopy(v[i])
		}
		return cp
	case time.Time:
		return v
	default:
		panic(fmt.Sprintf("unexpected type: %T", v))
	}
}
func deepWith(l, r any, merge bool) any {
	switch l := l.(type) {
	case map[string]any:
		if r, ok := r.(map[string]any); ok {
			if merge {
				for k := range r {
					l[k] = deepWith(l[k], r[k], k[0] == '$')
				}
				return l
			} else {
				for k := range r {
					r[k] = deepWith(l[k], r[k], k[0] == '$')
				}
				return r
			}
		}
		return r
	case []any:
		if r, ok := r.([]any); ok && merge {
			return append(l, r...)
		}
		return r
	default:
		return r
	}
}
func deepHave(l, r any) bool {
	switch source := l.(type) {
	case map[string]any:
		if part, ok := r.(map[string]any); ok {
			for k := range part {
				if !deepHave(source[k], part[k]) {
					return false
				}
			}
			return true
		}
	case []any:
		if part, ok := r.([]any); ok {
			for i := range part {
				var ok bool
				for j := range source {
					if ok = deepHave(source[j], part[i]); ok {
						break
					}
				}
				if !ok {
					return false
				}
			}
			return true
		}
	default:
		return equal(source, r)
	}
	return false
}
func equal[T1, T2 any](l T1, r T2) bool {
	switch l := any(l).(type) {
	case interface{ Equal(T2) bool }:
		if r, ok := any(r).(T2); ok && l.Equal(r) {
			return true
		}
	case int:
		if r, ok := any(r).(int); ok && l == r {
			return true
		}
	case byte:
		if r, ok := any(r).(byte); ok && l == r {
			return true
		}
	case int8:
		if r, ok := any(r).(int8); ok && l == r {
			return true
		}
	case int16:
		if r, ok := any(r).(int16); ok && l == r {
			return true
		}
	case int32:
		if r, ok := any(r).(int32); ok && l == r {
			return true
		}
	case int64:
		if r, ok := any(r).(int64); ok && l == r {
			return true
		}
	case bool:
		if r, ok := any(r).(bool); ok && l == r {
			return true
		}
	case string:
		if r, ok := any(r).(string); ok && l == r {
			return true
		}
	case float32:
		if r, ok := any(r).(float32); ok && l == r {
			return true
		}
	case float64:
		if r, ok := any(r).(float64); ok && l == r {
			return true
		}
	case []byte:
		if r, ok := any(r).([]byte); ok && bytes.Equal(l, r) {
			return true
		}
	case nil:
		if any(r) == nil {
			return true
		}
	case map[string]any:
		if r, ok := any(r).(map[string]any); ok {
			if len(l) != len(r) {
				return false
			}
			for k := range r {
				if !equal(l[k], r[k]) {
					return false
				}
			}
			return true
		}
	case []any:
		if r, ok := any(r).([]any); ok {
			if len(l) != len(r) {
				return false
			}
			for i := range l {
				if !equal(l[i], r[i]) {
					return false
				}
			}
			return true
		}
	case time.Time:
		if r, ok := any(r).(time.Time); ok && l.Equal(r) {
			return true
		}
	default:
		panic(fmt.Sprintf("unexpected type: %T", l))
	}

	return false
}
