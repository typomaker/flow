package build

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/evanw/esbuild/pkg/api"
)

const HTTPName = "esbuild-http"

var httpPlugin = api.Plugin{
	Name: HTTPName,
	Setup: func(build api.PluginBuild) {
		// Intercept import paths starting with "http:" and "https:" so
		// esbuild doesn't attempt to map them to a file system location.
		// Tag them with the HTTPName namespace to associate them with
		// this plugin.
		build.OnResolve(api.OnResolveOptions{Filter: `^https?://`},
			func(args api.OnResolveArgs) (api.OnResolveResult, error) {
				return api.OnResolveResult{
					Path:      args.Path,
					Namespace: HTTPName,
				}, nil
			})

		// We also want to intercept all import paths inside downloaded
		// files and resolve them against the original URL. All of these
		// files will be in the HTTPName namespace. Make sure to keep
		// the newly resolved URL in the HTTPName namespace so imports
		// inside it will also be resolved as URLs recursively.
		build.OnResolve(api.OnResolveOptions{Filter: ".*", Namespace: HTTPName},
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
					Namespace: HTTPName,
				}, nil
			})

		// When a URL is loaded, we want to actually download the content
		// from the internet. This has just enough logic to be able to
		// handle the example import from HTTPName but in reality this
		// would probably need to be more complex.
		build.OnLoad(api.OnLoadOptions{Filter: ".*", Namespace: HTTPName},
			func(args api.OnLoadArgs) (r api.OnLoadResult, err error) {
				var tmpname = filepath.Join(
					os.TempDir(),
					args.Namespace,
					args.Suffix,
					args.Path,
				)

				var content []byte
				if content, err = os.ReadFile(tmpname); err != nil {
					if !errors.Is(err, os.ErrNotExist) {
						return r, fmt.Errorf("%s: read temp %w", HTTPName, err)
					}
					res, err := http.Get(args.Path)
					if err != nil {
						return api.OnLoadResult{}, fmt.Errorf("%s: load data %w", HTTPName, err)
					}
					defer res.Body.Close()

					if content, err = io.ReadAll(res.Body); err != nil {
						return r, fmt.Errorf("%s: read body %w", HTTPName, err)
					}
					if err = os.MkdirAll(filepath.Dir(tmpname), os.ModePerm); err != nil {
						return r, fmt.Errorf("%s: mkdir temp %w", HTTPName, err)
					}
					if err = os.WriteFile(tmpname, content, os.ModePerm); err != nil {
						return r, fmt.Errorf("%s: write temp %w", HTTPName, err)
					}
				}
				var scontent = string(content)
				return api.OnLoadResult{Contents: &scontent}, nil
			})
	},
}
