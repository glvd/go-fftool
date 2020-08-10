package factory

import (
	"github.com/glvd/go-media-tool"
	"sync"
)

// Factory ...
type Factory struct {
	once  sync.Once
	probe *tool.FFProbe
	mpeg  *tool.FFMpeg
}

var _factory *Factory
var _once sync.Once

func init() {
	new()
}

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
		tool.DefaultMpegName = dop.MpegName
	}

	if dop.ProbeName != "" {
		tool.DefaultProbeName = dop.ProbeName
	}

	if dop.CommandPath != "" {
		tool.DefaultCommandPath = dop.CommandPath
	}

	_factory.once.Do(func() {
		_factory.probe = tool.NewFFProbe()
		_factory.mpeg = tool.NewFFMpeg()
	})
}

// Mpeg ...
func Mpeg() *tool.FFMpeg {
	return _factory.mpeg
}

// Probe ...
func Probe() *tool.FFProbe {
	return _factory.probe
}
