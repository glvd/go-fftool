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

// Version ...
func (ff *FFMpeg) Version() (string, error) {
	ff.init()
	return ff.cmd.Run("-version")
}

// Run ...
func (ff FFMpeg) Run(ctx context.Context, input, output string) (e error) {
	//c := ff.config.Clone()

	ff.config.Args(input, output)

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
