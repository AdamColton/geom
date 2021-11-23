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

// bytes is use for compression, converting the children to a byte slice that
// can be hashed.
func (c children) bytes(out []byte) {
	for _, child := range cIter {
		b := c[child]
		for i := byte(0); i < 4; i++ {
			out[child*4+i] = byte(b)
			b >>= 8
		}
	}
}
