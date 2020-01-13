package fftool

import (
	"context"
	"errors"
	"github.com/goextension/log"
	"os"
	"strings"
	"sync"
)

// FFMpeg ...
type FFMpeg struct {
	err  error
	cmd  *Command
	name string
}

// RunOptions ...
type RunOptions func(config *Config) *Config

// Name ...
func (ff FFMpeg) Name() string {
	return ff.name
}

// Version ...
func (ff *FFMpeg) Version() (string, error) {
	return ff.cmd.Run("-version")
}

// Run ...
func (ff FFMpeg) Run(ctx context.Context, input string, opts ...RunOptions) (e error) {

	cfg := DefaultConfig()
	for _, opt := range opts {
		cfg = opt(cfg)
	}

	stat, e := os.Stat(cfg.ProcessPath())
	if e != nil {
		if os.IsNotExist(e) {
			_ = os.MkdirAll(cfg.ProcessPath(), 0755)
		} else {
			return Err(e, "stat")
		}
	}
	if e == nil && !stat.IsDir() {
		return errors.New("target is not dir")
	}

	e = cfg.Action()
	if e != nil {
		return Err(e, "action do")
	}
	args := outputArgs(cfg, input)

	outlog := make(chan string)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		e = ff.cmd.RunContext(ctx, args, outlog)
	}()
	for i2 := range outlog {
		ss := strings.Split(i2, "\r")
		for _, i3 := range ss {
			log.Infow("runmsg", "log", strings.TrimSpace(i3))
		}
	}
	wg.Wait()
	return e
}

// Error ...
func (ff *FFMpeg) Error() error {
	return ff.err
}

// NewFFMpeg ...
func NewFFMpeg() *FFMpeg {
	ff := &FFMpeg{
		name: DefaultMpegName,
	}
	ff.cmd = NewCommand(ff.name)
	return ff
}
