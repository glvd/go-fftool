package fftool

import (
	"context"
	"path/filepath"
	"sync"
	"testing"

	"github.com/goextension/log"
)

func init() {
	DefaultCommandPath = filepath.Join(`D:\workspace\golang\project\go-fftool`, DefaultCommandPath)
}

// TestCommand_RunContext ...
func TestCommand_RunContext(t *testing.T) {
	type fields struct {
		path string
		Name string
	}
	type args struct {
		ctx  context.Context
		info chan string
		args string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "runcontext1",
			fields: fields{
				path: DefaultCommandPath,
				Name: "ffmpeg",
			},
			args: args{
				ctx:  context.Background(),
				args: "-version",
			},
			wantErr: false,
		}, {
			name: "runcontext2",
			fields: fields{
				path: DefaultCommandPath,
				Name: "ffmpeg",
			},
			args: args{
				ctx:  context.Background(),
				info: make(chan string),
				args: "-version",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Command{
				path: tt.fields.path,
				Name: tt.fields.Name,
			}
			wg := &sync.WaitGroup{}
			wg.Add(1)
			go func() {
				if err := c.RunContext(tt.args.ctx, tt.args.args, tt.args.info); (err != nil) != tt.wantErr {
					t.Errorf("RunContext() error = %v, wantErr %v", err, tt.wantErr)
				}
				wg.Done()
			}()
			if tt.args.info != nil {
				for v := range tt.args.info {
					log.Info(v)
				}
			}
			wg.Wait()
		})
	}
}

// TestCommand_Run ...
func TestCommand_Run(t *testing.T) {
	type fields struct {
		path string
		env  []string
		Name string
	}
	type args struct {
		args string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "testrun",
			fields: fields{
				path: DefaultCommandPath,
				env:  nil,
				Name: "ffmpeg",
			},
			args: args{
				args: "-version",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Command{
				path: tt.fields.path,
				env:  tt.fields.env,
				Name: tt.fields.Name,
			}
			got, err := c.Run(tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != "" {
				t.Logf("Run() got = %v", got)
			}
		})
	}
}
