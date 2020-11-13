package ffmpeg

import (
	"image"
	"io"
	"os"
	"os/exec"

	"golang.org/x/image/bmp"
)

// Proc is a running instance of ffmpeg.
type Proc struct {
	Settings
	cmd *exec.Cmd
	in  io.WriteCloser
}

// NewProc created from the Settings.
func (s *Settings) NewProc() (*Proc, error) {
	p := &Proc{
		Settings: *s,
	}

	p.cmd = exec.Command("ffmpeg", p.Args()...)
	var err error
	p.in, err = p.cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	p.cmd.Stdout = os.Stdout
	p.cmd.Stderr = os.Stdout
	err = p.cmd.Start()
	if err != nil {
		return nil, err
	}
	return p, nil
}

// AddFrame to the video created by ffmpeg
func (p *Proc) AddFrame(img image.Image) error {
	return bmp.Encode(p.in, img)
}

// Close the ffmpeg process. This must be called to finish the process.
func (p *Proc) Close() error {
	p.in.Close()
	return p.cmd.Wait()
}

// Run will start and close the ffmpeg Proc, passing the Proc into the provided
// func for use.
func (s *Settings) Run(fn func(*Proc) error) (err error) {
	p, err := s.NewProc()
	if err != nil {
		return
	}
	defer func() {
		closeErr := p.Close()
		if err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	err = fn(p)
	return
}

// RunPipeline creates a Pipeline from the Settings and passes it into the func
// provided.
func (s *Settings) RunPipeline(fn func(*Pipeline) error) (err error) {
	return s.Run(func(p *Proc) error {
		return fn(p.Pipeline())
	})
}
