package flow

import (
	"sync"
)

var syncpool struct {
	sliceError sync.Pool
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
