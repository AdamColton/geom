package grid

// ScanOption allows the scanner operation to be configured.
type ScanOption byte

const (
	// ScanVertical will scan columns then rows
	ScanVertical ScanOption = 1
	// ScanPrimaryReversed causes the scan in the primary direction to proceed
	// backward. The primary direction is X unless ScanVertical is used.
	ScanPrimaryReversed ScanOption = 2
	// ScanSecondaryReversed causes the scan in the secondary direction to
	// proceed backward. The secondary direction is Y unless ScanVertical is
	// used.
	ScanSecondaryReversed ScanOption = 4
)

func (s ScanOption) is(f ScanOption) bool {
	return s&f == f
}

// Scanner scans a rectangular region.
type Scanner struct {
	d     [2]Pt
	s     Pt
	r     Range
	ptIdx [3]int
	cur   struct {
		Pt
		Idx  int
		Done bool
	}
}

// == How Scanner Works ==
// d[0] : Will be one of the 4 cardinal directions, this is the primary scanning step.
// d[1] : This the step taken at the end of a row/col
// s : starting point, saved so we don't have to recalculate every time there's a reset
// r : The range
// ptIdx : Translating from a point to an index is linear, this holds a + bX + cY
// cur : the current or cursor value.
//
// The setup can be hard to follow because several operations are interspersed.
// Seeing them seperated out makes it easier. There are 8 possible scanning
// combinations. With no options, we Scan in X, then Y from Range[0] to
// Range[1]. The scan options can reverse those.
//
// The main if block at the end is the same logic repeat first if the primary
// scan direction is vertical then again if the primary scan direction is
// horizontal.
//
// The lines
//   pd, sd = pd*yd, sd*xd
// and
//   pd, sd = pd*xd, sd*yd
// are combining the primary scan direction, the range scan direction and the
// options to determine the primary scan direction (sd) and secondary scan
// direction (sd). The secondary direction point (s.d[1]) needs to reset the
// row or column hense
//   s.d[1].Y = (M.Y-m.Y)*(-pd) - pd
// or
//   s.d[1].X = (M.X-m.X)*(-pd) - pd

// NewScanner creates a scanner based on a Range. The r[0] value is incluse and
// the r[1] value is exclusive.
func NewScanner(r Range, opts ...ScanOption) *Scanner {
	var opt ScanOption
	for _, o := range opts {
		opt |= o
	}

	s := &Scanner{
		r: r,
	}
	m, M := r.Min(), r.Max()

	xd, yd := 1, 1
	if r[0].X > r[1].X {
		xd = -1
	}
	if r[0].Y > r[1].Y {
		yd = -1
	}

	pd, sd := 1, 1
	if opt.is(ScanPrimaryReversed) {
		pd = -1
	}
	if opt.is(ScanSecondaryReversed) {
		sd = -1
	}

	if opt.is(ScanVertical) {
		pd, sd = pd*yd, sd*xd

		s.d[0].Y = pd
		s.d[1].X = sd
		s.d[1].Y = (M.Y-m.Y)*(-pd) - pd
		s.ptIdx[2] = pd
		s.ptIdx[1] = -s.d[1].Y * sd * pd

		if pd == -1 {
			s.s.Y = M.Y
			s.ptIdx[0] = M.Y
		} else {
			s.s.Y = m.Y
			s.ptIdx[0] = -m.Y
		}

		if sd == -1 {
			s.s.X = M.X
			s.ptIdx[0] += pd * s.d[1].Y * -M.X
		} else {
			s.s.X = m.X
			s.ptIdx[0] += pd * s.d[1].Y * -m.X
		}
	} else {
		pd, sd = pd*xd, sd*yd

		s.d[0].X = pd
		s.d[1].Y = sd
		s.d[1].X = (M.X-m.X)*(-pd) - pd
		s.ptIdx[1] = pd
		s.ptIdx[2] = -s.d[1].X * sd * pd

		if pd == -1 {
			s.s.X = M.X
			s.ptIdx[0] = M.X
		} else {
			s.s.X = m.X
			s.ptIdx[0] = -m.X
		}

		if sd == -1 {
			s.s.Y = M.Y
			s.ptIdx[0] += pd * s.d[1].X * -M.Y
		} else {
			s.s.Y = m.Y
			s.ptIdx[0] += pd * s.d[1].X * -m.Y
		}
	}

	return s
}

// Pt the scanner is currently at
func (s *Scanner) Pt() Pt {
	return s.cur.Pt
}

// Idx is the Index value the scanner is currently at.
func (s *Scanner) Idx() int {
	return s.cur.Idx
}

// Done returns true when there is no more to scan
func (s *Scanner) Done() bool {
	return s.cur.Done
}

// Size returns a Pt that indicates the length and width of the scan area.
func (s *Scanner) Size() Pt {
	return s.r.Size()
}

// Reset moves the scanner back to the starting position.
func (s *Scanner) Reset() (done bool) {
	s.cur.Pt = s.s
	s.cur.Idx = 0
	s.cur.Done = !s.r.Contains(s.cur.Pt)
	return s.cur.Done
}

// Next moves the scanner to the next position
func (s *Scanner) Next() bool {
	s.cur.Idx++
	s.cur.Pt = s.cur.Pt.Add(s.d[0])
	if s.r.Contains(s.cur.Pt) {
		return false
	}

	s.cur.Pt = s.cur.Pt.Add(s.d[1])
	s.cur.Done = !s.r.Contains(s.cur.Pt)
	return s.cur.Done
}

// PtIdx takes a Pt and returns the Idx value when the scanner would reach
// that point. If the point is not inside the scan region, the response will not
// be meaningful.
func (s *Scanner) PtIdx(pt Pt) int {
	return s.ptIdx[0] + pt.X*s.ptIdx[1] + pt.Y*s.ptIdx[2]
}

// Contains checks if the Pt is inside the scanner region.
func (s *Scanner) Contains(pt Pt) bool {
	return s.r.Contains(pt)
}

// Iter makes a copy of the scanner and returns it as an Iterator.
func (s *Scanner) Iter() Iterator {
	cp := *s
	cp.Reset()
	return BaseIteratorWrapper{&cp}
}

// Scale returns the Scale that corresponds with the Range that defined the
// scanner.
func (s *Scanner) Scale() Scale {
	return s.r.Scale()
}
