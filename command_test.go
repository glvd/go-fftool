package tool

import (
	"context"
	"os"
	"reflect"
	"sync"
	"testing"
)

func init() {

}

// TestCommand_RunContext ...
func TestCommand_RunContext(t *testing.T) {
	type fields struct {
		path string
		bin  string
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
				bin:  "ffmpeg",
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
				bin:  "ffmpeg",
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
				bin: tt.fields.bin,
			}
			wg := &sync.WaitGroup{}
			wg.Add(1)

			c.message = func(s string) {
				log.Info(s)
			}
			go func() {
				if err := c.RunContext(tt.args.ctx, tt.args.args); (err != nil) != tt.wantErr {
					t.Errorf("RunContext() error = %v, wantErr %v", err, tt.wantErr)
				}
				wg.Done()
			}()

			wg.Wait()
		})
	}
}

// TestCommand_Run ...
func TestCommand_Run(t *testing.T) {
	type fields struct {
		bin string
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
				bin: "ffmpeg",
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
				bin: tt.fields.bin,
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

// TestEnviron ...
func TestEnviron(t *testing.T) {
	type args struct {
		env  []string
		path string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "",
			args: args{
				env:  []string{os.Getenv("PATH")},
				path: "bin",
			},
			want: nil,
		},
	}
	c := NewCommand("ffmpeg")
	c.environ()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.environ(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Environ() = %v \n want %v", got, os.Getenv("PATH"))
			}
		})
	}
}
