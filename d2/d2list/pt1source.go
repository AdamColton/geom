package d2list

import (
	"github.com/adamcolton/geom/d2"
)

type Pt1Source struct {
	d2.Pt1
	Start, D float64
	N        int
}

func (ps Pt1Source) Len() int {
	return ps.N
}

func (ps Pt1Source) Get(n int) d2.Pt {
	return ps.Pt1.Pt1(ps.Start + ps.D*float64(n))
}

func (ps *Pt1Source) Set(start, end float64, n int) {
	ps.Start = start
	ps.N = n
	ps.D = (end - start) / float64(n-1)
}
