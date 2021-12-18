package d3

type Point interface {
	Pt() Pt
}

type Vector interface {
	V() V
}

type Pt1 interface {
	Pt1(t0 float64) Pt
}

type V1 interface {
	V1(t0 float64) V
}

type Pt2 interface {
	Pt2(t0, t1 float64) Pt
}

type Pt2c1 interface {
	Pt2(t0 float64) Pt1
}

// TGen represents a type that can generate a Transform.
type TGen interface {
	T() *T
}

// TGenInv represents a type that can generate the inverse of a Transform.
type TGenInv interface {
	TGen
	TInv() *T
}

// TGenPair provides a way to get both the Transform and it's Inverse at the
// same time which can sometimes be more efficient.
type TGenPair interface {
	TGen
	Pair() [2]*T
}

// GetTInv of a TGen will call TInv if available.
func GetTInv(t TGen) *T {
	if p, ok := t.(TGenInv); ok {
		return p.TInv()
	}
	return t.T().TInv()
}

// GetTPair of a TGen will call Pair if available.
func GetTPair(t TGen) [2]*T {
	if p, ok := t.(TGenPair); ok {
		return p.Pair()
	}
	return [2]*T{t.T(), GetTInv(t)}
}
