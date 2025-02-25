package flow

import "log/slog"

func (it *Flow) setup(f *Flow) {
	f.stock = append(f.stock, it.stock...)
	f.plugin = append(f.plugin, it.plugin...)
	f.logger = it.logger
}
func (it Pipe) setup(f *Flow) {
	f.stock = append(f.stock, it)
}
func (it Plugin) setup(f *Flow) {
	f.plugin = append(f.plugin, it)
}
func Logger(v *slog.Logger) Option {
	return optionFunc(func(f *Flow) {
		f.logger = v
	})
}

type optionFunc func(*Flow)

func (it optionFunc) setup(f *Flow) {
	it(f)
}

type Option interface {
	setup(f *Flow)
}
