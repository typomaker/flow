package build

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/typomaker/flow"
)

func Goja(ctx context.Context, name string) (content []byte, err error) {
	var flowctx = flow.Context(ctx)
	var file fs.File
	if file, err = flowctx.FS().Open(name); err != nil {
		return nil, fmt.Errorf("build.goja: %w", err)
	}
	var b []byte
	if b, err = io.ReadAll(file); err != nil {
		return nil, fmt.Errorf("build.goja: %w", err)
	}

	var loader api.Loader
	if loader, err = matchLoader(name); err != nil {
		return nil, fmt.Errorf("build.goja: %w", err)
	}

	var r = api.Build(api.BuildOptions{
		Stdin: &api.StdinOptions{
			Contents:   string(b),
			Sourcefile: name,
			ResolveDir: ".",
			Loader:     loader,
		},
		Format:           api.FormatIIFE,
		GlobalName:       "entry",
		Bundle:           true,
		TreeShaking:      api.TreeShakingTrue,
		Sourcemap:        api.SourceMapInline,
		Target:           api.ES2020,
		PreserveSymlinks: true,
		Plugins: []api.Plugin{
			httpPlugin,
		},
	})
	if len(r.Errors) != 0 {
		var fmsg = api.FormatMessages(r.Errors, api.FormatMessagesOptions{Kind: api.ErrorMessage})
		return nil, fmt.Errorf("build.goja: %w", errors.New(strings.Join(fmsg, ";")))
	}
	if len(r.Warnings) != 0 {
		var fmsg = api.FormatMessages(r.Warnings, api.FormatMessagesOptions{Kind: api.WarningMessage})
		return nil, fmt.Errorf("build.goja: %w", errors.New(strings.Join(fmsg, ";")))
	}
	return r.OutputFiles[0].Contents, nil
}
