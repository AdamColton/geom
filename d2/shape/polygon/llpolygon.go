package polygon

import (
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
)

// LLNode is a node in a linked list polygon
type LLNode struct {
	PIdx, NextIdx uint32
}

// LLPolygon represents a polygon as a linked list stored in a slice. Start can
// be any index that is in the list. The list should be circular. This is useful
// when adding and removing vertexes in an algorithm.
type LLPolygon struct {
	Pts   []d2.Pt
	Nodes []LLNode
	Start uint32
}

// NewLL creates a linked list polygon from a regular polygon
func NewLL(p Polygon) *LLPolygon {
	ll := &LLPolygon{
		Pts:   p,
		Nodes: make([]LLNode, len(p)),
	}
	for i := range p {
		ll.Nodes[i] = LLNode{
			PIdx:    uint32(i),
			NextIdx: uint32(i + 1),
		}
	}
	ll.Nodes[len(p)-1].NextIdx = 0
	return ll
}

// Contains checks if a point in inside a linked list polygon
func (p *LLPolygon) Contains(pt d2.Pt) bool {
	// http://geomalgorithms.com/a03-_inclusion.html
	wn := 0
	curNd := p.Nodes[p.Start]
	for {
		nextNd := p.Nodes[curNd.NextIdx]
		cur, next := p.Pts[curNd.PIdx], p.Pts[nextNd.PIdx]
		c := line.New(cur, next).Cross(pt)
		if c == 0 &&
			((pt.X >= next.X && pt.X <= cur.X) || (pt.X <= next.X && pt.X >= cur.X)) &&
			((pt.Y >= next.Y && pt.Y <= cur.Y) || (pt.Y <= next.Y && pt.Y >= cur.Y)) {
			return true
		} else if cur.Y <= pt.Y {
			if c > 0 && next.Y > pt.Y {
				wn++
			}
		} else if c < 0 && next.Y <= pt.Y {
			wn--
		}
		if curNd.NextIdx == p.Start {
			break
		}
		curNd = nextNd
	}
	return wn != 0
}

// DoesIntersect returns true if the line intersects any side with a parametric
// value between 0 and 1.
func (p *LLPolygon) DoesIntersect(ln line.Line) bool {
	curNd := p.Nodes[p.Start]
	for {
		nextNd := p.Nodes[curNd.NextIdx]

		side := line.New(p.Pts[curNd.PIdx], p.Pts[nextNd.PIdx])
		i0, i1, ok := side.Intersection(ln)
		if ok && i0 >= small && i0 < 1.0-small && i1 >= 0 && i1 < 1.0 {
			return true
		}

		if curNd.NextIdx == p.Start {
			break
		}
		curNd = nextNd
	}
	return false
}
