package build

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/typomaker/flow"
)

func newImportFlow(ctx context.Context) api.Plugin {
	const namespace = "import-flow"
	var flowctx = flow.Context(ctx)
	return api.Plugin{
		Name: namespace,
		Setup: func(build api.PluginBuild) {
			build.OnResolve(
				api.OnResolveOptions{Filter: `^flow:+`},
				func(args api.OnResolveArgs) (api.OnResolveResult, error) {
					return api.OnResolveResult{
						Path:      args.Path,
						Namespace: namespace,
					}, nil
				},
			)
			build.OnResolve(
				api.OnResolveOptions{Filter: ".*", Namespace: namespace},
				func(args api.OnResolveArgs) (api.OnResolveResult, error) {
					base, err := url.Parse(args.Importer)
					if err != nil {
						return api.OnResolveResult{}, err
					}
					relative, err := url.Parse(args.Path)
					if err != nil {
						return api.OnResolveResult{}, err
					}
					return api.OnResolveResult{
						Path:      base.ResolveReference(relative).String(),
						Namespace: namespace,
					}, nil
				},
			)
			build.OnLoad(
				api.OnLoadOptions{Filter: ".*", Namespace: namespace},
				func(args api.OnLoadArgs) (r api.OnLoadResult, err error) {
					var path = strings.TrimPrefix(args.Path, "flow:")
					var fsfile fs.File
					if fsfile, err = flowctx.FS().Open(path); err != nil {
						return r, fmt.Errorf("%s: %w", namespace, err)
					}
					var fsinfo fs.FileInfo
					if fsinfo, err = fsfile.Stat(); err != nil {
						return r, fmt.Errorf("%s: %w", namespace, err)
					}
					if fsinfo.IsDir() {
						return r, fmt.Errorf("%s: %s is directory", namespace, args.Path)
					}
					var content []byte
					if content, err = io.ReadAll(fsfile); err != nil {
						return r, fmt.Errorf("%s: %w", namespace, err)
					}
					var scontent = string(content)
					return api.OnLoadResult{Contents: &scontent}, nil
				},
			)
		},
	}
}
