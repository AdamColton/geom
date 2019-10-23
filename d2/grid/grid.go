package grid

type Generator func(pt Pt) interface{}

type Grid interface {
	Get(Pt) (interface{}, error)
	Set(Pt, interface{}) error
	Iterator
}
