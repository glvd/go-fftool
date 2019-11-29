package fftool

import (
	"testing"
)

// TestFFProbe_StreamFormat ...
func TestFFProbe_StreamFormat(t *testing.T) {
	type fields struct {
		cmd  *Command
		Name string
	}
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				file: `d:\video\周杰伦 唱歌贼难听.2019.1080P.h264.aac.Japanese.None.mp4`,
			},
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				file: `d:\video\女大学生的沙龙室.Room.Salon.College.Girls.2018.HD720P.X264.AAC.Korean.CHS.mp4`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ff := NewFFProbe()
			got, err := ff.StreamFormat(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("StreamFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				t.Logf("StreamFormat() got = %v", got)
			}
		})
	}
}

// Test_getResolutionIndex ...
func Test_getResolutionIndex(t *testing.T) {
	type args struct {
		n   int64
		sta int
		end int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				n:   180,
				sta: 0,
				end: -1,
			},
			want: 240,
		},
		{
			name: "2",
			args: args{
				n:   240,
				sta: 0,
				end: -1,
			},
			want: 240,
		},
		{
			name: "3",
			args: args{
				n:   1080,
				sta: 0,
				end: -1,
			},
			want: 1080,
		},
		{
			name: "4",
			args: args{
				n:   4800,
				sta: 0,
				end: -1,
			},
			want: 4800,
		},
		{
			name: "5",
			args: args{
				n:   4900,
				sta: 0,
				end: -1,
			},
			want: 4800,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getResolution(tt.args.n, tt.args.sta, tt.args.end); got != tt.want {
				t.Errorf("getResolution() = %v, want %v", got, tt.want)
			}
		})
	}
}
