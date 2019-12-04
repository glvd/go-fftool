package fftool

import (
	"context"
	"fmt"
	"github.com/goextension/log"
	"github.com/google/uuid"
	"path/filepath"
	"strings"
	"sync"
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

// OptimizeWithFormat ...
func (ff *FFMpeg) OptimizeWithFormat(sfmt *StreamFormat) (newFF *FFMpeg, e error) {
	cfg := ff.config.Clone()
	e = OptimizeWithFormat(&cfg, sfmt)
	if e != nil {
		return nil, e
	}
	newFF = NewFFMpeg(&cfg)
	newFF.Name = ff.Name
	return
}

// Run ...
func (ff FFMpeg) Run(ctx context.Context, input string) (e error) {
	ff.init()
	args := outputArgs(ff.config, input)

	outlog := make(chan string)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		e = ff.cmd.RunContext(ctx, args, outlog)
	}()
	for i2 := range outlog {
		log.Infow("run", "info", strings.TrimSpace(i2))
	}
	wg.Wait()
	return e
}

// NewFFMpeg ...
func NewFFMpeg(config *Config) *FFMpeg {
	ff := &FFMpeg{
		config: config,
		Name:   "ffmpeg",
	}

	return ff
}

func outputArgs(c *Config, input string) string {
	var exts []interface{}

	if c.ProcessCore != ProcessCPU && c.videoFormat != "copy" {
		c.videoFormat = "h264_nvenc"
	}

	if c.Scale != -1 {
		log.Infow("scale", "scale", c.Scale, "value", scaleVale(c.Scale))
		if c.ProcessCore != ProcessCUVID {
			exts = append(exts, fmt.Sprintf(scaleOutputTemplate, scaleVale(c.Scale)))
		} else {
			exts = append(exts, fmt.Sprintf(cuvidScaleOutputTemplate, scaleVale(c.Scale)))
		}
	}
	if c.BitRate != 0 {
		exts = append(exts, fmt.Sprintf(bitRateOutputTemplate, c.BitRate/1024))
	}
	if c.FrameRate != 0 {
		exts = append(exts, fmt.Sprintf(frameRateOutputTemplate, c.FrameRate))
	}

	output := ""
	if c.NeedSlice {
		if filepath.Ext(c.OutputName) != "" {
			//fix slice output name
			log.Infow("runme")
			c.OutputName = uuid.New().String()

		}
		output = filepath.Join(c.AbsOutput(), c.OutputName)
		output = fmt.Sprintf(sliceOutputTemplate, c.HLSTime, filepath.Join(output, c.SegmentFileName), filepath.Join(output, c.M3U8Name))
	} else {
		if filepath.Ext(c.OutputName) == "" {
			//fix media output name
			c.OutputName += ".mp4"
		}
		output = filepath.Join(c.AbsOutput(), c.OutputName)
	}

	return outputTemplate(c.ProcessCore, input, c.videoFormat, c.audioFormat, output, exts...)
}

func outputTemplate(p ProcessCore, input, cv, ca, output string, exts ...interface{}) string {
	var outExt []string
	exts = append(exts, output)
	for range exts {
		outExt = append(outExt, "%s")
	}
	var def string
	if p == ProcessCPU {
		def = fmt.Sprintf(defaultTemplate, "", input, cv, ca, strings.Join(outExt, " "))
	} else if p == ProcessCUDA {
		def = fmt.Sprintf(defaultTemplate, cudaOutputTemplate, input, cv, ca, strings.Join(outExt, " "))
	} else if p == ProcessCUVID {
		def = fmt.Sprintf(defaultTemplate, cuvidOutputTemplate, input, cv, ca, strings.Join(outExt, " "))
	}
	log.Infow("format", "def", def)
	return fmt.Sprintf(def, exts...)
}
