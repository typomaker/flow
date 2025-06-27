package flow

import (
	"testing"
	"time"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/require"
	"github.com/typomaker/option"
)

func TestGetNodeUUID(t *testing.T) {
	rm := goja.New()
	node := Node{
		UUID: option.Some(MustUUID("84d835bf-ccca-42ca-90aa-2207372f33dd")),
	}
	var gojaValue goja.Value
	err := convert(rm, node, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("node", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		if (!("uuid" in node)) throw "expected uuid"
		if (node.uuid !== "84d835bf-ccca-42ca-90aa-2207372f33dd") throw "unexpected uuid"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &node)
	require.NoError(t, err)

	require.Equal(t, MustUUID("84d835bf-ccca-42ca-90aa-2207372f33dd"), node.UUID.Get())
}
func TestGetNodeUUIDNull(t *testing.T) {
	rm := goja.New()
	node := Node{
		UUID: option.None[UUID](),
	}
	var gojaValue goja.Value
	err := convert(rm, node, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("node", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		if (!("uuid" in node)) throw "expected uuid"
		if (node.uuid !== null) throw "unexpected uuid"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &node)
	require.NoError(t, err)

	require.True(t, node.UUID.IsNone())
}
func TestGetNodeUUIDUndefined(t *testing.T) {
	rm := goja.New()
	node := Node{
		UUID: option.Option[UUID]{},
	}
	var gojaValue goja.Value
	err := convert(rm, node, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("node", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		if ("uuid" in node) throw "unexpected uuid"
		if (node.uuid !== undefined) throw "unexpected uuid"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &node)
	require.NoError(t, err)

	require.True(t, node.UUID.IsZero())
}
func TestSetNodeUUID(t *testing.T) {
	rm := goja.New()
	node := Node{}
	var gojaValue goja.Value
	err := convert(rm, node, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("node", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		node.uuid = "2bd42b00-d96d-4360-b767-ca4bd8279576"
		if (!("uuid" in node)) throw "expected uuid"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &node)
	require.NoError(t, err)

	require.Equal(t, MustUUID("2bd42b00-d96d-4360-b767-ca4bd8279576"), node.UUID.Get())
}
func TestSetNodeUUIDNull(t *testing.T) {
	rm := goja.New()
	node := Node{
		UUID: option.Some(MustUUID("094db975-1196-45f0-9e32-d96497507a2d")),
	}
	var gojaValue goja.Value
	err := convert(rm, node, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("node", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		node.uuid = null
		if (!("uuid" in node)) throw "expected uuid"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &node)
	require.NoError(t, err)

	require.True(t, node.UUID.IsNone())
}
func TestSetNodeUUIDUndefined(t *testing.T) {
	rm := goja.New()
	node := Node{
		UUID: option.Some(MustUUID("094db975-1196-45f0-9e32-d96497507a2d")),
	}
	var gojaValue goja.Value
	err := convert(rm, node, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("node", gojaValue)
	require.NoError(t, err)
	_, err = rm.RunString(`
		node.uuid = undefined
		if ("uuid" in node) throw "unexpected uuid"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &node)
	require.NoError(t, err)

	require.True(t, node.UUID.IsZero())
}

func TestGetNodeMeta(t *testing.T) {
	rm := goja.New()
	node := Node{
		Meta: option.Some(Meta{}),
	}
	var gojaValue goja.Value
	err := convert(rm, node, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("node", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		if (!("meta" in node)) throw "expected meta"
		if (typeof node.meta !== "object") throw "unexpected meta"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &node)
	require.NoError(t, err)

	require.Equal(t, Meta{}, node.Meta.Get())
}
func TestGetNodeMetaNull(t *testing.T) {
	rm := goja.New()
	node := Node{
		Meta: option.None[Meta](),
	}
	var gojaValue goja.Value
	err := convert(rm, node, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("node", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		if (!("meta" in node)) throw "expected meta"
		if (node.meta !== null) throw "unexpected meta"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &node)
	require.NoError(t, err)

	require.True(t, node.Meta.IsNone())
}
func TestGetNodeMetaUndefined(t *testing.T) {
	rm := goja.New()
	node := Node{
		Meta: option.Option[Meta]{},
	}
	var gojaValue goja.Value
	err := convert(rm, node, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("node", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		if ("meta" in node) throw "unexpected meta"
		if (node.meta !== undefined) throw "unexpected meta"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &node)
	require.NoError(t, err)

	require.True(t, node.Meta.IsZero())
}
func TestSetNodeMeta(t *testing.T) {
	rm := goja.New()
	node := Node{}
	var gojaValue goja.Value
	err := convert(rm, node, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("node", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		node.meta = {}
		if (!("meta" in node)) throw "expected meta"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &node)
	require.NoError(t, err)

	require.Equal(t, Meta{}, node.Meta.Get())
}
func TestSetNodeMetaNull(t *testing.T) {
	rm := goja.New()
	node := Node{
		Meta: option.Some(Meta{}),
	}
	var gojaValue goja.Value
	err := convert(rm, node, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("node", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		node.meta = null
		if (!("meta" in node)) throw "expected meta"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &node)
	require.NoError(t, err)

	require.True(t, node.Meta.IsNone())
}
func TestSetNodeMetaUndefined(t *testing.T) {
	rm := goja.New()
	node := Node{
		Meta: option.Some(Meta{}),
	}
	var gojaValue goja.Value
	err := convert(rm, node, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("node", gojaValue)
	require.NoError(t, err)
	_, err = rm.RunString(`
		node.meta = undefined
		if ("meta" in node) throw "unexpected meta"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &node)
	require.NoError(t, err)

	require.True(t, node.Meta.IsZero())
}

func TestGetNodeHook(t *testing.T) {
	rm := goja.New()
	node := Node{
		Hook: option.Some(Hook{}),
	}
	var gojaValue goja.Value
	err := convert(rm, node, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("node", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		if (!("hook" in node)) throw "expected hook"
		if (typeof node.hook !== "object") throw "unexpected hook"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &node)
	require.NoError(t, err)

	require.Equal(t, Hook{}, node.Hook.Get())
}
func TestGetNodeHookNull(t *testing.T) {
	rm := goja.New()
	node := Node{
		Hook: option.None[Hook](),
	}
	var gojaValue goja.Value
	err := convert(rm, node, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("node", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		if (!("hook" in node)) throw "expected hook"
		if (node.hook !== null) throw "unexpected hook"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &node)
	require.NoError(t, err)

	require.True(t, node.Hook.IsNone())
}
func TestGetFlowHookUndefined(t *testing.T) {
	rm := goja.New()
	node := Node{
		Hook: option.Option[Hook]{},
	}
	var gojaValue goja.Value
	err := convert(rm, node, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("node", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		if ("hook" in node) throw "unexpected hook"
		if (node.hook !== undefined) throw "unexpected hook"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &node)
	require.NoError(t, err)

	require.True(t, node.Hook.IsZero())
}
func TestSetNodeHook(t *testing.T) {
	rm := goja.New()
	node := Node{}
	var gojaValue goja.Value
	err := convert(rm, node, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("node", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		node.hook = {}
		if (!("hook" in node)) throw "expected hook"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &node)
	require.NoError(t, err)

	require.Equal(t, Hook{}, node.Hook.Get())
}
func TestSetNodeHookNull(t *testing.T) {
	rm := goja.New()
	node := Node{
		Hook: option.Some(Hook{}),
	}
	var gojaValue goja.Value
	err := convert(rm, node, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("node", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		node.hook = null
		if (!("hook" in node)) throw "expected hook"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &node)
	require.NoError(t, err)

	require.True(t, node.Hook.IsNone())
}
func TestSetNodeHookUndefined(t *testing.T) {
	rm := goja.New()
	node := Node{
		Hook: option.Some(Hook{}),
	}
	var gojaValue goja.Value
	err := convert(rm, node, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("node", gojaValue)
	require.NoError(t, err)
	_, err = rm.RunString(`
		node.hook = undefined
		if ("hook" in node) throw "unexpected hook"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &node)
	require.NoError(t, err)

	require.True(t, node.Hook.IsZero())
}

func TestGetNodeLive(t *testing.T) {
	rm := goja.New()
	node := Node{
		Live: option.Some(Live{}),
	}
	var gojaValue goja.Value
	err := convert(rm, node, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("node", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		if (!("live" in node)) throw "expected live"
		if (typeof node.live !== "object") throw "unexpected live"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &node)
	require.NoError(t, err)

	require.Equal(t, Live{}, node.Live.Get())
}
func TestGetNodeLiveNull(t *testing.T) {
	rm := goja.New()
	node := Node{
		Live: option.None[Live](),
	}
	var gojaValue goja.Value
	err := convert(rm, node, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("node", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		if (!("live" in node)) throw "expected live"
		if (node.live !== null) throw "unexpected live"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &node)
	require.NoError(t, err)

	require.True(t, node.Live.IsNone())
}
func TestGetNodeLiveUndefined(t *testing.T) {
	rm := goja.New()
	node := Node{
		Live: option.Option[Live]{},
	}
	var gojaValue goja.Value
	err := convert(rm, node, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("node", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		if ("live" in node) throw "unexpected live"
		if (node.live !== undefined) throw "unexpected live"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &node)
	require.NoError(t, err)

	require.True(t, node.Live.IsZero())
}
func TestSetNodeLive(t *testing.T) {
	rm := goja.New()
	node := Node{}
	var gojaValue goja.Value
	err := convert(rm, node, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("node", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		node.live = {}
		if (!("live" in node)) throw "expected live"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &node)
	require.NoError(t, err)

	require.Equal(t, Live{}, node.Live.Get())
}
func TestSetNodeLiveNull(t *testing.T) {
	rm := goja.New()
	node := Node{
		Live: option.Some(Live{}),
	}
	var gojaValue goja.Value
	err := convert(rm, node, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("node", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		node.live = null
		if (!("live" in node)) throw "expected live"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &node)
	require.NoError(t, err)

	require.True(t, node.Live.IsNone())
}
func TestSetNodeLiveUndefined(t *testing.T) {
	rm := goja.New()
	node := Node{
		Live: option.Some(Live{}),
	}
	var gojaValue goja.Value
	err := convert(rm, node, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("node", gojaValue)
	require.NoError(t, err)
	_, err = rm.RunString(`
		node.live = undefined
		if ("live" in node) throw "unexpected live"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &node)
	require.NoError(t, err)

	require.True(t, node.Live.IsZero())
}

func TestGetLiveSince(t *testing.T) {
	rm := goja.New()
	now := time.Date(2021, 11, 16, 6, 0, 0, 0, time.UTC)
	live := Live{
		Since: option.Some(now),
	}
	var gojaValue goja.Value
	err := convert(rm, live, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("live", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		if (!("since" in live)) throw "expected live.since"
		if (live.since.getTime() != new Date("2021-11-16T06:00:00Z").getTime()) throw "unexpected live.since"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &live)
	require.NoError(t, err)

	require.WithinDuration(t, now, live.Since.Get(), time.Second)
}
func TestGetLiveSinceNull(t *testing.T) {
	rm := goja.New()
	live := Live{
		Since: option.None[time.Time](),
	}
	var gojaValue goja.Value
	err := convert(rm, live, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("live", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		if (!("since" in live)) throw "expected live.since"
		if (live.since !== null) throw "unexpected live.since"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &live)
	require.NoError(t, err)

	require.True(t, live.Since.IsNone())
}
func TestGetLiveSinceUndefined(t *testing.T) {
	rm := goja.New()
	live := Live{
		Since: option.Option[time.Time]{},
	}
	var gojaValue goja.Value
	err := convert(rm, live, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("live", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		if ("since" in live) throw "unexpected live.since"
		if (live.since !== undefined) throw "unexpected live.since"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &live)
	require.NoError(t, err)

	require.True(t, live.Since.IsZero())
}
func TestSetLiveSince(t *testing.T) {
	rm := goja.New()
	live := Live{}
	var gojaValue goja.Value
	err := convert(rm, live, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("live", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		live.since = new Date("2021-11-16T06:00:00Z")
		if (!("since" in live)) throw "expected live.since"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &live)
	require.NoError(t, err)

	now := time.Date(2021, 11, 16, 6, 0, 0, 0, time.UTC)
	require.WithinDuration(t, now, live.Since.Get(), time.Second)
}
func TestSetLiveSinceNull(t *testing.T) {
	rm := goja.New()
	now := time.Date(2021, 11, 16, 6, 0, 0, 0, time.UTC)
	live := Live{
		Since: option.Some(now),
	}
	var gojaValue goja.Value
	err := convert(rm, live, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("live", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		live.since = null
		if (!("since" in live)) throw "expected live.since"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &live)
	require.NoError(t, err)

	require.True(t, live.Since.IsNone())
}
func TestSetLiveSinceUndefined(t *testing.T) {
	rm := goja.New()
	now := time.Date(2021, 11, 16, 6, 0, 0, 0, time.UTC)
	live := Live{
		Since: option.Some(now),
	}
	var gojaValue goja.Value
	err := convert(rm, live, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("live", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		live.since = undefined
		if ("since" in live) throw "unexpected live.since"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &live)
	require.NoError(t, err)

	require.True(t, live.Since.IsZero())
}

func TestGetLiveUntil(t *testing.T) {
	rm := goja.New()
	now := time.Date(2021, 11, 16, 6, 0, 0, 0, time.UTC)
	live := Live{
		Until: option.Some(now),
	}
	var gojaValue goja.Value
	err := convert(rm, live, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("live", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		if (!("until" in live)) throw "expected live.until"
		if (live.until.getTime() != new Date("2021-11-16T06:00:00Z").getTime()) throw "unexpected live.until"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &live)
	require.NoError(t, err)

	require.WithinDuration(t, now, live.Until.Get(), time.Second)
}
func TestGetLiveUntilNull(t *testing.T) {
	rm := goja.New()
	live := Live{
		Until: option.None[time.Time](),
	}
	var gojaValue goja.Value
	err := convert(rm, live, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("live", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		if (!("until" in live)) throw "expected live.until"
		if (live.until !== null) throw "unexpected live.until"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &live)
	require.NoError(t, err)

	require.True(t, live.Until.IsNone())
}
func TestGetLiveUntilUndefined(t *testing.T) {
	rm := goja.New()
	live := Live{
		Until: option.Option[time.Time]{},
	}
	var gojaValue goja.Value
	err := convert(rm, live, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("live", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		if ("until" in live) throw "unexpected live.until"
		if (live.until !== undefined) throw "unexpected live.until"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &live)
	require.NoError(t, err)

	require.True(t, live.Until.IsZero())
}
func TestSetLiveUntil(t *testing.T) {
	rm := goja.New()
	live := Live{}
	var gojaValue goja.Value
	err := convert(rm, live, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("live", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		live.until = new Date("2021-11-16T06:00:00Z")
		if (!("until" in live)) throw "expected live.until"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &live)
	require.NoError(t, err)

	now := time.Date(2021, 11, 16, 6, 0, 0, 0, time.UTC)
	require.WithinDuration(t, now, live.Until.Get(), time.Second)
}
func TestSetLiveUntilNull(t *testing.T) {
	rm := goja.New()
	now := time.Date(2021, 11, 16, 6, 0, 0, 0, time.UTC)
	live := Live{
		Until: option.Some(now),
	}
	var gojaValue goja.Value
	err := convert(rm, live, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("live", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		live.until = null
		if (!("until" in live)) throw "expected live.until"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &live)
	require.NoError(t, err)

	require.True(t, live.Until.IsNone())
}
func TestSetLiveUntilUndefined(t *testing.T) {
	rm := goja.New()
	now := time.Date(2021, 11, 16, 6, 0, 0, 0, time.UTC)
	live := Live{
		Until: option.Some(now),
	}
	var gojaValue goja.Value
	err := convert(rm, live, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("live", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		live.until = undefined
		if ("until" in live) throw "unexpected live.until"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &live)
	require.NoError(t, err)

	require.True(t, live.Until.IsZero())
}
func TestGetMetaUndefined(t *testing.T) {
	rm := goja.New()
	meta := Meta{}
	var gojaValue goja.Value
	err := convert(rm, meta, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("meta", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		if (("value" in meta)) throw "expected meta.value"
		if (meta.value !== undefined) throw "unexpected meta.value"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &meta)
	require.NoError(t, err)
	require.NotContains(t, meta, "value")
}
func TestSetMetaUndefined(t *testing.T) {
	rm := goja.New()
	meta := Meta{}
	var gojaValue goja.Value
	err := convert(rm, meta, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("meta", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		meta.value = undefined
		if (!("value" in meta)) throw "expected meta.value"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &meta)
	require.NoError(t, err)

	require.NotContains(t, meta, "value")
}
func TestGetMetaNull(t *testing.T) {
	rm := goja.New()
	meta := Meta{
		"value": nil,
	}
	var gojaValue goja.Value
	err := convert(rm, meta, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("meta", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		if (!("value" in meta)) throw "expected meta.value"
		if (meta.value !== null) throw "unexpected meta.value"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &meta)
	require.NoError(t, err)
	require.Contains(t, meta, "value")
	require.Equal(t, nil, meta["value"])
}
func TestSetMetaNull(t *testing.T) {
	rm := goja.New()
	meta := Meta{}
	var gojaValue goja.Value
	err := convert(rm, meta, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("meta", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		meta.value = null
		if (!("value" in meta)) throw "expected meta.value"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &meta)
	require.NoError(t, err)

	require.Contains(t, meta, "value")
	require.Equal(t, nil, meta["value"])
}
func TestGetMetaString(t *testing.T) {
	rm := goja.New()
	meta := Meta{
		"value": "foo",
	}
	var gojaValue goja.Value
	err := convert(rm, meta, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("meta", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		if (!("value" in meta)) throw "expected meta.value"
		if (meta.value != "foo") throw "unexpected meta.value"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &meta)
	require.NoError(t, err)
	require.Equal(t, "foo", meta["value"])
}
func TestSetMetaString(t *testing.T) {
	rm := goja.New()
	meta := Meta{}
	var gojaValue goja.Value
	err := convert(rm, meta, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("meta", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		meta.value = "foo"
		if (!("value" in meta)) throw "expected meta.value"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &meta)
	require.NoError(t, err)

	require.Equal(t, "foo", meta["value"])
}
func TestGetMetaNumber(t *testing.T) {
	rm := goja.New()
	meta := Meta{
		"value": 5,
	}
	var gojaValue goja.Value
	err := convert(rm, meta, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("meta", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		if (!("value" in meta)) throw "expected meta.value"
		if (meta.value !== 5) throw "unexpected meta.value"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &meta)
	require.NoError(t, err)
	require.Contains(t, meta, "value")
	require.Equal(t, 5., meta["value"])
}
func TestSetMetaNumber(t *testing.T) {
	rm := goja.New()
	meta := Meta{}
	var gojaValue goja.Value
	err := convert(rm, meta, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("meta", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		meta.value = 5
		if (!("value" in meta)) throw "expected meta.value"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &meta)
	require.NoError(t, err)

	require.Equal(t, 5., meta["value"])
}
func TestGetMetaBoolean(t *testing.T) {
	rm := goja.New()
	meta := Meta{
		"value": true,
	}
	var gojaValue goja.Value
	err := convert(rm, meta, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("meta", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		if (!("value" in meta)) throw "expected meta.value"
		if (meta.value !== true) throw "unexpected meta.value"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &meta)
	require.NoError(t, err)
	require.Contains(t, meta, "value")
	require.Equal(t, true, meta["value"])
}
func TestSetMetaBoolean(t *testing.T) {
	rm := goja.New()
	meta := Meta{}
	var gojaValue goja.Value
	err := convert(rm, meta, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("meta", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		meta.value = true
		if (!("value" in meta)) throw "expected meta.value"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &meta)
	require.NoError(t, err)

	require.Equal(t, true, meta["value"])
}
func TestGetMetaArray(t *testing.T) {
	rm := goja.New()
	meta := Meta{
		"value": []any{"foo", 1, false, []any{}, map[string]any{}},
	}
	var gojaValue goja.Value
	err := convert(rm, meta, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("meta", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		if (!("value" in meta)) throw "expected meta.value"
		if (meta.value.length !== 5) throw "unexpected meta.value.length"
		if (meta.value[0] !== "foo") throw "unexpected meta.value[0]"
		if (meta.value[1] !== 1) throw "unexpected meta.value[1]"
		if (meta.value[2] !== false) throw "unexpected meta.value[2]"
		if (meta.value[3].length !== 0) throw "unexpected meta.value[3]"
		if (Object.keys(meta.value[4]).length !== 0) throw "unexpected meta.value[4]"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &meta)
	require.NoError(t, err)
	require.Contains(t, meta, "value")
	require.Equal(t, []any{"foo", 1., false, []any{}, map[string]any{}}, meta["value"])
}
func TestSetMetaArray(t *testing.T) {
	rm := goja.New()
	meta := Meta{}
	var gojaValue goja.Value
	err := convert(rm, meta, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("meta", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		meta.value = ["foo", 1, false, [], {}]
		if (!("value" in meta)) throw "expected meta.value"
		if (meta.value.length !== 5) throw "unexpected meta.value.length"
		if (meta.value[0] !== "foo") throw "unexpected meta.value[0]"
		if (meta.value[1] !== 1) throw "unexpected meta.value[1]"
		if (meta.value[2] !== false) throw "unexpected meta.value[2]"
		if (meta.value[3].length !== 0) throw "unexpected meta.value[3]"
		if (Object.keys(meta.value[4]).length !== 0) throw "unexpected meta.value[4]"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &meta)
	require.NoError(t, err)

	require.Equal(t, []any{"foo", 1., false, []any{}, map[string]any{}}, meta["value"])
}
func TestGetMetaObject(t *testing.T) {
	rm := goja.New()
	meta := Meta{
		"value": map[string]any{"a": "foo", "b": 1, "c": false, "d": []any{}, "e": map[string]any{}},
	}
	var gojaValue goja.Value
	err := convert(rm, meta, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("meta", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		if (!("value" in meta)) throw "expected meta.value"
		if (Object.keys(meta.value).length !== 5) throw "unexpected meta.value.length"
		if (meta.value.a !== "foo") throw "unexpected meta.value.a"
		if (meta.value.b !== 1) throw "unexpected meta.value.b"
		if (meta.value.c !== false) throw "unexpected meta.value.c"
		if (meta.value.d.length !== 0) throw "unexpected meta.value.d"
		if (Object.keys(meta.value.e).length !== 0) throw "unexpected meta.value.e"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &meta)
	require.NoError(t, err)
	require.Contains(t, meta, "value")
	require.Equal(t, map[string]any{"a": "foo", "b": 1., "c": false, "d": []any{}, "e": map[string]any{}}, meta["value"])
}
func TestSetMetaObject(t *testing.T) {
	rm := goja.New()
	meta := Meta{}
	var gojaValue goja.Value
	err := convert(rm, meta, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("meta", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		meta.value = {a:"foo",b:1,c:false,d:[],e:{}}
		if (!("value" in meta)) throw "expected meta.value"
		if (Object.keys(meta.value).length !== 5) throw "unexpected meta.value.length"
		if (meta.value.a !== "foo") throw "unexpected meta.value.a"
		if (meta.value.b !== 1) throw "unexpected meta.value.b"
		if (meta.value.c !== false) throw "unexpected meta.value.c"
		if (meta.value.d.length !== 0) throw "unexpected meta.value.d"
		if (Object.keys(meta.value.e).length !== 0) throw "unexpected meta.value.e"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &meta)
	require.NoError(t, err)
	require.Contains(t, meta, "value")
	require.Equal(t, map[string]any{"a": "foo", "b": 1., "c": false, "d": []any{}, "e": map[string]any{}}, meta["value"])
}
func TestSetMetaLazyList(t *testing.T) {
	rm := goja.New()
	t1 := time.Date(2021, 1, 2, 3, 4, 5, 9, time.UTC)
	t2 := time.Date(2022, 1, 2, 3, 4, 5, 9, time.UTC)
	goList := []Node{
		{
			UUID: option.Some(MustUUID("a50507cf-5015-4685-8eab-6f03f6be59e8")),
			Meta: option.Some(Meta{"a": "foo", "b": 1, "c": false, "d": []any{}, "e": map[string]any{}}),
			Hook: option.Some(Hook{"a": "foo", "b": 1, "c": false, "d": []any{}, "e": map[string]any{}}),
			Live: option.Some(Live{
				Since: option.Some(t1),
				Until: option.Some(t2),
			}),
			origin: &Node{
				UUID: option.None[UUID](),
				Meta: option.None[Meta](),
				Hook: option.None[Hook](),
				Live: option.None[Live](),
			},
		},
		{
			UUID: option.None[UUID](),
			Meta: option.None[Meta](),
			Hook: option.None[Hook](),
			Live: option.None[Live](),
		},
		{},
	}
	jsList := goja.Value(nil)
	err := convert(rm, goList, &jsList)
	require.NoError(t, err)
	require.IsType(t, (*lazyList)(nil), jsList.Export())

	goMeta := Meta{}
	jsMeta := goja.Value(nil)
	err = convert(rm, goMeta, &jsMeta)
	require.NoError(t, err)

	err = rm.Set("list", jsList)
	err = rm.Set("meta", jsMeta)
	require.NoError(t, err)

	_, err = rm.RunString(`
		meta.value = list
	`)
	require.NoError(t, err)

	err = convert(rm, jsMeta, &goMeta)
	require.NoError(t, err)
	require.Contains(t, goMeta, "value")
	require.Equal(t,
		[]any{
			map[string]any{
				"uuid": "a50507cf-5015-4685-8eab-6f03f6be59e8",
				"meta": map[string]any{"a": "foo", "b": 1., "c": false, "d": []any{}, "e": map[string]any{}},
				"hook": map[string]any{"a": "foo", "b": 1., "c": false, "d": []any{}, "e": map[string]any{}},
				"live": map[string]any{
					"since": "2021-01-02T03:04:05Z",
					"until": "2022-01-02T03:04:05Z",
				},
				"origin": map[string]any{
					"uuid": nil,
					"meta": nil,
					"hook": nil,
					"live": nil,
				},
			},
			map[string]any{
				"uuid": nil,
				"meta": nil,
				"hook": nil,
				"live": nil,
			},
			map[string]any{},
		},
		goMeta["value"],
	)
}
func TestSetMetaLazyLive(t *testing.T) {
	rm := goja.New()
	meta := Meta{}
	var gojaValue goja.Value
	err := convert(rm, meta, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("meta", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		meta.value = {a:"foo",b:1,c:false,d:[],e:{}}
		if (!("value" in meta)) throw "expected meta.value"
		if (Object.keys(meta.value).length !== 5) throw "unexpected meta.value.length"
		if (meta.value.a !== "foo") throw "unexpected meta.value.a"
		if (meta.value.b !== 1) throw "unexpected meta.value.b"
		if (meta.value.c !== false) throw "unexpected meta.value.c"
		if (meta.value.d.length !== 0) throw "unexpected meta.value.d"
		if (Object.keys(meta.value.e).length !== 0) throw "unexpected meta.value.e"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &meta)
	require.NoError(t, err)
	require.Contains(t, meta, "value")
	require.Equal(t, map[string]any{"a": "foo", "b": 1., "c": false, "d": []any{}, "e": map[string]any{}}, meta["value"])
}

// todo: the tests listed below need to be revised
func TestExportFlowMeta(t *testing.T) {
	rm := goja.New()
	src := Node{
		Meta: option.Some(Meta{
			"value": "foo",
		}),
	}
	var dst goja.Value
	err := convert(rm, src, &dst)
	require.NoError(t, err)

	dst = dst.(*goja.Object).Get("meta")
	require.IsType(t, (*goja.Object)(nil), dst)
	require.Implements(t, (*goja.DynamicObject)(nil), dst.Export())
}
func TestGetFlowMetaString(t *testing.T) {
	rm := goja.New()
	node := Node{
		Meta: option.Some(Meta{
			"value": "foo",
		}),
	}
	var dst goja.Value
	_ = convert(rm, node, &dst)
	dst = dst.(*goja.Object).Get("meta")
	require.Equal(t, rm.ToValue("foo"), dst.(*goja.Object).Get("value"))
}
func TestGetFlowMetaBoolean(t *testing.T) {
	rm := goja.New()
	src := Node{
		Meta: option.Some(Meta{
			"value": true,
		}),
	}
	var dst goja.Value
	_ = convert(rm, src, &dst)

	dst = dst.(*goja.Object).Get("meta")
	require.Equal(t, rm.ToValue(true), dst.(*goja.Object).Get("value"))
}
func TestGetFlowMetaInteger(t *testing.T) {
	rm := goja.New()
	src := Node{
		Meta: option.Some(Meta{
			"value": 1,
		}),
	}
	var dst goja.Value
	_ = convert(rm, src, &dst)

	dst = dst.(*goja.Object).Get("meta")
	require.Equal(t, rm.ToValue(1), dst.(*goja.Object).Get("value"))
}
func TestGetFlowMetaFloat(t *testing.T) {
	rm := goja.New()
	src := Node{
		Meta: option.Some(Meta{
			"value": 1.1,
		}),
	}
	var dst goja.Value
	_ = convert(rm, src, &dst)

	dst = dst.(*goja.Object).Get("meta")
	require.Equal(t, rm.ToValue(1.1), dst.(*goja.Object).Get("value"))
}
func TestExportFlowMetaArray(t *testing.T) {
	rm := goja.New()
	src := Node{
		Meta: option.Some(Meta{
			"value": []any{
				1,
				true,
				"foo",
				nil,
				map[string]any{"foo": "bar"},
			},
		}),
	}
	var dst goja.Value
	_ = convert(rm, src, &dst)

	value := dst.(*goja.Object).
		Get("meta").(*goja.Object).
		Get("value").(*goja.Object)
	require.IsType(t, (*goja.Object)(nil), value)
	require.Implements(t, (*goja.DynamicArray)(nil), value.Export())
}
func TestUnchangeFlowMetaArray(t *testing.T) {
	rm := goja.New()
	flowItem := Node{
		Meta: option.Some(Meta{
			"value": []any{
				1,
				true,
				"foo",
				nil,
				map[string]any{"foo": "bar"},
			},
		}),
	}
	var gojaValue goja.Value
	_ = convert(rm, flowItem, &gojaValue)

	err := rm.Set("item", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`item.object = true`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &flowItem)
	require.NoError(t, err)
	require.Equal(t,
		[]any{
			1,
			true,
			"foo",
			nil,
			map[string]any{"foo": "bar"},
		},
		flowItem.Meta.GetOrZero()["value"],
	)
}
func TestDeleteFlowMetaArray(t *testing.T) {
	rm := goja.New()
	flowItem := Node{
		Meta: option.Some(Meta{
			"value": []any{
				1,
				true,
				"foo",
				nil,
				map[string]any{"foo": "bar"},
			},
		}),
	}
	gojaValue := (goja.Value)(nil)
	err := convert(rm, flowItem, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("node", gojaValue)
	require.NoError(t, err)
	_, err = rm.RunString(`
		delete node.meta.value[0]
		delete node.meta.value[2]
		delete node.meta.value[4]
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &flowItem)
	require.NoError(t, err)
	require.Equal(t, []any{true, nil}, flowItem.Meta.GetOrZero()["value"])
}
func TestCreateEmptyFlowMetaArray(t *testing.T) {
	rm := goja.New()
	flowItem := Node{}
	var jsval goja.Value
	_ = convert(rm, flowItem, &jsval)
	err := rm.Set("node", jsval)
	require.NoError(t, err)

	_, err = rm.RunString(`node.meta = {value: []}`)
	require.NoError(t, err)

	err = convert(rm, rm.Get("node"), &flowItem)
	require.NoError(t, err)

	require.Equal(t, []any{}, flowItem.Meta.GetOrZero()["value"])
}
func TestCreateFlowMetaArray(t *testing.T) {
	rm := goja.New()
	inval := Meta{}
	var jsval goja.Value
	_ = convert(rm, inval, &jsval)
	err := rm.Set("meta", jsval)
	require.NoError(t, err)

	_, err = rm.RunString(`meta = {value: [1]}`)
	require.NoError(t, err)

	outval := Meta{}
	err = convert(rm, rm.Get("meta"), &outval)
	require.NoError(t, err)

	require.Equal(t, []any{1.}, outval["value"])
}
func TestPushFlowMetaArray(t *testing.T) {
	rm := goja.New()
	inval := Meta{"value": []any{1}}
	var jsval goja.Value
	_ = convert(rm, inval, &jsval)
	err := rm.Set("flow", jsval)
	require.NoError(t, err)

	_, err = rm.RunString(`flow.value.push(2)`)
	require.NoError(t, err)

	outval := Meta{}
	err = convert(rm, rm.Get("flow"), &outval)
	require.NoError(t, err)
	require.Equal(t, []any{1, 2.}, outval["value"])
}
func TestSetValueFlowMetaArray(t *testing.T) {
	rm := goja.New()
	inval := Meta{"value": []any{1}}
	var jsval goja.Value
	_ = convert(rm, inval, &jsval)
	err := rm.Set("flow", jsval)
	require.NoError(t, err)

	_, err = rm.RunString(`flow.value[0] = [true]`)
	require.NoError(t, err)

	outval := Meta{}
	err = convert(rm, rm.Get("flow"), &outval)
	require.NoError(t, err)
	require.Equal(t, []any{[]any{true}}, outval["value"])
}
func TestSetNullFlowMetaArray(t *testing.T) {
	rm := goja.New()
	inval := Meta{"value": []any{1}}
	var jsval goja.Value
	_ = convert(rm, inval, &jsval)
	err := rm.Set("flow", jsval)
	require.NoError(t, err)

	_, err = rm.RunString(`flow.value[0] = null`)
	require.NoError(t, err)

	outval := Meta{}
	err = convert(rm, rm.Get("flow"), &outval)
	require.NoError(t, err)
	require.Equal(t, []any{nil}, outval["value"])
}
func TestSetUndefinedFlowMetaArray(t *testing.T) {
	rm := goja.New()
	inval := Meta{"value": []any{1}}
	var jsval goja.Value
	_ = convert(rm, inval, &jsval)
	err := rm.Set("flow", jsval)
	require.NoError(t, err)

	_, err = rm.RunString(`flow.value[0] = undefined`)
	require.NoError(t, err)

	outval := Meta{}
	err = convert(rm, rm.Get("flow"), &outval)
	require.NoError(t, err)
	require.Equal(t, []any{}, outval["value"])
}
func TestGetFlowMetaArray(t *testing.T) {
	rm := goja.New()
	inval := Meta{"value": []any{1}}
	var jsval goja.Value
	_ = convert(rm, inval, &jsval)
	err := rm.Set("flow", jsval)
	require.NoError(t, err)

	_, err = rm.RunString(`if (flow.value[0] !== 1) throw "invalid value"`)
	require.NoError(t, err)

	outval := Meta{}
	err = convert(rm, rm.Get("flow"), &outval)
	require.NoError(t, err)
	require.Equal(t, []any{1.}, outval["value"])
}
func TestReferenceFlowMetaArray(t *testing.T) {
	rm := goja.New()
	flowItem := Node{
		Meta: option.Some(Meta{"value": []any{1}}),
	}
	var gojaValue goja.Value
	_ = convert(rm, flowItem, &gojaValue)

	err := rm.Set("node", gojaValue)
	require.NoError(t, err)
	_, err = rm.RunString(`let shared = {a: 1}`)
	require.NoError(t, err)
	_, err = rm.RunString(`node.meta.value = shared`)
	require.NoError(t, err)
	_, err = rm.RunString(`node.meta.value.a++`)
	require.NoError(t, err)
	_, err = rm.RunString(`shared.a++`)
	require.NoError(t, err)

	err = convert(rm, rm.Get("node"), &flowItem)
	require.NoError(t, err)
	require.Equal(t, map[string]any{"a": 3.}, flowItem.Meta.Get()["value"])
}
func TestLengthFlowMetaArray(t *testing.T) {
	rm := goja.New()
	flowItem := Node{
		Meta: option.Some(Meta{"value": []any{1}}),
	}
	var gojaValue goja.Value
	_ = convert(rm, flowItem, &gojaValue)

	err := rm.Set("node", gojaValue)
	require.NoError(t, err)
	_, err = rm.RunString(`if(node.meta.value.length !== 1) throw "unexpected length"`)
	require.NoError(t, err)
	_, err = rm.RunString(`node.meta.value.push(2)`)
	require.NoError(t, err)
	_, err = rm.RunString(`if(node.meta.value.length !== 2) throw "unexpected length"`)
	require.NoError(t, err)
	_, err = rm.RunString(`delete node.meta.value[0]`)
	require.NoError(t, err)
	_, err = rm.RunString(`if(node.meta.value.length !== 2) throw "unexpected length"`)
	require.NoError(t, err)
	_, err = rm.RunString(`node.meta.value.length = 3`)
	require.NoError(t, err)
	_, err = rm.RunString(`if(node.meta.value.length !== 3) throw "unexpected length"`)
	require.NoError(t, err)

	err = convert(rm, rm.Get("node"), &flowItem)
	require.NoError(t, err)
	require.Equal(t, []any{2.}, flowItem.Meta.Get()["value"])
}

func TestGetFlowHookObject(t *testing.T) {
	rm := goja.New()
	flowItem := Node{
		Hook: option.Some(Hook{
			"bool":   true,
			"number": 1,
			"object": map[string]any{"foo": []any{nil}},
		}),
	}
	var gojaValue goja.Value
	_ = convert(rm, flowItem, &gojaValue)

	err := rm.Set("node", gojaValue)
	require.NoError(t, err)
	_, err = rm.RunString(`
		if(node.hook.bool !== true) throw "unexpected value at bool"
		if(node.hook.number !== 1) throw "unexpected value at number"
		if(node.hook.object.foo[0] !== null) throw "unexpected value at object.foo[0]"
		if(node.hook.unknown !== undefined) throw "unexpected value at object.unknown"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &flowItem)
	require.NoError(t, err)
	require.Equal(t,
		Hook{
			"bool":   true,
			"number": 1.,
			"object": map[string]any{"foo": []any{nil}},
		},
		flowItem.Hook.Get(),
	)
}
func TestSetFlowHookObject(t *testing.T) {
	rm := goja.New()
	flowItem := Node{
		Hook: option.Some(Hook{
			"bool":   true,
			"number": 1,
			"object": map[string]any{"foo": []any{nil}},
		}),
	}
	var gojaValue goja.Value
	_ = convert(rm, flowItem, &gojaValue)

	err := rm.Set("node", gojaValue)
	require.NoError(t, err)
	_, err = rm.RunString(`
		node.hook.bool = false
		if(node.hook.bool !== false) throw "unexpected value at bool"
		
		node.hook.number = 2
		if(node.hook.number !== 2) throw "unexpected value at number"
		
		node.hook.object = 3
		if(node.hook.object !== 3) throw "unexpected value at object"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &flowItem)
	require.NoError(t, err)
	require.Equal(t,
		Hook{
			"bool":   false,
			"number": 2.,
			"object": 3.,
		},
		flowItem.Hook.Get(),
	)
}
func TestDeleteFlowHookObject(t *testing.T) {
	rm := goja.New()
	flowItem := Node{
		Hook: option.Some(Hook{
			"bool":   true,
			"number": 1,
			"object": map[string]any{"foo": []any{nil}},
		}),
	}
	var gojaValue goja.Value
	_ = convert(rm, flowItem, &gojaValue)

	err := rm.Set("node", gojaValue)
	require.NoError(t, err)
	_, err = rm.RunString(`
		delete node.hook.bool
		node.hook.number = undefined
		delete node.hook.object
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &flowItem)
	require.NoError(t, err)
	require.Equal(t,
		Hook{},
		flowItem.Hook.Get(),
	)
}
func TestKeysFlowHookObject(t *testing.T) {
	rm := goja.New()
	flowItem := Node{
		Hook: option.Some(Hook{
			"bool":   true,
			"number": 1,
			"object": map[string]any{"foo": []any{nil}},
		}),
	}
	var gojaValue goja.Value
	_ = convert(rm, flowItem, &gojaValue)

	err := rm.Set("node", gojaValue)
	require.NoError(t, err)
	_, err = rm.RunString(`
		node.hook.number = undefined
		node.hook.boop = 3

		const keys = Object.keys(node.hook).sort()
		if(keys.length !== 3) throw "unexpected keys length"
		if(keys[0] !== "bool") throw "unexpected keys 0"
		if(keys[1] !== "boop") throw "unexpected keys 1"
		if(keys[2] !== "object") throw "unexpected keys 2"
	`)
	require.NoError(t, err)
}
func TestGetFlowHookUnused(t *testing.T) {
	rm := goja.New()
	flowItem := Node{
		Hook: option.Some(Hook{"0": "foo", "1": false, "2": 1, "3": nil}),
	}
	var gojaValue goja.Value
	_ = convert(rm, flowItem, &gojaValue)

	err := rm.Set("node", gojaValue)
	require.NoError(t, err)
	_, err = rm.RunString(`
		node.hook["4"] = 12
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &flowItem)
	require.NoError(t, err)
	require.Equal(t, Hook{"0": "foo", "1": false, "2": 1, "3": nil, "4": 12.}, flowItem.Hook.Get())
}

func TestGetFlowList(t *testing.T) {
	rm := goja.New()
	flowList := []Node{
		{UUID: option.Some(MustUUID("086a9b1f-6519-4cd2-ab32-2218184ef863"))},
		{UUID: option.Some(MustUUID("df196843-7ba0-44d9-8c56-fadbbed6a4e3"))},
	}
	var gojaValue goja.Value
	err := convert(rm, flowList, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("node", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		if (node.length !== 2) throw "unexpected length"
		if (node[0].uuid != "086a9b1f-6519-4cd2-ab32-2218184ef863") throw "unexpected 0 uuid"
		if (node[1].uuid != "df196843-7ba0-44d9-8c56-fadbbed6a4e3") throw "unexpected 1 uuid"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &flowList)
	require.NoError(t, err)

	require.Equal(t,
		[]Node{
			{UUID: option.Some(MustUUID("086a9b1f-6519-4cd2-ab32-2218184ef863"))},
			{UUID: option.Some(MustUUID("df196843-7ba0-44d9-8c56-fadbbed6a4e3"))},
		},
		flowList,
	)
}
func TestSetFlowList(t *testing.T) {
	rm := goja.New()
	flowList := []Node{
		{UUID: option.Some(MustUUID("086a9b1f-6519-4cd2-ab32-2218184ef863"))},
		{UUID: option.Some(MustUUID("df196843-7ba0-44d9-8c56-fadbbed6a4e3"))},
		{UUID: option.Some(MustUUID("31253d1d-1002-4f63-a97c-a19de970970f"))},
	}
	var gojaValue goja.Value
	err := convert(rm, flowList, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("node", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		node[0] = null
		node[1].uuid = "0b101df0-b197-493d-9385-8eadca8d9a11"
		node[2] = undefined
		node.push({uuid:"65be679d-b665-4135-9230-bd580e815840"})
		
		if (node.length !== 4) throw "unexpected length"
		if (node[0] != null) throw "unexpected 0 uuid"
		if (node[1].uuid != "0b101df0-b197-493d-9385-8eadca8d9a11") throw "unexpected 0 uuid"
		if (node[3].uuid != "65be679d-b665-4135-9230-bd580e815840") throw "unexpected 3 uuid"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &flowList)
	require.NoError(t, err)

	require.Equal(t,
		[]Node{
			{},
			{UUID: option.Some(MustUUID("0b101df0-b197-493d-9385-8eadca8d9a11"))},
			{UUID: option.Some(MustUUID("65be679d-b665-4135-9230-bd580e815840"))},
		},
		flowList,
	)
}
func TestUnchangeFlowList(t *testing.T) {
	rm := goja.New()
	now := time.Now()
	inFlowList := []Node{
		{
			UUID: option.Some(MustUUID("fe355596-7105-419f-b766-8290c47e4988")),
			Meta: option.Some(Meta{"key": "val"}),
			Hook: option.Some(Hook{"k": "v"}),
			Live: option.Some(Live{Since: option.Some(now), Until: option.None[time.Time]()}),
		},
	}
	var gojaValue goja.Value
	_ = convert(rm, inFlowList, &gojaValue)

	outFlowList := []Node{}
	err := convert(rm, gojaValue, &outFlowList)
	require.NoError(t, err)
	require.Equal(t,
		inFlowList,
		outFlowList,
	)
}
func TestDeleteFlowList(t *testing.T) {
	rm := goja.New()
	flowList := []Node{
		{UUID: option.Some(MustUUID("0e355596-7105-419f-b766-8290c47e4988"))},
		{UUID: option.Some(MustUUID("1e355596-7105-419f-b766-8290c47e4988"))},
		{UUID: option.Some(MustUUID("2e355596-7105-419f-b766-8290c47e4988"))},
	}
	var gojaValue goja.Value
	_ = convert(rm, flowList, &gojaValue)

	err := rm.Set("flows", gojaValue)
	require.NoError(t, err)
	_, err = rm.RunString(`
		flows[0].uuid = "e0969f35-1aaa-49a6-b353-be23da2c5c57"
		delete flows[0]
		flows[1] = undefined
		delete flows[2]
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &flowList)
	require.NoError(t, err)
	require.Len(t, flowList, 0)
}
func TestUnchangeFirstFlowList(t *testing.T) {
	rm := goja.New()
	flowList := []Node{
		{UUID: option.Some(MustUUID("0e355596-7105-419f-b766-8290c47e4988"))},
		{UUID: option.Some(MustUUID("1e355596-7105-419f-b766-8290c47e4988"))},
		{UUID: option.Some(MustUUID("2e355596-7105-419f-b766-8290c47e4988"))},
	}
	var gojaValue goja.Value
	_ = convert(rm, flowList, &gojaValue)

	err := rm.Set("flows", gojaValue)
	require.NoError(t, err)
	_, err = rm.RunString(`
		flows[2].uuid = "e0969f35-1aaa-49a6-b353-be23da2c5c57"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &flowList)
	require.NoError(t, err)
	require.Equal(t,
		[]Node{
			{UUID: option.Some(MustUUID("0e355596-7105-419f-b766-8290c47e4988"))},
			{UUID: option.Some(MustUUID("1e355596-7105-419f-b766-8290c47e4988"))},
			{UUID: option.Some(MustUUID("e0969f35-1aaa-49a6-b353-be23da2c5c57"))},
		},
		flowList,
	)
}
func TestSetOriginflowFirstFlowList(t *testing.T) {
	rm := goja.New()
	flowList := []Node{
		{UUID: option.Some(MustUUID("0e355596-7105-419f-b766-8290c47e4988"))},
	}
	var gojaValue goja.Value
	_ = convert(rm, flowList, &gojaValue)

	err := rm.Set("flows", gojaValue)
	require.NoError(t, err)
	_, err = rm.RunString(`
		flows.length = 3
		flows[1]=null
		flows[2]= {uuid: "e0969f35-1aaa-49a6-b353-be23da2c5c57"}
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &flowList)
	require.NoError(t, err)
	require.Equal(t,
		[]Node{
			{UUID: option.Some(MustUUID("0e355596-7105-419f-b766-8290c47e4988"))},
			{},
			{UUID: option.Some(MustUUID("e0969f35-1aaa-49a6-b353-be23da2c5c57"))},
		},
		flowList,
	)
}
func TestGetFlowMetaTime(t *testing.T) {
	rm := goja.New()
	now := time.Date(2021, 1, 1, 17, 0, 0, 0, time.UTC)
	flowList := []Node{
		{
			UUID: option.Some(MustUUID("0e355596-7105-419f-b766-8290c47e4988")),
			Meta: option.Some(Meta{
				"time": now,
			}),
		},
	}
	var gojaValue goja.Value
	_ = convert(rm, flowList, &gojaValue)

	err := rm.Set("flows", gojaValue)
	require.NoError(t, err)
	_, err = rm.RunString(`
		if (flows[0].meta.time != "2021-01-01T17:00:00Z") 
			throw "unexpected time"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &flowList)
	require.NoError(t, err)
	require.Equal(t,
		[]Node{
			{
				UUID: option.Some(MustUUID("0e355596-7105-419f-b766-8290c47e4988")),
				Meta: option.Some(Meta{
					"time": "2021-01-01T17:00:00Z",
				}),
			},
		},
		flowList,
	)
}
func TestSetFlowMetaTime(t *testing.T) {
	rm := goja.New()
	flowList := []Node{
		{
			UUID: option.Some(MustUUID("0e355596-7105-419f-b766-8290c47e4988")),
			Meta: option.Some(Meta{}),
		},
	}
	var gojaValue goja.Value
	_ = convert(rm, flowList, &gojaValue)

	err := rm.Set("flows", gojaValue)
	require.NoError(t, err)
	_, err = rm.RunString(`
		flows[0].meta.time = new Date("2021-01-01T17:00:00Z")
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &flowList)
	require.NoError(t, err)
	require.Equal(t,
		[]Node{
			{
				UUID: option.Some(MustUUID("0e355596-7105-419f-b766-8290c47e4988")),
				Meta: option.Some(Meta{
					"time": "2021-01-01T17:00:00Z",
				}),
			},
		},
		flowList,
	)
}
func TestGetFlowMetaTimeNull(t *testing.T) {
	rm := goja.New()
	flowList := []Node{
		{
			UUID: option.Some(MustUUID("0e355596-7105-419f-b766-8290c47e4988")),
			Meta: option.Some(Meta{
				"time": nil,
			}),
		},
	}
	var gojaValue goja.Value
	_ = convert(rm, flowList, &gojaValue)

	err := rm.Set("flows", gojaValue)
	require.NoError(t, err)
	_, err = rm.RunString(`
		if (flows[0].meta.time!==null) 
			throw "unexpected time"
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &flowList)
	require.NoError(t, err)
	require.Equal(t,
		[]Node{
			{
				UUID: option.Some(MustUUID("0e355596-7105-419f-b766-8290c47e4988")),
				Meta: option.Some(Meta{
					"time": nil,
				}),
			},
		},
		flowList,
	)
}
func TestSetFlowMetaTimeNull(t *testing.T) {
	rm := goja.New()
	flowList := []Node{
		{
			UUID: option.Some(MustUUID("0e355596-7105-419f-b766-8290c47e4988")),
			Meta: option.Some(Meta{
				"time": time.Date(2021, 1, 1, 17, 0, 0, 0, time.UTC),
			}),
		},
	}
	var gojaValue goja.Value
	_ = convert(rm, flowList, &gojaValue)

	err := rm.Set("flows", gojaValue)
	require.NoError(t, err)
	_, err = rm.RunString(`
		flows[0].meta.time=null
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &flowList)
	require.NoError(t, err)
	require.Equal(t,
		[]Node{
			{
				UUID: option.Some(MustUUID("0e355596-7105-419f-b766-8290c47e4988")),
				Meta: option.Some(Meta{
					"time": nil,
				}),
			},
		},
		flowList,
	)
}
func TestSetFlowMetaTimeInvalidDate(t *testing.T) {
	rm := goja.New()

	flowItem := Node{
		Meta: option.Some(Meta{
			"time": time.Now(),
		}),
	}

	var gojaValue goja.Value
	err := convert(rm, flowItem, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("node", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		node.meta.time = new Date("invalid")
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &flowItem)
	require.NoError(t, err)
	require.Equal(t,
		Node{
			Meta: option.Some(Meta{
				"time": nil,
			}),
		},
		flowItem,
	)
}
func TestGetFlowOrigin(t *testing.T) {
	rm := goja.New()

	flowItem := Node{
		UUID: option.Some(MustUUID("9e2e7f50-9885-4fcf-b78e-804c8a6d8740")),
		Meta: option.Some(Meta{
			"a": []any{
				map[string]any{
					"b": 1,
				},
			},
		}),
		Hook: option.Some(Hook{
			"a": []any{
				map[string]any{
					"b": 1,
				},
			},
		}),
	}
	flowItem.SetOrigin(flowItem.Copy())

	var gojaValue goja.Value
	err := convert(rm, flowItem, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("node", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		node.uuid="d04e0183-e84a-427f-b14e-a7523e5885d0"
		if (node.origin.uuid !== "9e2e7f50-9885-4fcf-b78e-804c8a6d8740") 
			throw "unexpected uuid"
		
		node.meta.a[0].b=2
		if (node.origin.meta.a[0].b !== 1) 
			throw "unexpected meta"
		
		node.hook.a[0].b=2
		if (node.origin.hook.a[0].b !== 1) 
			throw "unexpected hook"
	`)
	require.NoError(t, err)
}
func TestSetFlowOrigin(t *testing.T) {
	rm := goja.New()

	flowItem := Node{
		UUID: option.Some(MustUUID("9e2e7f50-9885-4fcf-b78e-804c8a6d8740")),
		Meta: option.Some(Meta{
			"a": []any{
				map[string]any{
					"b": 1,
				},
			},
		}),
		Hook: option.Some(Hook{
			"a": []any{
				map[string]any{
					"b": 1,
				},
			},
		}),
	}
	flowItem.SetOrigin(flowItem.Copy())

	var gojaValue goja.Value
	err := convert(rm, flowItem, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("node", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		node.origin.uuid="d04e0183-e84a-427f-b14e-a7523e5885d0"
		node.origin.meta.a[0].b="2"
		node.origin.meta.a[0].c=true
		node.origin.hook.a[0].b="2"
		node.origin.hook.a[0].c=false
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &flowItem)
	require.NoError(t, err)

	expected := Node{
		UUID: option.Some(MustUUID("9e2e7f50-9885-4fcf-b78e-804c8a6d8740")),
		Meta: option.Some(Meta{
			"a": []any{
				map[string]any{
					"b": 1,
				},
			},
		}),
		Hook: option.Some(Hook{
			"a": []any{
				map[string]any{
					"b": 1,
				},
			},
		}),
	}
	expected.SetOrigin(Node{
		UUID: option.Some(MustUUID("d04e0183-e84a-427f-b14e-a7523e5885d0")),
		Meta: option.Some(Meta{
			"a": []any{
				map[string]any{
					"b": "2",
					"c": true,
				},
			},
		}),
		Hook: option.Some(Hook{
			"a": []any{
				map[string]any{
					"b": "2",
					"c": false,
				},
			},
		}),
	})
	require.Equal(t,
		expected,
		flowItem,
	)
}
func TestSetFlowOriginNull(t *testing.T) {
	rm := goja.New()

	flowItem := Node{
		UUID: option.Some(MustUUID("9e2e7f50-9885-4fcf-b78e-804c8a6d8740")),
		Meta: option.Some(Meta{
			"a": []any{
				map[string]any{
					"b": 1,
				},
			},
		}),
		Hook: option.Some(Hook{
			"a": []any{
				map[string]any{
					"b": 1,
				},
			},
		}),
	}
	flowItem.SetOrigin(flowItem.Copy())

	var gojaValue goja.Value
	err := convert(rm, flowItem, &gojaValue)
	require.NoError(t, err)

	err = rm.Set("node", gojaValue)
	require.NoError(t, err)

	_, err = rm.RunString(`
		node.origin = null
	`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &flowItem)
	require.NoError(t, err)

	require.Equal(t,
		Node{
			UUID: option.Some(MustUUID("9e2e7f50-9885-4fcf-b78e-804c8a6d8740")),
			Meta: option.Some(Meta{
				"a": []any{
					map[string]any{
						"b": 1,
					},
				},
			}),
			Hook: option.Some(Hook{
				"a": []any{
					map[string]any{
						"b": 1,
					},
				},
			}),
		},
		flowItem,
	)
}
func TestGetMetaWithLive(t *testing.T) {
	rm := goja.New()
	goMeta := Meta{}
	jsMeta := goja.Value(nil)
	err := convert(rm, goMeta, &jsMeta)
	require.NoError(t, err)
	err = rm.Set("meta", jsMeta)
	require.NoError(t, err)

	goLiveSome := Live{Since: option.Some(time.Date(2021, 2, 3, 4, 5, 6, 0, time.UTC)), Until: option.Some(time.Date(2022, 2, 3, 4, 5, 6, 0, time.UTC))}
	jsLiveSome := goja.Value(nil)
	err = convert(rm, goLiveSome, &jsLiveSome)
	require.NoError(t, err)
	err = rm.Set("liveSome", jsLiveSome)
	require.NoError(t, err)
	_, err = rm.RunString(`meta.liveSome = liveSome`)
	require.NoError(t, err)

	goLiveNone := Live{Since: option.None[time.Time](), Until: option.None[time.Time]()}
	jsLiveNone := goja.Value(nil)
	err = convert(rm, goLiveNone, &jsLiveNone)
	require.NoError(t, err)
	err = rm.Set("liveNone", jsLiveNone)
	require.NoError(t, err)
	_, err = rm.RunString(`meta.liveNone = liveNone`)
	require.NoError(t, err)

	goLiveZero := Live{}
	jsLiveZero := goja.Value(nil)
	err = convert(rm, goLiveZero, &jsLiveZero)
	require.NoError(t, err)
	err = rm.Set("liveZero", jsLiveZero)
	require.NoError(t, err)
	_, err = rm.RunString(`meta.liveZero = liveZero`)
	require.NoError(t, err)

	err = convert(rm, jsMeta, &goMeta)
	require.NoError(t, err)
	require.Equal(t,
		Meta{
			"liveSome": map[string]any{"since": "2021-02-03T04:05:06Z", "until": "2022-02-03T04:05:06Z"},
			"liveNone": map[string]any{"since": nil, "until": nil},
			"liveZero": map[string]any{},
		},
		goMeta,
	)
}

func TestSetMetaWithLive(t *testing.T) {
	rm := goja.New()
	goMeta := Meta{}
	jsMeta := goja.Value(nil)
	err := convert(rm, goMeta, &jsMeta)
	require.NoError(t, err)
	err = rm.Set("meta", jsMeta)
	require.NoError(t, err)

	goLiveSome := Live{Since: option.Some(time.Date(2021, 2, 3, 4, 5, 6, 0, time.UTC)), Until: option.Some(time.Date(2022, 2, 3, 4, 5, 6, 0, time.UTC))}
	jsLiveSome := goja.Value(nil)
	err = convert(rm, goLiveSome, &jsLiveSome)
	require.NoError(t, err)
	err = rm.Set("liveSome", jsLiveSome)
	require.NoError(t, err)
	_, err = rm.RunString(`meta.liveSome = liveSome`)
	require.NoError(t, err)
	_, err = rm.RunString(`liveSome.since=undefined; delete liveSome.until`)
	require.NoError(t, err)

	goLiveNone := Live{Since: option.None[time.Time](), Until: option.None[time.Time]()}
	jsLiveNone := goja.Value(nil)
	err = convert(rm, goLiveNone, &jsLiveNone)
	require.NoError(t, err)
	err = rm.Set("liveNone", jsLiveNone)
	require.NoError(t, err)
	_, err = rm.RunString(`meta.liveNone = liveNone`)
	require.NoError(t, err)
	_, err = rm.RunString(`liveNone.since=new Date("2021-02-03T04:05:06Z"); liveNone.until=new Date("2022-02-03T04:05:06Z")`)
	require.NoError(t, err)

	goLiveZero := Live{}
	jsLiveZero := goja.Value(nil)
	err = convert(rm, goLiveZero, &jsLiveZero)
	require.NoError(t, err)
	err = rm.Set("liveZero", jsLiveZero)
	require.NoError(t, err)
	_, err = rm.RunString(`meta.liveZero = liveZero`)
	require.NoError(t, err)
	_, err = rm.RunString(`liveZero.since=null; liveZero.until=null`)
	require.NoError(t, err)

	err = convert(rm, jsMeta, &goMeta)
	require.NoError(t, err)
	require.Equal(t,
		Meta{
			"liveSome": map[string]any{},
			"liveNone": map[string]any{"since": "2021-02-03T04:05:06Z", "until": "2022-02-03T04:05:06Z"},
			"liveZero": map[string]any{"since": nil, "until": nil},
		},
		goMeta,
	)
}
func TestGetMetaWithFlowMeta(t *testing.T) {
	rm := goja.New()
	goMeta := Meta{}
	jsMeta := goja.Value(nil)
	err := convert(rm, goMeta, &jsMeta)
	require.NoError(t, err)
	err = rm.Set("meta", jsMeta)
	require.NoError(t, err)

	goFlowMeta := Meta{"foo": 1.}
	jsFlowMeta := goja.Value(nil)
	err = convert(rm, goFlowMeta, &jsFlowMeta)
	require.NoError(t, err)
	err = rm.Set("flowMeta", jsFlowMeta)
	require.NoError(t, err)
	_, err = rm.RunString(`meta.flowMeta = flowMeta`)
	require.NoError(t, err)

	err = convert(rm, jsMeta, &goMeta)
	require.NoError(t, err)
	require.Equal(t,
		Meta{
			"flowMeta": map[string]any{"foo": 1.},
		},
		goMeta,
	)
}
func TestSetMetaWithFlowMeta(t *testing.T) {
	rm := goja.New()
	goMeta := Meta{}
	jsMeta := goja.Value(nil)
	err := convert(rm, goMeta, &jsMeta)
	require.NoError(t, err)
	err = rm.Set("meta", jsMeta)
	require.NoError(t, err)

	goFlowMeta := Meta{"foo": 1}
	jsFlowMeta := goja.Value(nil)
	err = convert(rm, goFlowMeta, &jsFlowMeta)
	require.NoError(t, err)
	err = rm.Set("flowMeta", jsFlowMeta)
	require.NoError(t, err)
	_, err = rm.RunString(`meta.flowMeta = flowMeta`)
	require.NoError(t, err)
	_, err = rm.RunString(`flowMeta.foo=1;flowMeta.bar=2`)
	require.NoError(t, err)

	err = convert(rm, jsMeta, &goMeta)
	require.NoError(t, err)
	require.Equal(t,
		Meta{
			"flowMeta": map[string]any{"foo": 1., "bar": 2.},
		},
		goMeta,
	)
}
func TestGetMetaWithFlowHook(t *testing.T) {
	rm := goja.New()
	goMeta := Meta{}
	jsMeta := goja.Value(nil)
	err := convert(rm, goMeta, &jsMeta)
	require.NoError(t, err)
	err = rm.Set("meta", jsMeta)
	require.NoError(t, err)

	goFlowHook := Hook{"foo": []any{1.}}
	jsFlowHook := goja.Value(nil)
	err = convert(rm, goFlowHook, &jsFlowHook)
	require.NoError(t, err)
	err = rm.Set("flowHook", jsFlowHook)
	require.NoError(t, err)
	_, err = rm.RunString(`meta.flowHook = flowHook`)
	require.NoError(t, err)

	err = convert(rm, jsMeta, &goMeta)
	require.NoError(t, err)
	require.Equal(t,
		Meta{
			"flowHook": map[string]any{"foo": []any{1.}},
		},
		goMeta,
	)
}
func TestSetMetaWithFlowHook(t *testing.T) {
	rm := goja.New()
	goMeta := Meta{}
	jsMeta := goja.Value(nil)
	err := convert(rm, goMeta, &jsMeta)
	require.NoError(t, err)
	err = rm.Set("meta", jsMeta)
	require.NoError(t, err)

	goFlowHook := Hook{"foo": []any{1}}
	jsFlowHook := goja.Value(nil)
	err = convert(rm, goFlowHook, &jsFlowHook)
	require.NoError(t, err)
	err = rm.Set("flowHook", jsFlowHook)
	require.NoError(t, err)
	_, err = rm.RunString(`meta.flowHook = flowHook`)
	require.NoError(t, err)
	_, err = rm.RunString(`flowHook.foo[0]=2;flowHook.bar=2`)
	require.NoError(t, err)

	err = convert(rm, jsMeta, &goMeta)
	require.NoError(t, err)
	require.Equal(t,
		Meta{
			"flowHook": map[string]any{"foo": []any{2.}, "bar": 2.},
		},
		goMeta,
	)
}
func TestSetMetaWithFlowListProto(t *testing.T) {
	rm := goja.New()
	goMeta := Meta{}
	jsMeta := goja.Value(nil)
	err := convert(rm, goMeta, &jsMeta)
	require.NoError(t, err)
	err = rm.Set("meta", jsMeta)
	require.NoError(t, err)

	goFlowList := []Node{
		{UUID: option.Some(MustUUID("1bbab2c0-38d7-4b25-81ec-ebca2291ce25"))},
		{UUID: option.Some(MustUUID("2bbab2c0-38d7-4b25-81ec-ebca2291ce25")), Hook: option.Some(Hook{"a": 1})},
		{UUID: option.Some(MustUUID("3bbab2c0-38d7-4b25-81ec-ebca2291ce25")), Meta: option.Some(Meta{"a": 1})},
		{UUID: option.Some(MustUUID("5bbab2c0-38d7-4b25-81ec-ebca2291ce25")),
			Live: option.Some(Live{Since: option.Some(time.Date(2021, 1, 1, 1, 1, 1, 0, time.UTC)), Until: option.None[time.Time]()})},
		{UUID: option.Some(MustUUID("6bbab2c0-38d7-4b25-81ec-ebca2291ce25")),
			Live: option.None[Live](), Hook: option.None[Hook](), Meta: option.None[Meta]()},
	}
	jsFlowList := goja.Value(nil)
	err = convert(rm, goFlowList, &jsFlowList)
	require.NoError(t, err)
	err = rm.Set("flowList", jsFlowList)
	require.NoError(t, err)
	_, err = rm.RunString(`meta.flowList = flowList`)
	require.NoError(t, err)

	err = convert(rm, jsMeta, &goMeta)
	require.NoError(t, err)
	require.Equal(t,
		Meta{
			"flowList": []any{
				map[string]any{"uuid": "1bbab2c0-38d7-4b25-81ec-ebca2291ce25"},
				map[string]any{"uuid": "2bbab2c0-38d7-4b25-81ec-ebca2291ce25", "hook": map[string]any{"a": 1.}},
				map[string]any{"uuid": "3bbab2c0-38d7-4b25-81ec-ebca2291ce25", "meta": map[string]any{"a": 1.}},
				map[string]any{"uuid": "5bbab2c0-38d7-4b25-81ec-ebca2291ce25", "live": map[string]any{"since": "2021-01-01T01:01:01Z", "until": nil}},
				map[string]any{"uuid": "6bbab2c0-38d7-4b25-81ec-ebca2291ce25", "live": nil, "hook": nil, "meta": nil},
			},
		},
		goMeta,
	)
}
func TestSetMetaWithFlowList(t *testing.T) {
	rm := goja.New()
	goMeta := Meta{}
	jsMeta := goja.Value(nil)
	err := convert(rm, goMeta, &jsMeta)
	require.NoError(t, err)
	err = rm.Set("meta", jsMeta)
	require.NoError(t, err)

	goFlowList := []Node{
		{UUID: option.Some(MustUUID("1bbab2c0-38d7-4b25-81ec-ebca2291ce25"))},
		{UUID: option.Some(MustUUID("2bbab2c0-38d7-4b25-81ec-ebca2291ce25")), Hook: option.Some(Hook{"a": 1})},
		{UUID: option.Some(MustUUID("3bbab2c0-38d7-4b25-81ec-ebca2291ce25")), Meta: option.Some(Meta{"a": 1})},
		{UUID: option.Some(MustUUID("5bbab2c0-38d7-4b25-81ec-ebca2291ce25")),
			Live: option.Some(Live{Since: option.Some(time.Date(2021, 1, 1, 1, 1, 1, 0, time.UTC)), Until: option.None[time.Time]()})},
		{UUID: option.Some(MustUUID("6bbab2c0-38d7-4b25-81ec-ebca2291ce25")),
			Live: option.None[Live](), Hook: option.None[Hook](), Meta: option.None[Meta]()},
	}
	jsFlowList := goja.Value(nil)
	err = convert(rm, goFlowList, &jsFlowList)
	require.NoError(t, err)
	err = rm.Set("flowList", jsFlowList)
	require.NoError(t, err)
	_, err = rm.RunString(`meta.flowList = flowList`)
	require.NoError(t, err)
	_, err = rm.RunString(`flowList.push({uuid:"1524cfa4-12a3-4061-820f-6ff0babc44bf"})`)
	require.NoError(t, err)
	_, err = rm.RunString(`flowList.push({uuid:"2524cfa4-12a3-4061-820f-6ff0babc44bf",meta:{a:2}})`)
	require.NoError(t, err)
	_, err = rm.RunString(`flowList.push({uuid:"3524cfa4-12a3-4061-820f-6ff0babc44bf",hook:{a:2}})`)
	require.NoError(t, err)
	_, err = rm.RunString(`flowList.push({uuid:"5524cfa4-12a3-4061-820f-6ff0babc44bf",live:{since:null,until:new Date("2019-03-04T01:01:01Z")}})`)
	require.NoError(t, err)
	_, err = rm.RunString(`flowList.push({uuid:"6524cfa4-12a3-4061-820f-6ff0babc44bf",live:null,hook:null,meta:null})`)
	require.NoError(t, err)

	err = convert(rm, jsMeta, &goMeta)
	require.NoError(t, err)
	require.Equal(t,
		Meta{
			"flowList": []any{
				map[string]any{"uuid": "1bbab2c0-38d7-4b25-81ec-ebca2291ce25"},
				map[string]any{"uuid": "2bbab2c0-38d7-4b25-81ec-ebca2291ce25", "hook": map[string]any{"a": 1.}},
				map[string]any{"uuid": "3bbab2c0-38d7-4b25-81ec-ebca2291ce25", "meta": map[string]any{"a": 1.}},
				map[string]any{"uuid": "5bbab2c0-38d7-4b25-81ec-ebca2291ce25", "live": map[string]any{"since": "2021-01-01T01:01:01Z", "until": nil}},
				map[string]any{"uuid": "6bbab2c0-38d7-4b25-81ec-ebca2291ce25", "live": nil, "hook": nil, "meta": nil},
				map[string]any{"uuid": "1524cfa4-12a3-4061-820f-6ff0babc44bf"},
				map[string]any{"uuid": "2524cfa4-12a3-4061-820f-6ff0babc44bf", "meta": map[string]any{"a": 2.}},
				map[string]any{"uuid": "3524cfa4-12a3-4061-820f-6ff0babc44bf", "hook": map[string]any{"a": 2.}},
				map[string]any{"uuid": "5524cfa4-12a3-4061-820f-6ff0babc44bf", "live": map[string]any{"since": nil, "until": "2019-03-04T01:01:01Z"}},
				map[string]any{"uuid": "6524cfa4-12a3-4061-820f-6ff0babc44bf", "live": nil, "hook": nil, "meta": nil},
			},
		},
		goMeta,
	)
}
func TestSetArrayLenToZero(t *testing.T) {
	rm := goja.New()
	array := any([]any{1, 2, 3})
	var gojaValue goja.Value
	_ = convert(rm, array, &gojaValue)

	err := rm.Set("array", gojaValue)
	require.NoError(t, err)
	_, err = rm.RunString(`array.length=0`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &array)
	require.NoError(t, err)
	require.Equal(t,
		[]any{},
		array,
	)
}
func TestSetArrayLenToHalf(t *testing.T) {
	rm := goja.New()
	array := any([]any{1, 2, 3})
	var gojaValue goja.Value
	_ = convert(rm, array, &gojaValue)

	err := rm.Set("array", gojaValue)
	require.NoError(t, err)
	_, err = rm.RunString(`array.length=2`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &array)
	require.NoError(t, err)
	require.Equal(t,
		[]any{1, 2},
		array,
	)
}
func TestSetArrayLenToOrigin(t *testing.T) {
	rm := goja.New()
	array := any([]any{1, 2, 3})
	var gojaValue goja.Value
	_ = convert(rm, array, &gojaValue)

	err := rm.Set("array", gojaValue)
	require.NoError(t, err)
	_, err = rm.RunString(`array.length=4`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &array)
	require.NoError(t, err)
	require.Equal(t,
		[]any{1, 2, 3},
		array,
	)
}
func TestSetArrayLenUnset(t *testing.T) {
	rm := goja.New()
	array := any([]any{1, 2, 3})
	var gojaValue goja.Value
	_ = convert(rm, array, &gojaValue)

	err := rm.Set("array", gojaValue)
	require.NoError(t, err)
	_, err = rm.RunString(`array.length=1`)
	require.NoError(t, err)
	_, err = rm.RunString(`array.length=4`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &array)
	require.NoError(t, err)
	require.Equal(t,
		[]any{1},
		array,
	)
}
func TestSetFlowListLenToZero(t *testing.T) {
	rm := goja.New()
	flowList := []Node{{UUID: option.Some(MustUUID("fa589db7-0347-4ece-b2d6-609d7d5b5d9c"))}}
	var gojaValue goja.Value
	_ = convert(rm, flowList, &gojaValue)

	err := rm.Set("flow", gojaValue)
	require.NoError(t, err)
	_, err = rm.RunString(`flow.length=0`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &flowList)
	require.NoError(t, err)
	require.Equal(t,
		[]Node{},
		flowList,
	)
}
func TestSetFlowListLenToHalf(t *testing.T) {
	rm := goja.New()
	flowList := []Node{
		{UUID: option.Some(MustUUID("fa589db7-0347-4ece-b2d6-609d7d5b5d9c"))},
		{UUID: option.Some(MustUUID("d890e872-4655-4311-9f7c-31ecfc0fcc47"))},
		{UUID: option.Some(MustUUID("7bfc69d8-9947-4f01-8d56-cd0cccada87a"))},
	}
	var gojaValue goja.Value
	_ = convert(rm, flowList, &gojaValue)

	err := rm.Set("flow", gojaValue)
	require.NoError(t, err)
	_, err = rm.RunString(`flow.length=2`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &flowList)
	require.NoError(t, err)
	require.Equal(t,
		[]Node{
			{UUID: option.Some(MustUUID("fa589db7-0347-4ece-b2d6-609d7d5b5d9c"))},
			{UUID: option.Some(MustUUID("d890e872-4655-4311-9f7c-31ecfc0fcc47"))},
		},
		flowList,
	)
}
func TestSetFlowListLenToOrigin(t *testing.T) {
	rm := goja.New()
	flowList := []Node{
		{UUID: option.Some(MustUUID("fa589db7-0347-4ece-b2d6-609d7d5b5d9c"))},
		{UUID: option.Some(MustUUID("d890e872-4655-4311-9f7c-31ecfc0fcc47"))},
		{UUID: option.Some(MustUUID("7bfc69d8-9947-4f01-8d56-cd0cccada87a"))},
	}
	var gojaValue goja.Value
	_ = convert(rm, flowList, &gojaValue)

	err := rm.Set("flow", gojaValue)
	require.NoError(t, err)
	_, err = rm.RunString(`flow.length=4`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &flowList)
	require.NoError(t, err)
	require.Equal(t,
		[]Node{
			{UUID: option.Some(MustUUID("fa589db7-0347-4ece-b2d6-609d7d5b5d9c"))},
			{UUID: option.Some(MustUUID("d890e872-4655-4311-9f7c-31ecfc0fcc47"))},
			{UUID: option.Some(MustUUID("7bfc69d8-9947-4f01-8d56-cd0cccada87a"))},
		},
		flowList,
	)
}
func TestSetFlowListLenUnset(t *testing.T) {
	rm := goja.New()
	flowList := []Node{
		{UUID: option.Some(MustUUID("fa589db7-0347-4ece-b2d6-609d7d5b5d9c"))},
		{UUID: option.Some(MustUUID("d890e872-4655-4311-9f7c-31ecfc0fcc47"))},
		{UUID: option.Some(MustUUID("7bfc69d8-9947-4f01-8d56-cd0cccada87a"))},
	}
	var gojaValue goja.Value
	_ = convert(rm, flowList, &gojaValue)

	err := rm.Set("flow", gojaValue)
	require.NoError(t, err)
	_, err = rm.RunString(`flow.length=1`)
	require.NoError(t, err)
	_, err = rm.RunString(`flow.length=4`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &flowList)
	require.NoError(t, err)
	require.Equal(t,
		[]Node{
			{UUID: option.Some(MustUUID("fa589db7-0347-4ece-b2d6-609d7d5b5d9c"))},
		},
		flowList,
	)
}
func TestResetLenLazyArray(t *testing.T) {
	rm := goja.New()
	val := any([]any{1, 2, 3})
	var gojaValue goja.Value
	_ = convert(rm, val, &gojaValue)

	err := rm.Set("val", gojaValue)
	require.NoError(t, err)
	_, err = rm.RunString(`val.length=0`)
	require.NoError(t, err)
	_, err = rm.RunString(`if(val?.[0] !== undefined) throw "unexpected 0"`)
	require.NoError(t, err)
	_, err = rm.RunString(`if(val?.[1] !== undefined) throw "unexpected 1"`)
	require.NoError(t, err)
	_, err = rm.RunString(`if(val?.[2] !== undefined) throw "unexpected 2"`)
	require.NoError(t, err)
	_, err = rm.RunString(`if(val?.[3] !== undefined) throw "unexpected 3"`)
	require.NoError(t, err)
	_, err = rm.RunString(`if(val?.[4] !== undefined) throw "unexpected 4"`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &val)
	require.NoError(t, err)
	require.Equal(t,
		[]any{},
		val,
	)
}
func TestResetLenFlowList(t *testing.T) {
	rm := goja.New()
	val := []Node{
		{UUID: option.Some(MustUUID("3d7f252c-404c-43dc-902c-7c548620f8b9"))},
		{UUID: option.Some(MustUUID("4d7f252c-404c-43dc-902c-7c548620f8b9"))},
		{UUID: option.Some(MustUUID("5d7f252c-404c-43dc-902c-7c548620f8b9"))},
	}
	var gojaValue goja.Value
	_ = convert(rm, val, &gojaValue)

	err := rm.Set("val", gojaValue)
	require.NoError(t, err)
	_, err = rm.RunString(`val.length=0`)
	require.NoError(t, err)
	_, err = rm.RunString(`if(val?.[0] !== undefined) throw "unexpected 0"`)
	require.NoError(t, err)
	_, err = rm.RunString(`if(val?.[1] !== undefined) throw "unexpected 1"`)
	require.NoError(t, err)
	_, err = rm.RunString(`if(val?.[2] !== undefined) throw "unexpected 2"`)
	require.NoError(t, err)
	_, err = rm.RunString(`if(val?.[3] !== undefined) throw "unexpected 3"`)
	require.NoError(t, err)
	_, err = rm.RunString(`if(val?.[4] !== undefined) throw "unexpected 4"`)
	require.NoError(t, err)

	err = convert(rm, gojaValue, &val)
	require.NoError(t, err)
	require.Equal(t,
		[]Node{},
		val,
	)
}
