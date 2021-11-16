package barycentric

import (
	"strconv"
	"strings"

	"github.com/adamcolton/geom/calc/cmpr"
	"github.com/adamcolton/geom/geomerr"
)

// B is a barycentric coordinate
type B struct {
	U, V float64
}

// Edge will return true if B is within d of an edge.
func (b B) Edge(d float64) bool {
	if b.U < -d || b.V < -d {
		return false
	}
	if b.U < d || b.V < d {
		return true
	}
	sum := 1 - b.U - b.V
	return sum < d && sum > -d
}

// Inside returns true if B is inside it's triangle.
func (b B) Inside() bool {
	return b.U >= 0 && b.V >= 0 && b.U+b.V <= 1
}

// String fulfills Stringer
func (b B) String() string {
	return strings.Join([]string{
		"B(",
		strconv.FormatFloat(b.U, 'f', 2, 64),
		", ",
		strconv.FormatFloat(b.V, 'f', 2, 64),
		")",
	}, "")
}

// BIterator iterates over points inside a triangle. A single iterator can be
// used to map between multiple triangles.
type BIterator struct {
	// Origin and U are the indexes to use, V is computed from these.
	Origin, U int
	Step      [2]B
	Cur       B
	// Reset is the point to reset to after one "line scan" completes
	Reset B
	Idx   int
}

// TODO: BIterator needs a significant change. It should always scan the whole
// triangle. The math for this will get a little tricky. It needs to work out
// the orientation of the scan and start in the correct corner and when doing a
// Step[1], it may need to take multiple Step[0] to be inside the triangle.

// V returns the index of the V coordinate
func (bs *BIterator) V() int { return 3 - (bs.Origin ^ bs.U) }

// Next returns the next B and if it is done
func (bs *BIterator) Next() (b B, done bool) {
	bs.Idx++
	bs.Cur.U += bs.Step[0].U
	bs.Cur.V += bs.Step[0].V
	if bs.Cur.Inside() {
		return bs.Cur, false
	}
	bs.Reset.U += bs.Step[1].U
	bs.Reset.V += bs.Step[1].V
	bs.Cur = bs.Reset
	return bs.Cur, !bs.Cur.Inside()
}

// Start resets the iterator
func (bs *BIterator) Start() (b B, done bool) {
	if bs == nil {
		return B{}, true
	}
	bs.Cur = B{}
	bs.Reset = B{}
	bs.Idx = 0
	return bs.Cur, false
}

// AssertEqual fulfils geomtest.AssertEqualizer
func (b B) AssertEqual(actual interface{}, t cmpr.Tolerance) error {
	b2, ok := actual.(B)
	if !ok {
		return geomerr.TypeMismatch(b, actual)
	}
	if !t.Zero(b.U-b2.U) || !t.Zero(b.V-b2.V) {
		return geomerr.NotEqual(b, b2)
	}
	return nil
}
