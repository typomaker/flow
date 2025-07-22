package build

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/evanw/esbuild/pkg/api"
)

func newImportHTTP(_ context.Context) api.Plugin {
	const namespace = "import-http"
	return api.Plugin{
		Name: namespace,
		Setup: func(build api.PluginBuild) {
			build.OnResolve(
				api.OnResolveOptions{Filter: `^https?://`},
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
					var tmpname = filepath.Join(
						os.TempDir(),
						args.Namespace,
						args.Path,
						args.Suffix,
					)

					var content []byte
					if content, err = os.ReadFile(tmpname); err != nil {
						if !errors.Is(err, os.ErrNotExist) {
							return r, fmt.Errorf("%s: %w", namespace, err)
						}
						res, err := http.Get(args.Path)
						if err != nil {
							return api.OnLoadResult{}, fmt.Errorf("%s: %w", namespace, err)
						}
						defer res.Body.Close()

						if content, err = io.ReadAll(res.Body); err != nil {
							return r, fmt.Errorf("%s: %w", namespace, err)
						}
						if err = os.MkdirAll(filepath.Dir(tmpname), os.ModePerm); err != nil {
							return r, fmt.Errorf("%s: %w", namespace, err)
						}
						if err = os.WriteFile(tmpname, content, os.ModePerm); err != nil {
							return r, fmt.Errorf("%s: %w", namespace, err)
						}
					}
					var scontent = string(content)
					return api.OnLoadResult{Contents: &scontent}, nil
				},
			)
		},
	}
}
