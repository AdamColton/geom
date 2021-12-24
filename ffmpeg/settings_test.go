package ffmpeg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSettings(t *testing.T) {
	tt := map[string]struct {
		Settings *Settings
		expected []string
	}{
		"New": {
			Settings: New("New").Set(101, 201),
			expected: []string{"-y", "-f", "image2pipe", "-vcodec", "bmp", "-r", "24", "-i", "-", "-vf", "scale=100x200", "-vcodec", "libx264", "-crf", "25", "-pix_fmt", "yuv420p", "New.mp4"},
		},
		"NewByAspect": {
			Settings: NewByAspect("NewByAspect", 200, 1.75),
			expected: []string{"-y", "-f", "image2pipe", "-vcodec", "bmp", "-r", "24", "-i", "-", "-vf", "scale=200x350", "-vcodec", "libx264", "-crf", "25", "-pix_fmt", "yuv420p", "NewByAspect.mp4"},
		},
		"NewSquare": {
			Settings: NewSquare("NewSquare", 250),
			expected: []string{"-y", "-f", "image2pipe", "-vcodec", "bmp", "-r", "24", "-i", "-", "-vf", "scale=250x250", "-vcodec", "libx264", "-crf", "25", "-pix_fmt", "yuv420p", "NewSquare.mp4"},
		},
		"NewWidescreen": {
			Settings: NewWidescreen("NewWidescreen", 150),
			expected: []string{"-y", "-f", "image2pipe", "-vcodec", "bmp", "-r", "24", "-i", "-", "-vf", "scale=150x84", "-vcodec", "libx264", "-crf", "25", "-pix_fmt", "yuv420p", "NewWidescreen.mp4"},
		},
		"Framerate": {
			Settings: (&Settings{
				Framerate: 26,
				Name:      "Framerate",
			}).Set(125, 456),
			expected: []string{"-y", "-f", "image2pipe", "-vcodec", "bmp", "-r", "26", "-i", "-", "-vf", "scale=124x456", "-vcodec", "libx264", "-crf", "25", "-pix_fmt", "yuv420p", "Framerate.mp4"},
		},
		"ConstantRateFactor": {
			Settings: (&Settings{
				ConstantRateFactor: 30,
				Name:               "ConstantRateFactor",
			}).Set(125, 456),
			expected: []string{"-y", "-f", "image2pipe", "-vcodec", "bmp", "-r", "24", "-i", "-", "-vf", "scale=124x456", "-vcodec", "libx264", "-crf", "30", "-pix_fmt", "yuv420p", "ConstantRateFactor.mp4"},
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.Settings.Args())
		})
	}
}
