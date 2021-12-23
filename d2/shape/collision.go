package shape

import (
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
	"github.com/adamcolton/geom/d2/shape/box"
)

func Collision(s1, s2 Shape) bool {
	iPts := make([]d2.Pt, 0, 4)
	b1 := box.New(s1.BoundingBox())
	b2 := box.New(s2.BoundingBox())

	// Should this logic get moved to package box.
	for i := 0; i < 4; i++ {
		if p := b1.Vertex(i); b2.Contains(p) {
			iPts = append(iPts, p)
		}
		if p := b2.Vertex(i); b1.Contains(p) {
			iPts = append(iPts, p)
		}
	}
	if len(iPts) == 0 {
		return false
	}

	return iterateBoxDescent(s1, s2, box.New(iPts...), iPts[:0], make([]float64, 8), 10)
}

func iterateBoxDescent(s1, s2 Shape, b *box.Box, pBuf []d2.Pt, fBuf []float64, steps int) bool {
	c := b.Centroid()
	if s1.Contains(c) && s2.Contains(c) {
		return true
	}
	if steps <= 0 {
		return false
	}

	var sa, sb Shape
	for i := 0; i < 4; i++ {
		l := b.Side(i)
		p := line.Line{
			D: d2.V{
				X: -l.D.Y,
				Y: l.D.X,
			},
		}
		for sIdx := 0; sIdx < 2; sIdx++ {
			if sIdx == 0 {
				sa, sb = s1, s2
			} else {
				sa, sb = s2, s2
			}

			fBuf = sa.LineIntersections(l, fBuf)
			for _, t := range fBuf {
				if t >= 0 && t < 1.0 {
					p.T0 = l.Pt1(t)
					for _, t2 := range sb.LineIntersections(p, fBuf[len(fBuf):]) {
						if t2 >= 0 && t2 < 1.0 {
							pBuf = append(pBuf, p.Pt1(t2))
						}
					}
				}
			}
		}
	}

	if len(pBuf) < 4 {
		return false
	}

	return iterateBoxDescent(s1, s2, box.New(pBuf...), pBuf[:0], fBuf, steps-1)
}

// FindEdge finds the perimeter point closest to the center of the box
func FindEdge(s Shape, b *box.Box, steps int, pBuf []d2.Pt, fBuf []float64) d2.Pt {

}
