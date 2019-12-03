package fftool

import (
	"context"
	"github.com/google/uuid"
	"strings"
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
			ff := NewFFMpeg(DefaultConfig())

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
		ctx   context.Context
		input string
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
				ctx:   context.Background(),
				input: testVideo,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := DefaultConfig()
			ff := NewFFMpeg(cfg)
			newff, err := ff.OptimizeWithFormat(testStreamFormat)
			if err != nil {
				t.Errorf("OptimizeWithFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err := newff.Run(tt.args.ctx, tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestFFMpeg_OptimizeWithFormat ...
func TestFFMpeg_OptimizeWithFormat(t *testing.T) {
	type fields struct {
		Scale           Scale
		BitRate         int64
		FrameRate       float64
		Output          string
		VideoFormat     string
		AudioFormat     string
		M3U8Name        string
		SegmentFileName string
		HLSTime         int
	}
	type args struct {
		sfmt *StreamFormat
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "1",
			fields: fields{
				Scale:           2,
				BitRate:         0,
				FrameRate:       0,
				Output:          "",
				VideoFormat:     "",
				AudioFormat:     "",
				M3U8Name:        "",
				SegmentFileName: "",
				HLSTime:         0,
			},
			args: args{
				sfmt: sf,
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "2",
			fields: fields{
				Scale:           1,
				BitRate:         0,
				FrameRate:       0,
				Output:          "",
				VideoFormat:     "",
				AudioFormat:     "",
				M3U8Name:        "",
				SegmentFileName: "",
				HLSTime:         0,
			},
			args: args{
				sfmt: sf,
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "3",
			fields: fields{
				Scale:           1,
				BitRate:         0,
				FrameRate:       0,
				Output:          "",
				VideoFormat:     "",
				AudioFormat:     "",
				M3U8Name:        "",
				SegmentFileName: "",
				HLSTime:         0,
			},
			args: args{
				sfmt: nil,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Scale:           tt.fields.Scale,
				BitRate:         tt.fields.BitRate,
				FrameRate:       tt.fields.FrameRate,
				OutputPath:      tt.fields.Output,
				videoFormat:     tt.fields.VideoFormat,
				audioFormat:     tt.fields.AudioFormat,
				M3U8Name:        tt.fields.M3U8Name,
				SegmentFileName: tt.fields.SegmentFileName,
				HLSTime:         tt.fields.HLSTime,
			}

			ff := NewFFMpeg(c)
			_, err := ff.OptimizeWithFormat(tt.args.sfmt)
			if (err != nil) != tt.wantErr {
				t.Errorf("OptimizeWithFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%+v", c)
		})
	}
}

// TestConfig_Args ...
func TestConfig_Args(t *testing.T) {
	type fields struct {
		Scale           Scale
		BitRate         int64
		FrameRate       float64
		Output          string
		VideoFormat     string
		AudioFormat     string
		M3U8Name        string
		SegmentFileName string
		HLSTime         int
		OutputPath      string
	}
	type args struct {
		intput string
		output string
	}
	tests := []struct {
		name   string
		fields Config
		args   args
	}{
		{
			name:   "args1",
			fields: *DefaultConfig(),
			args: args{
				intput: "d:\\video\\周杰伦 唱歌贼难听.2019.1080P.h264.aac.Japanese.None.mp4",
			},
		},
		{
			name: "args2",
			fields: Config{
				Scale:           DefaultScale,
				ProcessCore:     ProcessCPU,
				NeedSlice:       DefaultSlice,
				BitRate:         0,
				FrameRate:       0,
				OutputPath:      DefaultOutputPath,
				OutputName:      "63ca3045-80cf-445c-a40d-d374e734350a",
				videoFormat:     "libx264",
				audioFormat:     "aac",
				M3U8Name:        DefaultM3U8Name,
				SegmentFileName: DefaultSegmentFileName,
				HLSTime:         DefaultHLSTime,
			},
			args: args{
				intput: "d:\\video\\周杰伦 唱歌贼难听.2019.1080P.h264.aac.Japanese.None.mp4",
			},
		},
		{
			name: "args3",
			fields: Config{
				Scale:           DefaultScale,
				ProcessCore:     DefaultProcessCore,
				NeedSlice:       DefaultSlice,
				BitRate:         0,
				FrameRate:       0,
				OutputPath:      DefaultOutputPath,
				OutputName:      uuid.New().String(),
				videoFormat:     "libx264",
				audioFormat:     "aac",
				M3U8Name:        DefaultM3U8Name,
				SegmentFileName: DefaultSegmentFileName,
				HLSTime:         DefaultHLSTime,
			},
			args: args{
				intput: "d:\\video\\周杰伦 唱歌贼难听.2019.1080P.h264.aac.Japanese.None.mp4",
			},
		},
		{
			name: "args4",
			fields: Config{
				Scale:           DefaultScale,
				ProcessCore:     DefaultProcessCore,
				NeedSlice:       true,
				BitRate:         0,
				FrameRate:       0,
				OutputPath:      DefaultOutputPath,
				OutputName:      DefaultOutputName,
				videoFormat:     "libx264",
				audioFormat:     "aac",
				M3U8Name:        DefaultM3U8Name,
				SegmentFileName: DefaultSegmentFileName,
				HLSTime:         DefaultHLSTime,
			},
			args: args{
				intput: "d:\\video\\周杰伦 唱歌贼难听.2019.1080P.h264.aac.Japanese.None.mp4",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := outputArgs(&tt.fields, tt.args.intput); got != "" {
				t.Logf("Args() = %v", strings.ReplaceAll(got, ",", " "))
			}
		})
	}
}
