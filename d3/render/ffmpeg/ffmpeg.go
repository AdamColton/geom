package ffmpeg

import (
	"errors"
	"image"
	"io"
	"os"
	"os/exec"
	"strconv"

	"golang.org/x/image/bmp"
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

	p.cmd = exec.Command("ffmpeg", "-y", "-f", "image2pipe", "-vcodec", "bmp", "-r", strconv.Itoa(framerate), "-i", "-", "-vf", "scale="+strconv.Itoa(p.Width)+"x"+strconv.Itoa(p.Height), "-vcodec", "libx264", "-crf", strconv.Itoa(crf), "-pix_fmt", "yuv420p", name+".mp4")
	var err error
	p.in, err = p.cmd.StdinPipe()
	if err != nil {
		return err
	}
	p.cmd.Stdout = os.Stdout
	p.cmd.Stderr = os.Stdout
	return p.cmd.Start()
}

func (p *Proc) AddFrame(img image.Image) error {
	return bmp.Encode(p.in, img)
}

func (p *Proc) Close() error {
	p.in.Close()
	return p.cmd.Wait()
}
