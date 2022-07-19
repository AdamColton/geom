package grid

type Aspect float64

// Common aspect rations
const (
	Widescreen Aspect = 9.0 / 16.0
	Fullscreen Aspect = 3.0 / 4.0
	Square     Aspect = 1.0
)

// Pt is created from an aspect ratio and a width.
func (a Aspect) Pt(w int) Pt {
	return Pt{
		X: w,
		Y: int(Aspect(w) * a),
	}
}
