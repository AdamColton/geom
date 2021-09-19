package material

import (
	"github.com/adamcolton/geom/angle"
)

type Material struct {
	Specular    angle.Rad
	Color       *Color
	Diffuse     angle.Rad
	Luminous    float64
	Border      float64
	BorderColor *Color
}
