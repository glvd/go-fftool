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
				file: "d:\\video\\周杰伦 唱歌贼难听.2019.1080P.h264.aac.Japanese.None.mp4",
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
