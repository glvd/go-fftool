package fftool

import (
	"context"
)

// FFMpeg ...
type FFMpeg struct {
	config Config
	cmd    *Command
	Name   string
}

func (ff *FFMpeg) init() {
	if ff.cmd == nil {
		ff.cmd = New(ff.Name)
	}
}

// Version ...
func (ff *FFMpeg) Version() (string, error) {
	ff.init()
	return ff.cmd.Run("-version")
}

// OptimizeWithFormat ...
func (ff FFMpeg) OptimizeWithFormat(format *StreamFormat) (*FFMpeg, error) {
	e := ff.config.OptimizeWithFormat(format)
	if e != nil {
		return nil, e
	}
	return &ff, nil
}

// Run ...
func (ff FFMpeg) Run(ctx context.Context, input, output string) (e error) {

	ff.config.Args(input, output)

	return ff.cmd.RunContext(ctx, "", nil)
}

// NewFFMpeg ...
func NewFFMpeg() *FFMpeg {
	ff := &FFMpeg{
		config: *DefaultConfig(),
		Name:   "ffmpeg",
	}

	return ff
}
