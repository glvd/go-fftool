package fftool

import (
	"bufio"
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	exio "github.com/goextension/io"
	"github.com/goextension/log"
)

// DefaultCommandPath ...
var DefaultCommandPath = "bin"

// Command ...
type Command struct {
	path string
	Name string
	Args []string
	//OutPath string
	//Opts    map[string][]string
}

// Path ...
func (c *Command) Path() string {
	if filepath.IsAbs(c.path) {
		return c.path
	}
	return filepath.Join(getCurrentDir(), c.path)
}

// SetPath ...
func (c *Command) SetPath(path string) {
	c.path = path
}

// New ...
func New(name string) *Command {
	return &Command{
		path: DefaultCommandPath,
		Name: name,
		//Opts: make(map[string][]string),
	}
}

// NewFFMpeg ...
func NewFFMpeg() *Command {
	return New("ffmpeg")
}

// NewFFProbe ...
func NewFFProbe() *Command {
	return New("ffprobe")
}

// SetArgs ...
func (c *Command) SetArgs(s string) {
	c.Args = strings.Split(s, ",")
}

// AddArgs ...
func (c *Command) AddArgs(s string) {
	c.Args = append(c.Args, s)
}

func getCurrentDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Errorw("current dir", "error", err)
		return ""
	}
	return dir
}

// Env ...
func (c *Command) Env() []string {
	if err := os.Setenv("PATH", strings.Join([]string{os.Getenv("PATH"), c.Path()}, ":")); err != nil {
		panic(err)
	}
	return os.Environ()
}

// Run ...
func (c *Command) Run() (string, error) {
	cmd := exec.Command(c.Name, c.Args...)
	cmd.Env = c.Env()
	//显示运行的命令
	log.Infow("run", "args", cmd.Args)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return string(stdout), errWrap(err, "run")
	}
	return string(stdout), nil
}

// RunContext ...
func (c *Command) RunContext(ctx context.Context, info chan<- string) (e error) {
	defer func() {
		if info != nil {
			close(info)
		}
	}()
	cmd := exec.CommandContext(ctx, c.Name, c.Args...)
	cmd.Env = c.Env()
	//显示运行的命令
	log.Infow("run context", "args", cmd.Args)
	stdout, e := cmd.StdoutPipe()
	if e != nil {
		return errWrap(e, "StdoutPipe")
	}

	stderr, e := cmd.StderrPipe()
	if e != nil {
		return errWrap(e, "StderrPipe")
	}

	e = cmd.Start()
	if e != nil {
		return errWrap(e, "start")
	}

	reader := bufio.NewReader(exio.MultiReader(stderr, stdout))
	for {
		select {
		case <-ctx.Done():
			return errWrap(ctx.Err(), "done")
		default:
			lines, _, e := reader.ReadLine()
			if e != nil || io.EOF == e {
				goto END
			}
			if strings.TrimSpace(string(lines)) != "" {
				if info != nil {
					info <- string(lines)
				}
			}
		}
	}
END:
	e = cmd.Wait()
	if e != nil {
		return e
	}
	return nil
}
