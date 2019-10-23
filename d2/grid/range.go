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
	}

	if r[0].Y < r[1].Y {
		if pt.Y < r[0].Y || pt.Y >= r[1].Y {
			return false
		}
	} else if r[1].Y < r[0].Y {
		if pt.Y > r[0].Y || pt.Y <= r[1].Y {
			return false
		}
	}

	return true
}

func (r Range) Size() Pt {
	return r[1].Subtract(r[0]).Abs()
}

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
