package flow

import (
	"fmt"
)

type Plugin[F any] struct{ Name string }

func (it Plugin[F]) New(v F) Option {
	return optionFunc(func(c *Config) {
		if e, ok := c.Plugin[it.Name]; ok {
			if e, ok := e.(interface{ With(v F) F }); ok {
				v = e.With(v)
			} else {
				panic(fmt.Sprintf("plugin %q already enabled", it.Name))
			}
		}
		c.Plugin[it.Name] = v
	})
}
func (it Plugin[F]) Of(f Flow) (plugin F, ok bool) {
	plugin, ok = f.config.Plugin[it.Name].(F)
	return plugin, ok
}
