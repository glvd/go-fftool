package factory

import "github.com/glvd/go-fftool"

// Option ...
type Option struct {
	MpegName    string
	ProbeName   string
	CommandPath string
}

// Options ...
type Options func(option *Option)

// DefaultOption ...
func DefaultOption() *Option {
	return &Option{
		MpegName:    fftool.DefaultMpegName,
		ProbeName:   fftool.DefaultProbeName,
		CommandPath: fftool.DefaultCommandPath,
	}
}
