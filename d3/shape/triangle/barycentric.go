package triangle

import (
	"github.com/adamcolton/geom/barycentric"
	"github.com/adamcolton/geom/d3"
)

// BT is a triangle expressed barycentrically
type BT struct {
	Pt   d3.Pt
	U, V d3.V
}

// BT translates a Triangle into a Barycentric representation of a triangle
func (t Triangle) BT(origin, u int, buf *BT) *BT {
	if origin < 0 || origin > 2 || u < 0 || u > 2 || origin == u {
		return nil
	}
	v := 3 - (origin ^ u)
	if buf == nil {
		return &BT{
			Pt: t[origin],
			U:  t[u].Subtract(t[origin]),
			V:  t[v].Subtract(t[origin]),
		}
	}
	buf.Pt = t[origin]
	buf.U = t[u].Subtract(t[origin])
	buf.V = t[v].Subtract(t[origin])
	return buf
}

// PtB translates a barycentric coordinate to a d3 Pt
func (bt *BT) PtB(b barycentric.B) d3.Pt {
	return bt.Pt.Add(bt.U.Multiply(b.U)).Add(bt.V.Multiply(b.V))
}

// Triangle translates from a Barycentric representation to 3 d3.Pts
func (bt *BT) Triangle() *Triangle {
	return &Triangle{
		bt.Pt,
		bt.Pt.Add(bt.U),
		bt.Pt.Add(bt.V),
	}
}
