package goja

import (
	"testing"

	"github.com/typomaker/flow"
	"github.com/typomaker/flow/flowtest"
)

func TestJS(t *testing.T) {
	flowtest.TestJS(t, func(t *testing.T) flow.Handler {
		return New(flowtest.JSFSPath)
	})
}
