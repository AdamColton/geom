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

type TGen interface {
	T() *T
	TInv() *T
	Pair() [2]*T
}
