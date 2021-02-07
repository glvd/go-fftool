package tool

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"os"
	"strings"
)

// FFMpeg ...
type FFMpeg struct {
	cmd    *Command
	config *Config
	name   string
	msg    func(s string)
}

// MpegOption ...
type MpegOption struct {
	Debug  bool
	Config *Config
}

// ConfigOptions ...
type ConfigOptions func(cfg *Config)

// Name ...
func (ff FFMpeg) Name() string {
	return ff.name
}

// Version ...
func (ff *FFMpeg) Version() (string, error) {
	return ff.cmd.Run("-version")
}

// HandleMessage ...
func (ff *FFMpeg) HandleMessage(msg func(s string)) {
	ff.msg = msg
}

// Run ...
func (ff FFMpeg) Run(ctx context.Context, input string) (e error) {
	pid := uuid.New().String()

	ff.config.processID = pid

	log.Infow("process id", "id", ff.config.ProcessID())
	stat, e := os.Stat(ff.config.ProcessPath())
	if e != nil {
		if os.IsNotExist(e) {
			_ = os.MkdirAll(ff.config.ProcessPath(), 0755)
		} else {
			return Err(e, "stat")
		}
	}
	if e == nil && !stat.IsDir() {
		return errors.New("target is not dir")
	}
	e = ff.config.Action()
	if e != nil {
		return Err(e, "action do")
	}
	args := outputArgs(ff.config, input)
	ff.cmd.Message(ff.messageCallback)
	e = ff.cmd.RunContext(ctx, args)

	return e
}

func (ff *FFMpeg) messageCallback(message string) {
	if ff.msg != nil {
		ff.msg(message)
	}
	if ff.config.LogOutput {
		str := strings.Split(message, "\r")
		for i := range str {
			log.Infow(strings.TrimSpace(str[i]))
		}
	}
}

// RunCommandString ...
func (ff *FFMpeg) RunCommandString() string {
	return ff.cmd.runArgs
}

// NewFFMpeg ...
func NewFFMpeg(opts ...ConfigOptions) *FFMpeg {
	ff := &FFMpeg{
		name:   DefaultMpegName,
		config: DefaultConfig(),
	}

	for _, opt := range opts {
		opt(ff.config)
	}

	ff.cmd = NewCommand(ff.name)
	return ff
}
