package zbuf

import (
	"math"
	"testing"

	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/render/scene"
	"github.com/adamcolton/geom/d3/shape/triangle"
	"github.com/adamcolton/geom/d3/solid/mesh"
	"github.com/adamcolton/geom/geomtest"
	"github.com/stretchr/testify/assert"
)

func TestCameraBasic(t *testing.T) {
	c := Camera{
		Camera: &scene.Camera{
			Pt:     d3.Pt{0, 0, 0},
			Q:      d3.Q{1, 0, 0, 0},
			Angle:  math.Pi / 2.0,
			Width:  1,
			Height: 1,
		},
		Near: 1,
		Far:  10,
	}

	assert.Equal(t, c.Q.T(), d3.Identity())

	ca, cb := c.ab()
	assert.Equal(t, c.Near, -ca*c.Near+cb)
	assert.Equal(t, -c.Far, -ca*c.Far+cb)

	tr := c.T()
	//assert.Equal(t, c.Perspective(), tr)

	//n, nw := tr.PtF(d3.Pt{0, 0, -c.Near})
	//assert.Equal(t, d3.Pt{0, 0, c.Near}, n)
	//assert.Equal(t, nw, c.Near)

	// since the perspective is square, edge is the distance in both x and y
	// from the z to the perimeter at the near point.
	edge := math.Tan(float64(c.Angle) / 2.0)

	testPoints := []d3.Pt{
		{0, 0, -c.Near},
		{0, 0, -c.Far},
		{edge, edge, -c.Near},
		{edge, -edge, -c.Near},
		{-edge, -edge, -c.Near},
		{-edge, edge, -c.Near},
		{10 * edge, 10 * edge, -c.Far},
		{10 * edge, 10 * -edge, -c.Far},
		{10 * -edge, 10 * -edge, -c.Far},
		{10 * -edge, 10 * edge, -c.Far},
	}

	expected := []d3.Pt{
		{0.5, 0.5, 0},
		{0.5, 0.5, 1},
		{1, 1, 0},
		{1, 0, 0},
		{0, 0, 0},
		{0, 1, 0},
		{1, 1, 1},
		{1, 0, 1},
		{0, 0, 1},
		{0, 1, 1},
	}

	for i, p := range testPoints {
		assert.Equal(t, expected[i], tr.PtScl(p))
	}
}

func TestCameraMesh(t *testing.T) {
	c := Camera{
		Camera: &scene.Camera{
			Pt:    d3.Pt{0, 0, 0},
			Q:     d3.Q{1, 0, 0, 0},
			Angle: math.Pi / 2.0,
		},
		Near: 2,
		Far:  10,
	}

	expected := make([]d3.Pt, 0, 8)
	f := []float64{1, -1}
	for _, z := range f {
		for _, y := range f {
			for _, x := range f {
				expected = append(expected, d3.Pt{x, y, z})
			}
		}
	}

	xy := c.Near * math.Tan(float64(c.Angle)/2.0)
	scale := c.Far / c.Near
	m := mesh.NewExtrusion(
		[]d3.Pt{
			{xy, xy, -c.Near},
			{-xy, xy, -c.Near},
			{xy, -xy, -c.Near},
			{-xy, -xy, -c.Near},
		}).
		Extrude(
			d3.Scale(d3.V{scale, scale, 1}).T().T(d3.Translate(d3.V{0, 0, -c.Far + c.Near}).T()),
		).
		Close()
	mt := m.T(c.T())

	for i, p := range mt.Pts {
		geomtest.Equal(t, expected[i], p)
	}
}

func TestCameraWH(t *testing.T) {
	c := Camera{
		Camera: &scene.Camera{
			Pt:     d3.Pt{0, 0, 0},
			Q:      d3.Q{1, 0, 0, 0},
			Angle:  math.Pi / 2.0,
			Width:  150,
			Height: 100,
		},
		Near: 1,
		Far:  10,
	}

	assert.Equal(t, c.Q.T(), d3.Identity())

	ca, cb := c.ab()
	assert.Equal(t, c.Near, -ca*c.Near+cb)
	assert.Equal(t, -c.Far, -ca*c.Far+cb)

	tr := c.T()
	//assert.Equal(t, c.Perspective(), tr)

	//n, nw := tr.PtF(d3.Pt{0, 0, -c.Near})
	//assert.Equal(t, d3.Pt{0, 0, c.Near}, n)
	//assert.Equal(t, nw, c.Near)

	edge := math.Tan(float64(c.Angle) / 2.0)
	y := float64(c.Height) / float64(c.Width)
	testPoints := []d3.Pt{
		{0, 0, -c.Near},
		{0, 0, -c.Far},
		{edge, y * edge, -c.Near},
		{edge, y * -edge, -c.Near},
		{-edge, y * -edge, -c.Near},
		{-edge, y * edge, -c.Near},
		{10 * edge, y * 10 * edge, -c.Far},
		{10 * edge, y * 10 * -edge, -c.Far},
		{10 * -edge, y * 10 * -edge, -c.Far},
		{10 * -edge, y * 10 * edge, -c.Far},
	}

	expected := []d3.Pt{
		{0.5, 0.5, 0},
		{0.5, 0.5, 1},
		{1, 1, 0},
		{1, 0, 0},
		{0, 0, 0},
		{0, 1, 0},
		{1, 1, 1},
		{1, 0, 1},
		{0, 0, 1},
		{0, 1, 1},
	}

	for i, p := range testPoints {
		assert.Equal(t, expected[i], tr.PtScl(p))
	}
}

// func TestCamera(t *testing.T) {
// 	t.Skip()
// 	fmt.Println("Generating Images")
// 	c := Camera{
// 		Pt:    d3.Pt{},
// 		Q:     d3.Q{1, 0, 0, 0},
// 		Near:  0.1,
// 		Far:   10,
// 		Angle: 3.1415 / 2.0,
// 	}

// 	s := d3.Pt{0, 0, 0}
// 	f := []d3.Pt{
// 		s,
// 		s.Add(d3.V{1, 1, 0}),
// 		s.Add(d3.V{2, 1, 0}),
// 		s.Add(d3.V{3, 0, 0}),
// 		s.Add(d3.V{0, -3, 0}),
// 		s.Add(d3.V{-3, 0, 0}),
// 		s.Add(d3.V{-2, 1, 0}),
// 		s.Add(d3.V{-1, 1, 0}),
// 	}

// 	m := mesh.NewExtrusion(f).
// 		Extrude(d3.Translate(d3.V{0.5, 0.5, -1.5}).T()).
// 		Close()

// 	fmt.Println(m.Pts)

// 	ref := d3.Pt{1, 0, -3}
// 	for i := 0.0; i < 21; i++ {
// 		//c.Pre = d3.NewTSet().AddBoth(d3.Translate(d3.V{0, 0, -3.5})).Add(d3.Q{1, i, 0, 0}.T()).Get()
// 		s, cs := math.Sincos(math.Pi * i / 10.0)
// 		r := d3.T{
// 			{cs, -s, 0, 0},
// 			{s, cs, 0, 0},
// 			{0, 0, 1, 0},
// 			{0, 0, 0, 1},
// 		}
// 		fmt.Println(r.Pt(ref))
// 		//c.Pre = r.T(d3.Translate(d3.V{0, 0, -3}).T())
// 		m := m.T(c.T())
// 		ctx := gg.NewContext(500, 500)
// 		ctx.SetRGB(0, 0, 0)
// 		ctx.DrawRectangle(0, 0, 500, 500)
// 		ctx.Fill()
// 		ctx.Stroke()
// 		ctx.SetLineWidth(4)
// 		ctx.SetRGB(1, 0, 0.1)
// 		tm, err := m.TriangleMesh()
// 		assert.NoError(t, err)
// 		Wireframe(ctx, tm, &([3]float64{1, 0, 0}), &([3]float64{0, 1, 0}))
// 		ctx.SavePNG("test" + strconv.FormatFloat(i, 'f', 0, 64) + ".png")
// 	}
// }

func TestScan(t *testing.T) {
	tr := &triangle.Triangle{{10, 10, 0}, {100, 50, 0}, {50, 100, 0}}
	bi, bt := Scan(tr, 0.5, 0.7)
	assert.NotNil(t, bi)
	assert.Equal(t, bt.Pt, d3.Pt{10, 10, 0})

	// m*U.X + n*V.X = dx
	// m*U.Y + n*V.Y = 0

	assert.Equal(t, 0.5, bt.U.X*bi.Step[0].U+bt.V.X*bi.Step[0].V)
	assert.Equal(t, 0.0, bt.U.Y*bi.Step[0].U+bt.V.Y*bi.Step[0].V)
}

func TestScanU(t *testing.T) {
	tt := map[string]struct {
		t         triangle.Triangle
		origin, u int
	}{
		"basic": {
			t:      triangle.Triangle{{0, 0, 0}, {1, 0, 0}, {0, 1, 0}},
			origin: 0,
			u:      2,
		},
		"invert": {
			t:      triangle.Triangle{{0, 0, 0}, {-1, 0, 0}, {0, -1, 0}},
			origin: 2,
			u:      0,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			bt := scanU(&(tc.t))
			assert.Equal(t, tc.origin, bt.Origin)
			assert.Equal(t, tc.u, bt.U)
		})
	}
}
