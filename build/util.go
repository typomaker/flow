package build

import (
	"fmt"
	"path/filepath"

	"github.com/evanw/esbuild/pkg/api"
)

func matchLoader(path string) (api.Loader, error) {
	var ext = filepath.Ext(path)
	switch ext {
	case ".js":
		return api.LoaderJS, nil
	case ".json":
		return api.LoaderJSON, nil
	case ".ts":
		return api.LoaderTS, nil
	default:
		return api.LoaderNone, fmt.Errorf("esbuild: unexpected file extension")
	}
}
