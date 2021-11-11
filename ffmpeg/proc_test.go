package ffmpeg

import (
	"bytes"
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockFramer struct {
	frames int
	size   image.Rectangle
}

func (mf mockFramer) Frame(idx int, img image.Image) (image.Image, error) {
	if img == nil {
		img = image.NewGray16(mf.size)
	}
	return img, nil
}

func (mf mockFramer) Frames() int {
	return mf.frames
}

func TestBasicPipeline(t *testing.T) {
	newCommander = newMockCommand
	defer func() {
		newCommander = newCommand
	}()

	o, e := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	s := NewSquare("BasicPipeline", 64).
		SetOut(o, e)

	mf := mockFramer{
		frames: 10,
		size:   image.Rect(0, 0, s.Size.X, s.Size.Y),
	}
	err := s.Framer(mf)
	assert.NoError(t, err)
	assert.Equal(t, "ffmpeg", lastMockCommand.name)
	expected := []string{
		"-y", "-f", "image2pipe", "-vcodec", "bmp", "-r", "24", "-i", "-",
		"-vf", "scale=64x64", "-vcodec", "libx264", "-crf", "25", "-pix_fmt",
		"yuv420p", "BasicPipeline.mp4",
	}
	assert.Equal(t, expected, lastMockCommand.args)
	header := 54
	bytesPerPx := 3
	assert.Equal(t, mf.frames*(bytesPerPx*s.Size.X*s.Size.Y+header), lastMockCommand.stdin.Len())
	assert.Equal(t, o, lastMockCommand.stdout)
	assert.Equal(t, e, lastMockCommand.stderr)
}
