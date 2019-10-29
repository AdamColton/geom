package ffmpeg

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"strconv"

	"github.com/fogleman/gg"
)

type Proc struct {
	Framerate          byte
	Name               string
	ConstantRateFactor byte
	Width, Height      int
	cmd                *exec.Cmd
	in                 io.WriteCloser
}

func (p *Proc) Start() error {
	if p.cmd != nil {
		return errors.New("Already running")
	}
	framerate := 24
	if p.Framerate != 0 {
		framerate = int(p.Framerate)
	}
	crf := 25
	if p.ConstantRateFactor != 0 {
		crf = int(p.ConstantRateFactor)
	}
	name := "out"
	if p.Name != "" {
		name = p.Name
	}
	p.cmd = exec.Command("ffmpeg", "-y", "-f", "image2pipe", "-vcodec", "png", "-r", strconv.Itoa(framerate), "-i", "-", "-vf", "scale="+strconv.Itoa(p.Width)+"x"+strconv.Itoa(p.Height), "-vcodec", "libx264", "-crf", strconv.Itoa(crf), "-pix_fmt", "yuv420p", name+".mp4")
	var err error
	p.in, err = p.cmd.StdinPipe()
	if err != nil {
		return err
	}
	p.cmd.Stdout = os.Stdout
	p.cmd.Stderr = os.Stdout
	return p.cmd.Start()
}

func (p *Proc) AddFrame(ctx *gg.Context) error {
	return ctx.EncodePNG(p.in)
}

func (p *Proc) Close() error {
	p.in.Close()
	return p.cmd.Wait()
}
