package grid

// Generator is used to fill grids from their point values.
type Generator func(pt Pt) interface{}

// Grid is used to store data in a 2 dimensional grid.
type Grid interface {
	Get(Pt) (interface{}, error)
	Set(Pt, interface{}) error
	Iterator
}
