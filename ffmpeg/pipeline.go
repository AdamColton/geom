package ffmpeg

import (
	"image"
	"sync"
)

// Pipeline allows the process of encoding and passing an image to the ffmpeg
// process to happen in parallel to the process of generating frames.
type Pipeline struct {
	sync.Mutex
	wg            sync.WaitGroup
	pool          []image.Image
	frames        map[int]image.Image
	proc          *Proc
	idxIn, idxOut int
	running       bool
}

// Pipeline that operates in 2 stages. Add allows images to be added to the
// after which they are encoded and sent to the ffmeg process in a parallel
// go routine. The images are recycled on a pool.
func (p *Proc) Pipeline() *Pipeline {
	return &Pipeline{
		frames: make(map[int]image.Image),
		proc:   p,
	}
}

// Add an image to the pipeline. A recycled image is returned, but it may be
// nil.
func (p *Pipeline) Add(img image.Image) image.Image {
	p.wg.Add(1)
	p.Lock()
	p.frames[p.idxIn] = img
	p.idxIn++
	if ln := len(p.pool); ln > 0 {
		img = p.pool[ln-1]
		p.pool = p.pool[:ln-1]
	} else {
		img = nil
	}
	if !p.running {
		go p.run()
	}
	p.Unlock()
	return img
}

// Wait for the pipeline to finish processing.
func (p *Pipeline) Wait() {
	p.wg.Wait()
}

func (p *Pipeline) run() {
	p.Lock()
	img := p.frames[p.idxOut]
	if p.running || img == nil {
		p.Unlock()
		return
	}
	delete(p.frames, p.idxOut)
	p.running = true
	p.Unlock()

	for {
		p.proc.AddFrame(img)
		p.idxOut++
		p.wg.Add(-1)
		p.Lock()
		p.pool = append(p.pool, img)
		img = p.frames[p.idxOut]
		if img == nil {
			p.running = false
			p.Unlock()
			return
		}
		delete(p.frames, p.idxOut)
		p.Unlock()
	}
}
