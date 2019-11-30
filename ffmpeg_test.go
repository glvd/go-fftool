package fftool_test

import (
	"testing"

	"github.com/glvd/go-fftool"
)

var _config = fftool.DefaultConfig()

func init() {

}

// TestFFMpeg_Version ...
func TestFFMpeg_Version(t *testing.T) {
	tests := []struct {
		name    string
		fields  fftool.FFMpeg
		wantErr bool
	}{
		{
			name:    "version",
			fields:  fftool.FFMpeg{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ff := fftool.NewFFMpeg(*_config)

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
