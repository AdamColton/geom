package ffmpeg

import (
	"fmt"
	"strconv"

	"github.com/adamcolton/geom/d2/grid"
)

// Settings for an ffmpeg process
type Settings struct {
	Framerate          byte
	Name               string
	ConstantRateFactor byte
	Size               grid.Pt
}

// Common aspect rations
const (
	Widescreen = 9.0 / 16.0
	Square     = 1.0
)

// Default Settings values
const (
	Framerate          = 24
	ConstantRateFactor = 25
	Name               = "out"
)

// NewWidescreen creates Settings with a Widescreen aspect ratio.
func NewWidescreen(name string, width int) *Settings {
	return NewByAspect(name, width, Widescreen)
}

// NewSquare creates Settings with a square aspect ratio.
func NewSquare(name string, width int) *Settings {
	return NewByAspect(name, width, Square)
}

// NewByAspect creates Settings with the provided aspect ratio.
func NewByAspect(name string, width int, aspect float64) *Settings {
	return (&Settings{
		Name: name,
	}).ByAspect(width, aspect)
}

// New creates Settings with the defined width and height/
func New(name string, w, h int) *Settings {
	return (&Settings{
		Name: name,
	}).Set(w, h)
}

// ByAspect update the width and height of the settings using the provided
// aspect ratio.
func (s *Settings) ByAspect(width int, aspect float64) *Settings {
	return s.Set(width, int(float64(width)*aspect))
}

// Set the width and height, which must both be multiples of 2.
func (s *Settings) Set(w, h int) *Settings {
	if w%2 == 1 {
		w++
	}
	if h%2 == 1 {
		h++
	}
	s.Size.X = w
	s.Size.Y = h
	return s
}

// Args to pass into the ffmpeg process.
func (s *Settings) Args() []string {
	framerate := Framerate
	if s.Framerate != 0 {
		framerate = int(s.Framerate)
	}
	crf := ConstantRateFactor
	if s.ConstantRateFactor != 0 {
		crf = int(s.ConstantRateFactor)
	}
	name := Name
	if s.Name != "" {
		name = s.Name
	}
	name += ".mp4"

	frStr := strconv.Itoa(framerate)
	scale := fmt.Sprintf("scale=%dx%d", s.Size.X, s.Size.Y)
	crfStr := strconv.Itoa(crf)
	return []string{"-y", "-f", "image2pipe", "-vcodec", "bmp", "-r", frStr, "-i", "-", "-vf", scale, "-vcodec", "libx264", "-crf", crfStr, "-pix_fmt", "yuv420p", name}
}
