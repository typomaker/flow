package flow

import (
	"context"
	"io/fs"
	"log/slog"
)

type Config struct {
	FS       fs.FS
	Logger   *slog.Logger
	Handler  Handler
	Modifier func(context.Context, Node) error
	Notifier func(context.Context, Case) error
	Plugin   map[string]any
}

func (it Config) With(o ...Option) Config {
	for i := range o {
		o[i].apply(&it)
	}
	return it
}
