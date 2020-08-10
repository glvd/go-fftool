package tool

import (
	"context"
	"errors"
	"github.com/google/uuid"
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
}

// RunOptions ...
type RunOptions func(cfg *Config) *Config

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
	pid := uuid.New().String()
	config := DefaultConfig()

	config.processID = pid
	for _, opt := range opts {
		config = opt(config)
	}
	if config.processID == "" {
		config.processID = pid
	}

	log.Infow("process id", "id", config.ProcessID())
	stat, e := os.Stat(config.ProcessPath())
	if e != nil {
		if os.IsNotExist(e) {
			_ = os.MkdirAll(config.ProcessPath(), 0755)
		} else {
			return Err(e, "stat")
		}
	}
	if e == nil && !stat.IsDir() {
		return errors.New("target is not dir")
	}
	e = config.Action()
	if e != nil {
		return Err(e, "action do")
	}
	args := outputArgs(config, input)

	var outLog chan string
	if config.LogOutput {
		outLog = make(chan string)
	}

	log.Infow("runmsg", "init", outLog == nil)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		e = ff.cmd.RunContext(ctx, args, outLog)
	}()

	if config.LogOutput {
		for i2 := range outLog {
			ss := strings.Split(i2, "\r")
			for _, i3 := range ss {
				log.Infow("runmsg", "outLog", strings.TrimSpace(i3))
			}
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
