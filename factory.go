package fftool

import "sync"

// Factory ...
type Factory struct {
	once  sync.Once
	probe *FFProbe
	mpeg  *FFMpeg
}

var _factory = New()
var _once sync.Once

// New ...
func New() *Factory {
	_once.Do(func() {
		if _factory == nil {
			_factory = &Factory{}
		}
	})
	return _factory
}
