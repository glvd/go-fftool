package fftool

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	exio "github.com/goextension/io"
	"github.com/goextension/log"
)

// DefaultCommandPath default point to current dir
var DefaultCommandPath = ""

// CommandRunner ...
type CommandRunner interface {
	Run() (string, error)
	RunContext(ctx context.Context, info chan<- string) (e error)
}

// Command ...
type Command struct {
	path string
	env  []string
	Name string
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

// NewCommand ...
func NewCommand(name string) *Command {
	return &Command{
		path: DefaultCommandPath,
		Name: name,
	}
}

func getCurrentDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.cmdArgs[0])去除最后一个元素的路径
	if err != nil {
		log.Errorw("current dir", "error", err)
		return ""
	}
	return dir
}

// environ ...
func (c *Command) init() []string {
	if c.env == nil {
		_, e := exec.LookPath(c.Path())
		if e != nil {
			log.Warnw("init", "path", c.Path())
			if err := os.Setenv("PATH", strings.Join([]string{os.Getenv("PATH"), c.Path()}, string(os.PathListSeparator))); err != nil {
				panic(err)
			}
		}
		c.env = os.Environ()
	}
	return c.env
}

// Run ...
func (c *Command) Run(args string) (string, error) {
	c.init()
	cmd := exec.Command(c.Name, cmdArgs(args)...)
	//显示运行的命令
	log.Infow("run", "outputArgs", cmd.Args)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return string(stdout), Err(err, "run")
	}
	return string(stdout), nil
}

// RunContext ...
func (c *Command) RunContext(ctx context.Context, args string, info chan<- string) (e error) {
	defer func() {
		if info != nil {
			close(info)
		}
	}()
	c.init()
	cmd := exec.CommandContext(ctx, c.Name, cmdArgs(args)...)
	//显示运行的命令
	log.Infow("run context", "outputArgs", strings.Join(cmd.Args, " "))
	stdout, e := cmd.StdoutPipe()
	if e != nil {
		return Err(e, "StdoutPipe")
	}

	stderr, e := cmd.StderrPipe()
	if e != nil {
		return Err(e, "StderrPipe")
	}

	e = cmd.Start()
	if e != nil {
		return Err(e, "start")
	}

	reader := bufio.NewReader(exio.MultiReader(stderr, stdout))
	var lines []byte
	for {
		select {
		case <-ctx.Done():
			return Err(ctx.Err(), "done")
		default:
			lines, _, e = reader.ReadLine()
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

// cmdArgs ...
func cmdArgs(args string) []string {
	return strings.Split(args, ",")
}

// Args ...
func Args(s ...string) string {
	return strings.Join(s, ",")
}

// formatArgs ...
func formatArgs(source string, args ...interface{}) string {
	return fmt.Sprintf(source, args...)
}
