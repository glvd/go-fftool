package factory

import "github.com/glvd/go-media-tool"

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
		MpegName:    tool.DefaultMpegName,
		ProbeName:   tool.DefaultProbeName,
		CommandPath: tool.DefaultCommandPath,
	}
}
