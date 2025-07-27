package flow

type Plugin[F any] struct {
	name string
}

func NewPlugin[F any](name string) Plugin[F] {
	return Plugin[F]{name: name}
}
func (it Plugin[F]) Set(v ...F) Option {
	return optionFunc(func(c *Config) {
		if f, ok := c.Plugin[it.name].([]F); ok {
			c.Plugin[it.name] = append(f, v...)
		} else {
			c.Plugin[it.name] = v
		}
	})
}
func (it Plugin[F]) Get(f Flow) []F {
	if v, ok := f.config.Plugin[it.name].([]F); ok {
		return v
	}
	var zero []F
	return zero
}
