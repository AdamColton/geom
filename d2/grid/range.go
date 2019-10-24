package grid

// Range represents a grid with the first point being inclusive and the second
// point exclusive.
type Range [2]Pt

// Contains checks if the point is in the Range with the first point being
// inclusive and the second point being exclusive.
func (r Range) Contains(pt Pt) bool {
	if r[0].X < r[1].X {
		if pt.X < r[0].X || pt.X >= r[1].X {
			return false
		}
	} else if r[1].X < r[0].X {
		if pt.X > r[0].X || pt.X <= r[1].X {
			return false
		}
	} else {
		return false
	}

	if r[0].Y < r[1].Y {
		if pt.Y < r[0].Y || pt.Y >= r[1].Y {
			return false
		}
	} else if r[1].Y < r[0].Y {
		if pt.Y > r[0].Y || pt.Y <= r[1].Y {
			return false
		}
	} else {
		return false
	}

	return true
}

// Size returns a Pt that indicates the width and height of the Range.
func (r Range) Size() Pt {
	return r[1].Subtract(r[0]).Abs()
}

// Iter fulfills IteratorFactory and returns a Scanner using the Range.
func (r Range) Iter() Iterator {
	return BaseIteratorWrapper{NewScanner(r)}
}

// Min returns the point with the lowest X and Y value that is contained in the
// range.
func (r Range) Min() Pt {
	var pt Pt

	if r[0].X < r[1].X {
		pt.X = r[0].X
	} else {
		pt.X = r[1].X + 1
	}

	if r[0].Y < r[1].Y {
		pt.Y = r[0].Y
	} else {
		pt.Y = r[1].Y + 1
	}

	return pt
}

// Max returns the point with the largest X and Y value that is contained in the
// range.
func (r Range) Max() Pt {
	var pt Pt

	if r[0].X > r[1].X {
		pt.X = r[0].X
	} else {
		pt.X = r[1].X - 1
	}

	if r[0].Y > r[1].Y {
		pt.Y = r[0].Y
	} else {
		pt.Y = r[1].Y - 1
	}

	return pt
}

// Scale the range so that r[0] returns t0=0, t1=0 and r[1] returns t0=1, t1=1.
func (r Range) Scale() Scale {
	var s Scale
	d := r[1].Subtract(r[0])
	if d.X > 1 {
		s.X = 1.0 / float64(d.X-1)
		s.DX = -float64(r[0].X) * s.X
	} else if d.X < -1 {
		s.X = 1.0 / float64(d.X+1)
		s.DX = -float64(r[0].X) * s.X
	} else {
		s.DX = 1
	}
	if d.Y > 1 {
		s.Y = 1.0 / float64(d.Y-1)
		s.DY = -float64(r[0].Y) * s.Y
	} else if d.Y < -1 {
		s.Y = 1.0 / float64(d.Y+1)
		s.DY = -float64(r[0].Y) * s.Y
	} else {
		s.DY = 1
	}
	return s
}
