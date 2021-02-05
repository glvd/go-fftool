package tool

import (
	"context"
	"fmt"
	extlog "github.com/goextension/log"
	"github.com/goextension/log/zap"
	"sync"
	"testing"
)

var testVideo = `D:\workspace\golang\project\go-media-tool\test\media.mp4`
var testStreamFormat *StreamFormat

func init() {
	DefaultMpegName = "D:\\workspace\\golang\\project\\ipvc\\bin\\ffmpeg.exe"
	DefaultProbeName = "D:\\workspace\\golang\\project\\ipvc\\bin\\ffprobe.exe"
	zap.InitZapSugar()
	log = extlog.Log()
	var err error
	p := NewFFProbe()
	testStreamFormat, err = p.StreamFormat(testVideo)
	if err != nil {

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
			ff := NewFFMpeg()
			got, err := ff.cmd.Run("-version")
			//got, err := ff.Version()
			if (err != nil) != tt.wantErr {
				t.Errorf("Version() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != "" {
				fmt.Println(got)
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
			name:   "run1",
			fields: fields{},
			args: args{
				ctx:   context.Background(),
				input: testVideo,
			},
			wantErr: false,
		},
	}
	wg := &sync.WaitGroup{}
	for _, tt := range tests {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cfg := DefaultConfig()
			cfg.Slice = true
			cfg.ProcessCore = ProcessHevcAMF
			//c := GenerateCrypto(NewOpenSSL(), true)
			//
			//cfg.SetCrypt(*c)
			cfg.LogOutput = true
			cfg.Slice = true
			ff := NewFFMpeg(cfg.ConfigOptions())
			ff.HandleMessage(func(s string) {
				fmt.Println(s)
			})
			e := OptimizeWithFormat(cfg, testStreamFormat)
			if e != nil {
				t.Errorf("OptimizeWithFormat() error = %v, wantErr %v", e, tt.wantErr)
				return
			}
			if err := ff.Run(tt.args.ctx, tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		}()
	}
	wg.Wait()
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
				VideoFormat:     tt.fields.VideoFormat,
				AudioFormat:     tt.fields.AudioFormat,
				M3U8Name:        tt.fields.M3U8Name,
				SegmentFileName: tt.fields.SegmentFileName,
				HLSTime:         tt.fields.HLSTime,
			}

			ff := NewFFMpeg()
			e := OptimizeWithFormat(c, tt.args.sfmt)
			if (e != nil) != tt.wantErr {
				t.Errorf("OptimizeWithFormat() error = %v, wantErr %v", e, tt.wantErr)
				return
			}
			if _, err := ff.Version(); err != nil {
				t.Errorf("Version() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("config:%+v\n", c)
		})
	}
}
