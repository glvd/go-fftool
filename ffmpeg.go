package fftool

import (
	"context"
	"errors"
	"fmt"
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

// MpegOption ...
type MpegOption struct {
	Debug  bool
	Config *Config
	Input  string
	Output string
}

// RunOptions ...
type RunOptions func(opts *MpegOption)

// Name ...
func (ff FFMpeg) Name() string {
	return ff.name
}

// Version ...
func (ff *FFMpeg) Version() (string, error) {
	return ff.cmd.Run("-version")
}

// Run ...
func (ff FFMpeg) Run(ctx context.Context, opts ...RunOptions) (e error) {
	m := &MpegOption{
		Config: DefaultConfig(),
		Input:  "",
		Output: "",
	}
	for _, opt := range opts {
		opt(m)
	}
	if m.Config.processID != "" {
		return fmt.Errorf("run with a exist id:%+v", m.Config.processID)
	}

	log.Infow("process id", "id", m.Config.newProcessID())
	m.Output = m.Config.ProcessPath()
	stat, e := os.Stat(m.Config.ProcessPath())
	if e != nil {
		if os.IsNotExist(e) {
			_ = os.MkdirAll(m.Config.ProcessPath(), 0755)
		} else {
			return Err(e, "stat")
		}
	}
	if e == nil && !stat.IsDir() {
		return errors.New("target is not dir")
	}
	e = m.Config.Action()
	if e != nil {
		return Err(e, "action do")
	}
	args := outputArgs(m.Config, m.Input)

	outlog := make(chan string)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		e = ff.cmd.RunContext(ctx, args, outlog)
	}()
	for i2 := range outlog {
		if !m.Debug {
			continue
		}
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
