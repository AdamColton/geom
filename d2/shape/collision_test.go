package shape

import (
	"testing"

	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/shape/ellipse"
	"github.com/adamcolton/geom/d2/shape/triangle"
	"github.com/stretchr/testify/assert"
)

func TestCollision(t *testing.T) {
	tt := map[string]struct {
		s1, s2   Shape
		expected bool
	}{
		"box-no-intersection": {
			s1:       ellipse.NewCircle(d2.Pt{}, 1),
			s2:       ellipse.NewCircle(d2.Pt{2.1, 0}, 1),
			expected: false,
		},
		"near-miss": {
			s1:       ellipse.NewCircle(d2.Pt{}, 0.707),
			s2:       ellipse.NewCircle(d2.Pt{1, 1}, 0.707),
			expected: false,
		},
		"small-overlap": {
			s1:       ellipse.NewCircle(d2.Pt{}, 1.12),
			s2:       ellipse.NewCircle(d2.Pt{2, 1}, 1.12),
			expected: true,
		},
		"triangle-circle": {
			s1:       ellipse.NewCircle(d2.Pt{}, 1),
			s2:       &triangle.Triangle{{-1, .999}, {1, .999}, {0, 2}},
			expected: true,
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc.expected, Collision(tc.s1, tc.s2))
		})
	}
}
