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
	env  []string
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

// environ ...
func (c *Command) init() []string {
	if c.env == nil {
		if err := os.Setenv("PATH", strings.Join([]string{os.Getenv("PATH"), c.Path()}, string(os.PathListSeparator))); err != nil {
			panic(err)
		}
		c.env = os.Environ()
	}
	return c.env
}

// Run ...
func (c *Command) Run() (string, error) {
	c.init()
	cmd := exec.Command(c.Name, c.Args...)
	//显示运行的命令
	log.Infow("run", "args", cmd.Args, "environ", c.env)
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
	c.init()
	cmd := exec.CommandContext(ctx, c.Name, c.Args...)
	//显示运行的命令
	log.Infow("run context", "args", cmd.Args, "environ", c.env)
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
