package tool

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	exio "github.com/goextension/io"
)

// DefaultCommandPath default point to current dir
var DefaultCommandPath = ""

// CommandRunner ...
type CommandRunner interface {
	Message(func(message string))
	Run(s string) (string, error)
	RunContext(ctx context.Context, s string) (e error)
}

// Command ...
type Command struct {
	runArgs string
	bin     string
	message func(s string)
}

var _ CommandRunner = &Command{}
var _env []string

// environ ...
func (c *Command) environ() []string {
	if _env == nil {
		_, e := exec.LookPath(filepath.Join(c.BinPath(), c.bin))
		if e == nil {
			if err := os.Setenv("PATH", strings.Join([]string{os.Getenv("PATH"), c.BinPath()}, string(os.PathListSeparator))); err != nil {
				panic(err)
			}
		}
		_env = os.Environ()
	}
	return _env
}

// BinPath ...
func (c *Command) BinPath() string {
	if filepath.IsAbs(c.bin) {
		return c.bin
	}
	if DefaultCommandPath != "" {
		if filepath.IsAbs(DefaultCommandPath) {
			filepath.Join(DefaultCommandPath, c.bin)
		}
	}
	return filepath.Join(getCurrentDir(), c.bin)
}

// NewCommand ...
func NewCommand(name string) *Command {
	if !filepath.IsAbs(name) {
		name = binaryExt(name)
	}
	c := &Command{
		bin: name,
	}
	c.environ()
	return c
}

func getCurrentDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.cmdArgs[0])去除最后一个元素的路径
	if err != nil {
		log.Errorw("current dir", "error", err)
		return ""
	}
	return dir
}

// Run ...
func (c *Command) Run(args string) (string, error) {
	cmd := exec.Command(c.BinPath(), cmdArgs(args)...)
	cmd.Env = c.environ()
	//显示运行的命令
	log.Infow("run", "outputArgs", cmd.Args)
	c.runArgs = strings.Join(cmd.Args, " ")
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return string(stdout), Err(err, "run")
	}
	return string(stdout), nil
}

// Message ...
func (c *Command) Message(f func(message string)) {
	c.message = f
}

// RunContext ...
func (c *Command) RunContext(ctx context.Context, args string) (e error) {
	cctx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()
	cmd := exec.CommandContext(cctx, c.BinPath(), cmdArgs(args)...)
	cmd.Env = c.environ()
	//显示运行的命令
	log.Infow("run context", "outputArgs", strings.Join(cmd.Args, " "))
	c.runArgs = strings.Join(cmd.Args, " ")
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

	mreader := exio.MultiReader(stderr, stdout)
	defer exio.Close(mreader)
	go func(ctx context.Context, mreader io.Reader) {
		reader := bufio.NewReader(mreader)
		var lines []byte
		for {
			select {
			case <-ctx.Done():
				return
			default:
				lines, _, e = reader.ReadLine()
				if e != nil {
					if e != io.EOF {
						log.Error(Err(e, "readline"))
					}
					return
				}
				if l := string(bytes.TrimSpace(lines)); l != "" {
					if c.message != nil {
						c.message(l)
					}
				}
			}
		}
	}(cctx, mreader)
	if err := cmd.Wait(); err != nil {
		if err != io.EOF {
			return Err(err, "wait")
		}
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
