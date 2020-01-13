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
func Initialize(opts ...Options) {
	dop := DefaultOption()
	for _, op := range opts {
		op(dop)
	}

	if dop.MpegName != "" {
		fftool.DefaultMpegName = dop.MpegName
	}

	if dop.ProbeName != "" {
		fftool.DefaultProbeName = dop.ProbeName
	}

	if dop.CommandPath != "" {
		fftool.DefaultCommandPath = dop.CommandPath
	}

	_factory.once.Do(func() {
		_factory.probe = fftool.NewFFProbe()
		_factory.mpeg = fftool.NewFFMpeg()
	})
}

// Mpeg ...
func Mpeg() *fftool.FFMpeg {
	return _factory.mpeg
}

// Probe ...
func Probe() *fftool.FFProbe {
	return _factory.probe
}
