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
	Name string
}

// RunOptions ...
type RunOptions func(config *Config) *Config

func (ff *FFMpeg) init() error {
	if ff.cmd == nil {
		ff.cmd = NewCommand(ff.Name)
	}
	if ff.err != nil {
		return ff.err
	}
	return nil
}

// Version ...
func (ff *FFMpeg) Version() (string, error) {
	if err := ff.init(); err != nil {
		return "", err
	}

	return ff.cmd.Run("-version")
}

// Run ...
func (ff FFMpeg) Run(ctx context.Context, input string, opts ...RunOptions) (e error) {
	if err := ff.init(); err != nil {
		return Err(err, "init")
	}
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
	return &FFMpeg{
		Name: DefaultMpegName,
	}
}
