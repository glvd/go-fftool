package factory

import (
	"github.com/glvd/go-fftool"
	"sync"
)

// Factory ...
type Factory struct {
	once  sync.Once
	probe *fftool.FFProbe
	mpeg  *fftool.FFMpeg
}

var _factory = new()
var _once sync.Once

// New ...
func new() *Factory {
	_once.Do(func() {
		if _factory == nil {
			_factory = &Factory{}
		}
	})
	return _factory
}

// Initialize ...
func Initialize() {

}
