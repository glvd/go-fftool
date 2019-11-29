package fftool

import (
	"context"
)

// FFMpeg ...
type FFMpeg struct {
	config *Config
	cmd    *Command
	Name   string
}

func (ff *FFMpeg) init() {
	if ff.cmd == nil {
		ff.cmd = New(ff.Name)
	}
}

// Run ...
func (ff *FFMpeg) Run(ctx context.Context, input, output string) (e error) {
	return ff.cmd.RunContext(ctx, "", nil)
}

// NewFFMpeg ...
func NewFFMpeg(config Config) *FFMpeg {
	ff := &FFMpeg{
		config: &config,
		Name:   "ffmpeg",
	}

	return ff
}

func optimizeWithFormat() {

}
