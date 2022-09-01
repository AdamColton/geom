package polygon

import (
	"math"
	"strings"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/curve/line"
)

// Polygon represents a Convex Polygon
type Polygon []d2.Pt

// New creates a polygon and orders the points to proceed counter clockwise.
func New(pts []d2.Pt) Polygon {
	p, c := NewPolar(pts)
	return p.Polygon(c)
}

// Pt2c1 returns line.Segments as d2.Pt1 that adheres to the Shape rules
func (p Polygon) Pt2c1(t0 float64) d2.Pt1 {
	n := (len(p) - 1)
	h := n - (n / 2) // half rounded up
	as := line.Segments(p[0:h])
	bs := line.Segments(p[h:])
	a := as.Pt1(t0)
	b := bs.Pt1(1 - t0)
	return line.New(a, b)
}

// Pt2 returns a point in the polygon and adheres to the Shape rules.
func (p Polygon) Pt2(t0, t1 float64) d2.Pt {
	return p.Pt2c1(t0).Pt1(t1)
}

// Pt1 returns a point on the perimeter.
func (p Polygon) Pt1(t0 float64) d2.Pt {
	ln := len(p)
	t0 -= math.Floor(t0) // [0,1)
	t0 *= float64(ln)
	idx := int(t0)
	t0 -= float64(idx)
	return line.New(p[idx], p[(idx+1)%ln]).Pt1(t0)
}

// String lists the points as a string.
func (p Polygon) String() string {
	pts := make([]string, len(p))
	for i, pt := range p {
		pts[i] = pt.String()
	}

	return strings.Join([]string{"Polygon{ ", strings.Join(pts, ", "), " }"}, "")
}

// SignedArea returns the Area and may be negative depending on the polarity.
func (p Polygon) SignedArea() float64 {
	var s float64
	prev := d2.V(p[len(p)-1])
	for _, cur := range p {
		v := d2.V(cur)
		s += prev.Cross(v)
		prev = v
	}
	return s / 2
}

// Area of the polygon
func (p Polygon) Area() float64 {
	return math.Abs(p.SignedArea())
}

// Centroid returns the center of mass of the polygon
func (p Polygon) Centroid() d2.Pt {
	var x, y, a float64
	prev := p[len(p)-1]
	for _, cur := range p {
		t := (prev.X*cur.Y - cur.X*prev.Y)
		x += (prev.X + cur.X) * t
		y += (prev.Y + cur.Y) * t
		a += t
		prev = cur
	}
	a = 1 / (3 * a)
	return d2.Pt{x * a, y * a}
}

// Contains returns true of the point f is inside of the polygon. This is done
// with the winding number algorithm and runs in O(n).
func (p Polygon) Contains(pt d2.Pt) bool {
	// https://en.wikipedia.org/wiki/Point_in_polygon#Winding_number_algorithm
	windings := 0
	prev := p[len(p)-1]
	for _, cur := range p {
		c := line.New(prev, cur).Cross(pt)
		if c == 0 &&
			((pt.X >= prev.X && pt.X <= cur.X) || (pt.X <= prev.X && pt.X >= cur.X)) &&
			((pt.Y >= prev.Y && pt.Y <= cur.Y) || (pt.Y <= prev.Y && pt.Y >= cur.Y)) {
			return true
		} else if prev.Y <= pt.Y {
			if c > 0 && cur.Y > pt.Y {
				windings++
			}
		} else if c < 0 && cur.Y <= pt.Y {
			windings--
		}
		prev = cur
	}
	return windings != 0
}

// Perimeter returngs the total length of the perimeter
func (p Polygon) Perimeter() float64 {
	var sum float64
	prev := p[0]
	for _, f := range p[1:] {
		sum += prev.Distance(f)
		prev = f
	}
	sum += prev.Distance(p[0])
	return sum
}

// CountAngles returns the number of counter clockwise and clockwise angles. If
// ccwOut or cwOut is not nil, they will be populated with the indexes of the
// verticies
func (p Polygon) CountAngles(out []int) (ccw int, cw int) {
	if out != nil && len(out) != len(p) {
		out = nil
	}
	prevF := p[len(p)-1]
	prevAngle := p[len(p)-2].Subtract(prevF).Angle()
	for idx, f := range p {
		curAngle := prevF.Subtract(f).Angle()
		a := curAngle - prevAngle
		if a < 0 {
			a += 2 * math.Pi
		}
		if a > math.Pi {
			if out != nil {
				out[len(out)-cw-1] = idx
			}
			cw++
		} else if a < math.Pi {
			if out != nil {
				out[ccw] = idx
			}
			ccw++
		}
		prevAngle, prevF = curAngle, f
	}
	return
}

// Convex returns True if the polygon contains a convex angle.
func (p Polygon) Convex() bool {
	ccw, cw := p.CountAngles(nil)
	return ccw == 0 || cw == 0
}

const small = 1e-5

// FindTriangles returns the index sets of the polygon broken up into triangles.
// Given a unit square it would return [[0,1,2], [0,2,3]] which means that
// the square can be broken up in to 2 triangles formed by the points at those
// indexes.
func (p Polygon) FindTriangles() [][3]uint32 {
	// This is the ear clipping, there are better algorithms
	out := make([][3]uint32, 0, len(p)-2)

	ll := NewLL(p)
	left := len(p)

	n0 := &(ll.Nodes[0])
	n1 := &(ll.Nodes[1])
	n2 := &(ll.Nodes[2])

	cur := make(Polygon, len(p))
	copy(cur, p)

	ctr := 0
	for ctr < len(p) {
		ctr++
		n0, n1, n2 = n1, n2, &(ll.Nodes[n2.NextIdx])
		ll.Start = n0.NextIdx
		if left == 3 {
			out = append(out, [3]uint32{n0.PIdx, n1.PIdx, n2.PIdx})
			break
		}

		ln := line.New(ll.Pts[n0.PIdx], ll.Pts[n2.PIdx])
		if !ll.Contains(ln.Pt1(0.5)) || ll.DoesIntersect(ln) {
			continue
		}
		ctr = 0
		left--
		out = append(out, [3]uint32{n0.PIdx, n1.PIdx, n2.PIdx})
		n0.NextIdx = n1.NextIdx
		n1, n2 = n2, &(ll.Nodes[n2.NextIdx])
	}
	return out
}

// Collision returns the first side that is intersected by the given
// lineSegment, returning the parametic t for the lineSegment, the index of the
// side and the parametric t of the side
func (p Polygon) Collision(lineSegment line.Line) (lineT float64, idx int, sideT float64) {
	idx = -1
	ln := len(p)
	for i, f := range p {
		side := line.New(f, p[(i+1)%ln])
		t0, t1, ok := line.DefaultRange.Check(side.Intersection(lineSegment))
		if ok && (idx == -1 || lineT > t1) {
			lineT = t1
			idx = i
			sideT = t0
		}
	}
	return
}

// LineIntersections fulfills line.LineIntersector
func (p Polygon) LineIntersections(ln line.Line, buf []float64) []float64 {
	max := len(buf)
	buf = buf[:0]
	prev := p[len(p)-1]
	for _, cur := range p {
		side := line.New(prev, cur)
		t0, t1, ok := ln.Intersection(side)
		if ok && t1 >= 0 && t1 < 1 {
			buf = append(buf, t0)
			if max > 0 && len(buf) == max {
				return buf
			}
		}
		prev = cur
	}
	return buf
}

// Collision between two polygons
type Collision struct {
	PIdx, P2Idx int
	PT, P2T     float64
}

// P returns the collision point using polygon P to compute.
func (c Collision) P(p Polygon) d2.Pt {
	return p.Side(c.PIdx).Pt1(c.PT)
}

// P2 returns the collision point using polygon P2 to compute.
func (c Collision) P2(p2 Polygon) d2.Pt {
	return p2.Side(c.P2Idx).Pt1(c.P2T)
}

// Collisions represent a set of collisions between two polygons.
type Collisions []Collision

// P returns the collision points using polygon P to compute.
func (cs Collisions) P(p Polygon) []d2.Pt {
	out := make([]d2.Pt, 0, len(cs))
	for _, c := range cs {
		out = append(out, c.P(p))
	}
	return out
}

// P2 returns the collision points using polygon P2 to compute.
func (cs Collisions) P2(p2 Polygon) []d2.Pt {
	out := make([]d2.Pt, 0, len(cs))
	for _, c := range cs {
		out = append(out, c.P2(p2))
	}
	return out
}

// PolygonIntersections finds the intersection points between two polygons.
func (p Polygon) PolygonCollisions(p2 Polygon) Collisions {
	ln := len(p)
	ln2 := len(p2)
	if ln < 2 || ln2 < 2 {
		return nil
	}
	sides := p.Sides()
	sides2 := p2.Sides()
	var out Collisions
	for idx, s := range sides {
		for idx2, s2 := range sides2 {
			t, t2, ok := line.DefaultRange.Check(s.Intersection(s2))
			if ok {
				out = append(out, Collision{
					PIdx:  idx,
					P2Idx: idx2,
					PT:    t,
					P2T:   t2,
				})
			}
		}
	}
	return out
}

// Sides converts the perimeter of the polygon to a slice of lines.
func (p Polygon) Sides() []line.Line {
	ln := len(p)
	side := make([]line.Line, ln)
	for i, pt := range p[:ln-1] {
		side[i] = line.New(pt, p[(i+1)])
	}
	side[ln-1] = line.New(p[ln-1], p[0])
	return side
}

// Side n of the polygon where the side is formed by p[n] to p[n+1].
func (p Polygon) Side(n int) line.Line {
	ln := len(p)
	n %= ln
	if n < 0 {
		n += ln
	}
	return line.New(p[n], p[(n+1)%ln])
}

// NonIntersecting returns false if any two sides intersect. This requires
// O(N^2) time to check.
func (p Polygon) NonIntersecting() bool {
	side := p.Sides()
	// Each side needs to be check against each non-adjacent side with a greater
	// index.
	for i, si := range side[:len(side)-2] {
		for _, sj := range side[i+2:] {
			t0, t1, ok := si.Intersection(sj)
			if ok &&
				t1 > 0 && t1 < 1 &&
				t0 > 0 && t0 < 1 {
				return false
			}
		}
	}
	return true
}

// Reverse the order of the points defining the polygon
func (p Polygon) Reverse() Polygon {
	out := make([]d2.Pt, len(p))
	l := len(p) - 1
	m := (len(p) + 1) / 2 //+1 causes round up
	for i := 0; i < m; i++ {
		out[i], out[l-i] = p[l-i], p[i]
	}
	return out
}

// BoundingBox fulfills BoundingBoxer returning a box that contains the polygon.
func (p Polygon) BoundingBox() (min, max d2.Pt) {
	return d2.MinMax(p...)
}

// ConvexHull fulfills shape.ConvexHuller. Returns the convex hull of the
// polygon using the ConvexHull function.
func (p Polygon) ConvexHull() []d2.Pt {
	return ConvexHull(p...)
}

// T applies the transform to the Polygon returning a new Polygon.
func (p Polygon) T(transform *d2.T) Polygon {
	return transform.Slice(p)
}
