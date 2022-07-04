package ffmpeg

import (
	"io"
	"os/exec"
)

type commander interface {
	StdinPipe() (io.WriteCloser, error)
	SetStdout(io.Writer)
	SetStderr(io.Writer)
	Start() error
	Wait() error
}

type command struct {
	cmd *exec.Cmd
}

func newCommand(name string, args ...string) commander {
	return &command{
		cmd: exec.Command(name, args...),
	}
}

func (c *command) StdinPipe() (io.WriteCloser, error) { return c.cmd.StdinPipe() }

func (c *command) SetStdout(w io.Writer) { c.cmd.Stdout = w }

func (c *command) SetStderr(w io.Writer) { c.cmd.Stderr = w }

func (c *command) Start() error { return c.cmd.Start() }

func (c *command) Wait() error { return c.cmd.Wait() }

var newCommander = newCommand
