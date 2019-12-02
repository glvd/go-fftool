package fftool

import (
	"context"
	"testing"
)

var testVideo = `D:\video\周杰伦 唱歌贼难听.2019.1080P.h264.aac.Japanese.None.mp4`
var testStreamFormat *StreamFormat

func init() {
	var err error
	p := NewFFProbe()
	testStreamFormat, err = p.StreamFormat(testVideo)
	if err != nil {
		panic(err)
	}
}

// TestFFMpeg_Version ...
func TestFFMpeg_Version(t *testing.T) {
	tests := []struct {
		name    string
		fields  FFMpeg
		wantErr bool
	}{
		{
			name:    "version",
			fields:  FFMpeg{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ff := NewFFMpeg(*DefaultConfig())

			got, err := ff.Version()
			if (err != nil) != tt.wantErr {
				t.Errorf("Version() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != "" {
				t.Logf("Version() got = %v", got)
			}
		})
	}
}

// TestFFMpeg_Run ...
func TestFFMpeg_Run(t *testing.T) {
	type fields struct {
		config Config
		cmd    *Command
		Name   string
	}
	type args struct {
		ctx    context.Context
		input  string
		output string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "run",
			fields: fields{},
			args: args{
				ctx:    context.Background(),
				input:  testVideo,
				output: "d:\\temp\\",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := DefaultConfig()
			cfg.UseGPU = true
			err := cfg.OptimizeWithFormat(testStreamFormat)
			if err != nil {
				t.Errorf("OptimizeWithFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			ff := NewFFMpeg(*cfg)

			if err := ff.Run(tt.args.ctx, tt.args.input, tt.args.output); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
