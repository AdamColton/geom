package boxmodel

const (
	unknownLeaf uint32 = iota
	insideLeaf
	outsideLeaf
	perimeterLeaf
	firstParent
)

// cIter is a helper because I prefer
//   for _, child := range cIter
// over
//   for child := byte(0); child < 4; child++
var cIter = [4]byte{0, 1, 2, 3}

type children [4]uint32
