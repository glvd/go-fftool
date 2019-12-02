package fftool

import (
	"context"
	"github.com/goextension/log"
	"sync"
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
	ff.init()
	args := ff.config.Args(input, output)

	outlog := make(chan string)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		e = ff.cmd.RunContext(ctx, args, outlog)
	}()
	for i2 := range outlog {
		log.Info("run", "info", i2)
	}
	wg.Wait()
	return e
}

// NewFFMpeg ...
func NewFFMpeg() *FFMpeg {
	ff := &FFMpeg{
		config: *DefaultConfig(),
		Name:   "ffmpeg",
	}

	return ff
}
