package material

import (
	"github.com/adamcolton/geom/angle"
)

// Material holds information that different renderers can use to shade a
// material.
type Material struct {
	Specular    angle.Rad
	Color       *Color
	Diffuse     angle.Rad
	Luminous    float64
	Border      float64
	BorderColor *Color
}
