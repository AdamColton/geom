package d2listwrap

import "github.com/adamcolton/geom/d2/d2list"

type Updater interface {
	Update()
}

func UpdateCurveList(cl d2list.CurveList) {
	ln := cl.Len()
	for i := 0; i < ln; i++ {
		c := cl.Idx(i)
		if u, ok := c.(Updater); ok {
			u.Update()
		}
	}
}
