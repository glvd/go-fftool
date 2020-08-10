package tool

import (
	"github.com/google/uuid"
	"strings"
	"testing"
)

var sf *StreamFormat

func init() {
	//DefaultCommandPath = `D:\workspace\golang\project\go-fftool\bin`
	var e error
	sf, e = NewFFProbe().StreamFormat(`D:\video\集锦-挪威混剪8.1-4k_360.mp4`)
	if e != nil {
		//ignore
	}
}

// Test_outputArgs ...
func Test_outputArgs(t *testing.T) {
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
				intput: "d:\\video\\集锦-挪威混剪8.1-4k_360.mp4",
			},
		},
		{
			name: "args2",
			fields: Config{
				Scale:           DefaultScale,
				ProcessCore:     ProcessCPU,
				BitRate:         0,
				FrameRate:       0,
				OutputPath:      DefaultOutputPath,
				OutputName:      "63ca3045-80cf-445c-a40d-d374e734350a",
				VideoFormat:     "libx264",
				AudioFormat:     "aac",
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
				BitRate:         0,
				FrameRate:       0,
				OutputPath:      DefaultOutputPath,
				OutputName:      uuid.New().String(),
				VideoFormat:     "libx264",
				AudioFormat:     "aac",
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
				output:          "",
				VideoFormat:     "libx264",
				AudioFormat:     "aac",
				crypto:          nil,
				Scale:           DefaultScale,
				ProcessCore:     DefaultProcessCore,
				BitRate:         0,
				FrameRate:       0,
				OutputPath:      DefaultOutputPath,
				OutputName:      DefaultOutputName,
				M3U8Name:        DefaultM3U8Name,
				SegmentFileName: DefaultSegmentFileName,
				HLSTime:         DefaultHLSTime,
			},
			args: args{
				intput: "d:\\video\\周杰伦 唱歌贼难听.2019.1080P.h264.aac.Japanese.None.mp4",
			},
		},
		{
			name: "args5",
			fields: Config{
				output:          "",
				VideoFormat:     "libx264",
				AudioFormat:     "aac",
				crypto:          GenerateCrypto(NewOpenSSL(), true),
				Scale:           DefaultScale,
				ProcessCore:     DefaultProcessCore,
				BitRate:         0,
				FrameRate:       0,
				OutputPath:      DefaultOutputPath,
				OutputName:      DefaultOutputName,
				M3U8Name:        DefaultM3U8Name,
				SegmentFileName: DefaultSegmentFileName,
				HLSTime:         DefaultHLSTime,
				KeyOutput:       false,
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
