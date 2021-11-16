package ffmpeg

import (
	"image"
)

// Framer provides an interface for rendering frames to ffmpeg.
type Framer interface {
	// Frame renders the frame at the given index. The img passed in allows for
	// recycling. If it is nil, a new image should be generated.
	Frame(idx int, img image.Image) (image.Image, error)

	// Frames returns the number of frames to render
	Frames() int
}

// Framer will iterate through each frame of each framers. They are passed to
// the ffpeg process through a Pipeline.
func (s *Settings) Framer(fs ...Framer) (err error) {
	return s.RunPipeline(func(p *Pipeline) (err error) {
		var img image.Image
		for _, f := range fs {
			fs := f.Frames()
			for i := 0; i < fs; i++ {
				img, err = f.Frame(i, img)
				if err != nil {
					return
				}
				img = p.Add(img)
			}
		}
		p.Wait()
		return
	})
}
