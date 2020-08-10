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
	path    string
	env     []string
	bin     string
	message func(s string)
}

var _ CommandRunner = &Command{}

// Path ...
func (c *Command) BinPath() string {
	if filepath.IsAbs(c.path) {
		return filepath.Join(c.path, c.bin)
	}
	return filepath.Join(getCurrentDir(), c.path, c.bin)
}

// SetPath ...
func (c *Command) SetPath(path string) {
	c.path = path
}

// NewCommand ...
func NewCommand(name string) *Command {
	return &Command{
		path: DefaultCommandPath,
		bin:  binaryExt(name),
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

// Run ...
func (c *Command) Run(args string) (string, error) {
	cmd := exec.Command(c.BinPath(), cmdArgs(args)...)
	//显示运行的命令
	log.Infow("run", "outputArgs", cmd.Args)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return string(stdout), Err(err, "run")
	}
	return string(stdout), nil
}

func (c *Command) Message(f func(message string)) {
	c.message = f
}

// RunContext ...
func (c *Command) RunContext(ctx context.Context, args string) (e error) {
	cmd := exec.CommandContext(ctx, c.BinPath(), cmdArgs(args)...)
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
			if lines = bytes.TrimSpace(lines); lines != nil {
				if c.message != nil {
					c.message(string(lines))
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
