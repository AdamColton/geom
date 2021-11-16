package ffmpeg

import (
	"bytes"
	"io"
)

type bufCloser struct {
	*bytes.Buffer
}

func (bc bufCloser) Close() error {
	return nil
}

type mockcommand struct {
	name           string
	args           []string
	stdin          bufCloser
	stdout, stderr io.Writer
}

func (c *mockcommand) StdinPipe() (io.WriteCloser, error) {
	if c.stdin.Buffer == nil {
		c.stdin.Buffer = bytes.NewBuffer(nil)
	}
	return c.stdin, nil
}

func (c *mockcommand) SetStdout(w io.Writer) { c.stdout = w }

func (c *mockcommand) SetStderr(w io.Writer) { c.stderr = w }

func (c *mockcommand) Start() error {
	return nil
}

func (c *mockcommand) Wait() error {
	return nil
}

var lastMockCommand *mockcommand

func newMockCommand(name string, args ...string) commander {
	lastMockCommand = &mockcommand{
		name: name,
		args: args,
	}
	return lastMockCommand
}
