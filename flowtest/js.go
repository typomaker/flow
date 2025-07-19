package flowtest

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
	"testing/fstest"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/typomaker/flow"
	"github.com/typomaker/option"
)

const JSFSPath = "foo/index.js"

func TestJS(t *testing.T, build Provider) {
	for _, jsCase := range jsCases {
		t.Run(jsCase.name, func(t *testing.T) {
			jsCase.test(t, build)
		})
	}
}

var jsCases = []struct {
	name string
	test func(t *testing.T, provide Provider)
}{
	{
		name: "empty main",
		test: func(t *testing.T, provide Provider) {
			f := flow.New(
				flow.FS(fstest.MapFS{
					JSFSPath: &fstest.MapFile{
						Data: []byte(`
							export default function main() {}
						`),
					},
				}),
			)
			err := f.Run(context.Background(), nil)
			require.NoError(t, err)
		},
	},
	{
		name: "console log second argument as object with node",
		test: func(t *testing.T, provide Provider) {
			b := bytes.Buffer{}
			f := flow.New(
				flow.FS(fstest.MapFS{
					JSFSPath: &fstest.MapFile{
						Data: []byte(`
							console.log("foo", {
								node: {
									uuid: "c546e4d0-3315-4bce-bfcf-661eb6710a03",
									meta: {foo:[1,2]},
									hook: {buz:{k:"v"}},
									live: {
										since: new Date("2021-01-01T16:00:00Z"),
										until: new Date("2022-01-01T16:00:00Z")
									}
								}
							})
							export default function main() {}
						`),
					},
				}),
				flow.Logger(
					slog.New(slog.NewJSONHandler(&b, slogJsonHandlerOptions)),
				),
				provide(t),
			)

			err := f.Run(context.Background(), nil)
			require.NoError(t, err)
			require.JSONEq(
				t,
				`{
					"level":"INFO",
					"msg":"foo",
					"node":{
						"uuid":"c546e4d0-3315-4bce-bfcf-661eb6710a03",
						"meta":{
							"foo":"[1,2]"
						},
						"hook":{
							"buz":"{\"k\":\"v\"}"
						},
						"live":{
							"since":"2021-01-01T16:00:00Z",
							"until":"2022-01-01T16:00:00Z"
						}
					}
				}`,
				b.String(),
			)
		},
	},
	{
		name: "console log only message ",
		test: func(t *testing.T, provide Provider) {
			b := bytes.Buffer{}
			f := flow.New(
				flow.FS(fstest.MapFS{
					JSFSPath: &fstest.MapFile{
						Data: []byte(`
							console.log("foo")
							export default function main() {}
						`),
					},
				}),
				flow.Logger(
					slog.New(slog.NewJSONHandler(&b, slogJsonHandlerOptions)),
				),
				provide(t),
			)
			err := f.Run(context.Background(), nil)
			require.NoError(t, err)
			require.JSONEq(
				t,
				`{
					"level":"INFO",
					"msg":"foo"
				}`,
				b.String(),
			)
		},
	},
	{
		name: "console log first argument not message",
		test: func(t *testing.T, provide Provider) {
			b := bytes.Buffer{}
			f := flow.New(
				flow.FS(fstest.MapFS{
					JSFSPath: &fstest.MapFile{
						Data: []byte(`
							console.log([1, true])
							export default function main() {}
						`),
					},
				}),
				flow.Logger(
					slog.New(slog.NewJSONHandler(&b, slogJsonHandlerOptions)),
				),
				provide(t),
			)
			err := f.Run(context.Background(), nil)
			require.NoError(t, err)
			require.JSONEq(
				t,
				`{
					"level":"INFO",
					"msg":"js print",
					"js":{
						"args": "[[1,true]]"
					}
				}`,
				b.String(),
			)
		},
	},
	{
		name: "console log multiple unexpected arguments",
		test: func(t *testing.T, provide Provider) {
			b := bytes.Buffer{}
			f := flow.New(
				flow.FS(fstest.MapFS{
					JSFSPath: &fstest.MapFile{
						Data: []byte(`
							console.log([1, true], true, "buz")
							export default function main() {}
						`),
					},
				}),
				flow.Logger(
					slog.New(slog.NewJSONHandler(&b, slogJsonHandlerOptions)),
				),
				provide(t),
			)
			err := f.Run(context.Background(), nil)
			require.NoError(t, err)
			require.JSONEq(
				t,
				`{
					"level":"INFO",
					"msg":"js print",
					"js":{
						"args": "[[1,true],true,\"buz\"]"
					}
				}`,
				b.String(),
			)
		},
	},
	{
		name: "console log first argument as object",
		test: func(t *testing.T, provide Provider) {
			b := bytes.Buffer{}
			f := flow.New(
				flow.FS(fstest.MapFS{
					JSFSPath: &fstest.MapFile{
						Data: []byte(`
							console.log({"foo":"bar"}, true, "buz")
							export default function main() {}
						`),
					},
				}),
				flow.Logger(
					slog.New(slog.NewJSONHandler(&b, slogJsonHandlerOptions)),
				),
				provide(t),
			)
			err := f.Run(context.Background(), nil)
			require.NoError(t, err)
			require.JSONEq(
				t,
				`{
					"level":"INFO",
					"msg":"js print",
					"js":{
						"foo":"\"bar\"",
						"args": "[true,\"buz\"]"
					}
				}`,
				b.String(),
			)
		},
	},
	{
		name: "exception throwing",
		test: func(t *testing.T, provide Provider) {
			b := bytes.Buffer{}
			f := flow.New(
				flow.FS(fstest.MapFS{
					JSFSPath: &fstest.MapFile{
						Data: []byte(`
							export default function main(nodes) {
								throw new Error("foo")
							}
						`),
					},
				}),
				flow.Logger(
					slog.New(slog.NewJSONHandler(&b, slogJsonHandlerOptions)),
				),
				provide(t),
			)
			target := []flow.Node{
				{},
			}
			err := f.Run(context.Background(), target)
			require.ErrorContains(t, err, "foo")
		},
	},
	{
		name: "set node ",
		test: func(t *testing.T, provide Provider) {
			b := bytes.Buffer{}
			f := flow.New(
				flow.FS(fstest.MapFS{
					JSFSPath: &fstest.MapFile{
						Data: []byte(`
							export default function main(nodes) {
								nodes[0].uuid = "e14492a1-8333-4439-ab2e-d211dc305734"
								nodes[0].meta = {buz: true}
								nodes[0].hook = {zod: "true"}
								nodes[0].live = {
									since: new Date("2021-01-01T16:00:00Z"),
									until: new Date("2022-01-01T16:00:00Z")
								}
							}
						`),
					},
				}),
				flow.Logger(
					slog.New(slog.NewJSONHandler(&b, slogJsonHandlerOptions)),
				),
				provide(t),
			)
			target := []flow.Node{
				{},
			}
			err := f.Run(context.Background(), target)
			require.NoError(t, err)
			require.Equal(t, flow.MustUUID("e14492a1-8333-4439-ab2e-d211dc305734"), target[0].UUID.GetOrZero())
			require.Equal(t, flow.Meta{"buz": true}, target[0].Meta.GetOrZero())
			require.Equal(t, flow.Hook{"zod": "true"}, target[0].Hook.GetOrZero())
			require.Equal(t,
				flow.Live{
					Since: option.Some(time.Date(2021, 1, 1, 16, 0, 0, 0, time.UTC)),
					Until: option.Some(time.Date(2022, 1, 1, 16, 0, 0, 0, time.UTC)),
				},
				target[0].Live.GetOrZero(),
			)
		},
	},
	{
		name: "node uuid",
		test: func(t *testing.T, provide Provider) {
			b := bytes.Buffer{}
			f := flow.New(
				flow.FS(fstest.MapFS{
					JSFSPath: &fstest.MapFile{
						Data: []byte(`
							export default function main(nodes) {
								if (nodes[0].uuid == "e14492a1-8333-4439-ab2e-d211dc305734") {
									nodes[0].uuid = "de208cd8-b993-4861-a354-6d218a578556"
								}
							}
						`),
					},
				}),
				flow.Logger(
					slog.New(slog.NewJSONHandler(&b, slogJsonHandlerOptions)),
				),
				provide(t),
			)
			target := []flow.Node{
				{UUID: option.Some(flow.MustUUID("e14492a1-8333-4439-ab2e-d211dc305734"))},
			}
			err := f.Run(context.Background(), target)
			require.NoError(t, err)
			require.Equal(t, flow.MustUUID("de208cd8-b993-4861-a354-6d218a578556"), target[0].UUID.GetOrZero())
		},
	},
	{
		name: "node meta",
		test: func(t *testing.T, provide Provider) {
			b := bytes.Buffer{}
			f := flow.New(
				flow.FS(fstest.MapFS{
					JSFSPath: &fstest.MapFile{
						Data: []byte(`
							export default function main(nodes) {
								if (nodes[0].meta?.foo == "foo") {
									nodes[0].meta.foo = true
									nodes[0].meta.bar = "bar"
									nodes[0].meta.buz = [1]
								}
							}
						`),
					},
				}),
				flow.Logger(
					slog.New(slog.NewJSONHandler(&b, slogJsonHandlerOptions)),
				),
				provide(t),
			)
			target := []flow.Node{
				{Meta: option.Some(flow.Meta{"foo": "foo"})},
			}
			err := f.Run(context.Background(), target)
			require.NoError(t, err)
			require.Equal(t,
				flow.Meta{"foo": true, "bar": "bar", "buz": []any{1.}},
				target[0].Meta.GetOrZero(),
			)
		},
	},
	{
		name: "node hook",
		test: func(t *testing.T, provide Provider) {
			b := bytes.Buffer{}
			f := flow.New(
				flow.FS(fstest.MapFS{
					JSFSPath: &fstest.MapFile{
						Data: []byte(`
							export default function main(nodes) {
								if (nodes[0].hook?.foo == "foo") {
									nodes[0].hook.foo = true
									nodes[0].hook.bar = "bar"
									nodes[0].hook.buz = [1]
								}
							}
						`),
					},
				}),
				flow.Logger(
					slog.New(slog.NewJSONHandler(&b, slogJsonHandlerOptions)),
				),
				provide(t),
			)
			target := []flow.Node{
				{Hook: option.Some(flow.Hook{"foo": "foo"})},
			}
			err := f.Run(context.Background(), target)
			require.NoError(t, err)
			require.Equal(t,
				flow.Hook{"foo": true, "bar": "bar", "buz": []any{1.}},
				target[0].Hook.GetOrZero(),
			)
		},
	},
	{
		name: "node live",
		test: func(t *testing.T, provide Provider) {
			b := bytes.Buffer{}
			f := flow.New(
				flow.FS(fstest.MapFS{
					JSFSPath: &fstest.MapFile{
						Data: []byte(`
							export default function main(nodes) {
								if (nodes[0].live.since?.getTime() == new Date("2021-01-01T00:00:00Z").getTime()) {
									nodes[0].live.since = new Date("2022-01-01T00:00:00Z")
								}
								if (nodes[0].live.until?.getTime() == new Date("2023-01-01T00:00:00Z").getTime()) {
									nodes[0].live.until = new Date("2024-01-01T00:00:00Z")
								}
							}
						`),
					},
				}),
				flow.Logger(
					slog.New(slog.NewJSONHandler(&b, slogJsonHandlerOptions)),
				),
				provide(t),
			)
			target := []flow.Node{
				{
					Live: option.Some(flow.Live{
						Since: option.Some(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
						Until: option.Some(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
					}),
				},
			}
			err := f.Run(context.Background(), target)
			require.NoError(t, err)
			require.Equal(t,
				flow.Live{
					Since: option.Some(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
					Until: option.Some(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
				target[0].Live.GetOrZero(),
			)
		},
	},
	{
		name: "modify stub by default",
		test: func(t *testing.T, provide Provider) {
			b := bytes.Buffer{}
			f := flow.New(
				flow.FS(fstest.MapFS{
					JSFSPath: &fstest.MapFile{
						Data: []byte(`
							export default function main(nodes) {
								this.modify()
							}
						`),
					},
				}),
				flow.Logger(
					slog.New(slog.NewJSONHandler(&b, slogJsonHandlerOptions)),
				),
				provide(t),
			)
			target := []flow.Node{
				{},
			}
			err := f.Run(context.Background(), target)
			require.NoError(t, err)
		},
	},
	{
		name: "modify with js object",
		test: func(t *testing.T, provide Provider) {
			b := bytes.Buffer{}
			f := flow.New(
				flow.FS(fstest.MapFS{
					JSFSPath: &fstest.MapFile{
						Data: []byte(`
							export default function main(nodes) {
								this.modify({
									uuid: "08a0cfc4-9dd8-4869-9eec-47ab946e5da3",
									meta: {foo: true},
									hook: {buz: "bar"},
									live: {
										since: new Date("2024-01-01T00:00:00Z"),
										until: new Date("2025-01-01T00:00:00Z"),
									}
								})
							}
						`),
					},
				}),
				flow.Logger(
					slog.New(slog.NewJSONHandler(&b, slogJsonHandlerOptions)),
				),
				provide(t),
			)
			target := []flow.Node{
				{},
			}
			modify := &modifier{}
			err := f.Run(context.Background(), target, modify)
			require.NoError(t, err)
			require.Len(t, modify.flowNode, 1)
			require.Equal(t,
				flow.MustUUID("08a0cfc4-9dd8-4869-9eec-47ab946e5da3"),
				modify.flowNode[0].UUID.GetOrZero(),
			)
			require.Equal(t,
				flow.Meta{"foo": true},
				modify.flowNode[0].Meta.GetOrZero(),
			)
			require.Equal(t,
				flow.Hook{"buz": "bar"},
				modify.flowNode[0].Hook.GetOrZero(),
			)
			require.Equal(t,
				flow.Live{
					Since: option.Some(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					Until: option.Some(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
				modify.flowNode[0].Live.GetOrZero(),
			)
		},
	},
	{
		name: "notify stub by default",
		test: func(t *testing.T, provide Provider) {
			b := bytes.Buffer{}
			f := flow.New(
				flow.FS(fstest.MapFS{
					JSFSPath: &fstest.MapFile{
						Data: []byte(`
							export default function main(nodes) {
								this.notify()
							}
						`),
					},
				}),
				flow.Logger(
					slog.New(slog.NewJSONHandler(&b, slogJsonHandlerOptions)),
				),
				provide(t),
			)
			target := []flow.Node{
				{},
			}
			err := f.Run(context.Background(), target)
			require.NoError(t, err)
		},
	},
	{
		name: "notify with js object",
		test: func(t *testing.T, provide Provider) {
			b := bytes.Buffer{}
			f := flow.New(
				flow.FS(fstest.MapFS{
					JSFSPath: &fstest.MapFile{
						Data: []byte(`
							export default function main(nodes) {
								this.notify({
									when: {
										uuid: [
											"08a0cfc4-9dd8-4869-9eec-47ab946e5da3", 
											"c72e5013-f676-4f24-8df4-b49ec7113b24",
										],
										hook: [{buz:"bar"}]
									},
									then: {
										uuid: "08a0cfc4-9dd8-4869-9eec-47ab946e5da3",
										meta: {foo: true},
										hook: {buz: "bar"},
										live: {
											since: new Date("2024-01-01T00:00:00Z"),
											until: new Date("2025-01-01T00:00:00Z"),
										}
									},
								})
							}
						`),
					},
				}),
				flow.Logger(
					slog.New(slog.NewJSONHandler(&b, slogJsonHandlerOptions)),
				),
				provide(t),
			)
			target := []flow.Node{
				{},
			}
			notify := &notifier{}
			err := f.Run(context.Background(), target, notify)
			require.NoError(t, err)
			require.Len(t, notify.flowCase, 1)
			require.Equal(t,
				[]flow.UUID{
					flow.MustUUID("08a0cfc4-9dd8-4869-9eec-47ab946e5da3"),
					flow.MustUUID("c72e5013-f676-4f24-8df4-b49ec7113b24"),
				},
				notify.flowCase[0].When.UUID.GetOrZero(),
			)
			require.Equal(t,
				flow.Meta{"foo": true},
				notify.flowCase[0].Then.Meta.GetOrZero(),
			)
			require.Equal(t,
				flow.Hook{"buz": "bar"},
				notify.flowCase[0].Then.Hook.GetOrZero(),
			)
			require.Equal(t,
				flow.Live{
					Since: option.Some(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					Until: option.Some(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
				notify.flowCase[0].Then.Live.GetOrZero(),
			)
		},
	},
}

type Provider func(t *testing.T) flow.Handler

var slogJsonHandlerOptions = &slog.HandlerOptions{
	ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		if len(groups) > 0 {
			switch groups[0] {
			case "runtime":
				return slog.Attr{}
			}
		}
		switch a.Key {
		case "time":
			return slog.Attr{}
		}
		return a
	},
}

type modifier struct {
	flowNode []flow.Node
	err      error
}

var _ flow.Modifier = (*modifier)(nil)

func (m *modifier) Modify(ctx context.Context, n flow.Node) error {
	m.flowNode = append(m.flowNode, n)
	return m.err
}
func (m *modifier) LogAttr() slog.Attr {
	return slog.Group("modifier", slog.String("foo", "bar"))
}

type notifier struct {
	flowCase []flow.Case
	err      error
}

var _ flow.Notifier = (*notifier)(nil)

func (m *notifier) Notify(ctx context.Context, n flow.Case) error {
	m.flowCase = append(m.flowCase, n)
	return m.err
}
func (m *notifier) LogAttr() slog.Attr {
	return slog.Group("notifier", slog.String("foo", "bar"))
}
