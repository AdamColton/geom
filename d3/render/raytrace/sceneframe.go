package raytrace

import (
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/geom/work"

	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/curve/line"
	"github.com/adamcolton/geom/d3/render/material"
	"github.com/adamcolton/geom/d3/render/scene"
	"github.com/adamcolton/geom/d3/shape/triangle"
	"github.com/adamcolton/geom/d3/solid/box"
)

type SceneFrame struct {
	*scene.SceneFrame
	*RayFrame
}

type RayFrame struct {
	Background   Shader
	Shaders      []Shader
	Depth        int
	RayMult      int
	ImageScale   float64
	Intersectors []ModelIntersector
}

func (sf *SceneFrame) PopulateShaders() {
	sf.Shaders = make([]Shader, len(sf.SceneFrame.Meshes))
	for i, m := range sf.SceneFrame.Meshes {
		shader, ok := m.Shader.(RayShader)
		if ok {
			sf.Shaders[i] = shader.RayShader
			continue
		}
		m, ok := m.Shader.(*material.Material)
		if ok {
			mw := *NewMaterialWrapper(*m)
			sf.Shaders[i] = mw.RayShader
		}
	}
}

type Intersection struct {
	Ray line.Line
	*scene.SceneFrame
	*triangle.Triangle
	*scene.TriangleIndex
	triangle.Intersection
}

type ModelIntersector struct {
	*box.Box
	Intersectors []Intersector
}

type Intersector struct {
	scene.TriangleIndex
	*triangle.TriangleIntersector
}

func (sf *SceneFrame) PopulateIntersectors() {
	sf.Intersectors = make([]ModelIntersector, len(sf.Meshes))
	t := &triangle.Triangle{}
	for mIdx, m := range sf.Meshes {
		mi := ModelIntersector{
			Box:          &box.Box{},
			Intersectors: make([]Intersector, m.Original.GetTriangleCount()),
		}
		idx := 0
		for pIdx, p := range m.Original.Polygons {
			for tIdx, ptIdxs := range p {
				t[0] = m.Space[ptIdxs[0]]
				t[1] = m.Space[ptIdxs[1]]
				t[2] = m.Space[ptIdxs[2]]
				mi.Box.Add(t[0], t[1], t[2])
				i := t.Intersector()
				mi.Intersectors[idx] = Intersector{
					TriangleIndex: scene.TriangleIndex{
						MeshIdx:     mIdx,
						PolygonIdx:  pIdx,
						TriangleIdx: tIdx,
					},
					TriangleIntersector: i,
				}
				idx++
			}
		}
		sf.Intersectors[mIdx] = mi
	}
}

func (sf *SceneFrame) intersect(ray line.Line) *Intersection {
	out := &Intersection{
		Ray:        ray,
		SceneFrame: sf.SceneFrame,
	}
	for _, mi := range sf.Intersectors {
		if _, doesIsect := mi.Box.LineIntersection(ray); !doesIsect {
			continue
		}
		for _, i := range mi.Intersectors {
			ri := i.RawIntersection(ray)
			// TODO: t0 check needs to come from a "near" cutoff value
			if ri.Does && ri.T > 0.001 && (out.TriangleIndex == nil || out.T > ri.T) {
				out.TriangleIndex = &i.TriangleIndex
				out.Intersection = ri
			}
		}
	}
	if out.TriangleIndex != nil {
		t := &triangle.Triangle{}
		out.Triangle = t
		m := sf.Meshes[out.MeshIdx]
		p := m.Original.Polygons[out.PolygonIdx]
		ptIdxs := p[out.TriangleIdx]
		t[0] = m.Space[ptIdxs[0]]
		t[1] = m.Space[ptIdxs[1]]
		t[2] = m.Space[ptIdxs[2]]
	}
	return out
}

func (sf *SceneFrame) trace(ray line.Line, depth int) *material.Color {
	ctx := &Context{
		Intersection: sf.intersect(ray),
	}
	if ctx.TriangleIndex == nil {
		return sf.Background(ctx).Color
	}
	m := sf.Shaders[ctx.MeshIdx](ctx)
	if m.Luminous > 0 {
		return m.Color.Scale(m.Luminous)
	}

	pt := ctx.Ray.Pt1(ctx.T)
	rv := reflect(ctx.Ray.D, ctx.Triangle.Normal().Normal())
	colors := make([]*material.Color, sf.RayMult*depth)
	for i := range colors {
		q := randomAngle(m.Diffuse)
		ray := line.Line{
			T0: pt,
			D:  q.T().V(rv),
		}
		colors[i] = sf.trace(ray, depth-1)
	}
	c := material.Avg(colors...)
	c = material.Reflect(m.Color, c)
	return c
}

func (sf *SceneFrame) Image(img *image.RGBA) *image.RGBA {
	if sf.Shaders == nil {
		sf.PopulateShaders()
	}
	if sf.Intersectors == nil {
		sf.PopulateIntersectors()
	}

	size := sf.Camera.Size.Multiply(sf.ImageScale)
	z := -(0.5 * float64(size.X)) / math.Tan(float64(sf.Camera.Angle)*0.5)
	t, ok := sf.Camera.Rot.Inversion()
	if !ok {
		panic("cannot find inversion")
	}
	t = t.T(d3.Translate(sf.Camera.Pt).TInv())

	if img == nil {
		img = image.NewRGBA(image.Rect(0, 0, size.X, size.Y))
	}
	ln := size.Area()
	//w64, h64 := float64(sf.Width), float64(sf.Height)
	d := size.Multiply(0.5)

	work.RunRange(ln, func(idx, _ int) {
		pt := size.Index(idx).Subtract(d)
		iPt := pt.Add(d)
		iPt.Y = size.Y - iPt.Y - 1
		// if iPt.X < 225 || iPt.X > 275 || iPt.Y < 200 || iPt.Y > 220 {
		// 	return
		// }

		v := t.V(d3.V{float64(pt.X), float64(pt.Y), z})
		l := line.Line{
			T0: sf.Camera.Pt,
			D:  v,
		}

		c := sf.trace(l, sf.Depth)
		img.SetRGBA(iPt.X, iPt.Y, color.RGBA{uint8(c.R * 255), uint8(c.G * 255), uint8(c.B * 255), 255})
	})

	return img
}

// reflect a ray off a surface. The ray is represented by V and the surface is
// represented by n - it's normal vector.
func reflect(v, n d3.V) d3.V {
	return v.Subtract(n.Multiply(2 * v.Dot(n)))
}

// randomAngle returns a Quaternion that deviates from the identity Quaternion
// by no more than the angle.
func randomAngle(ang angle.Rad) d3.Q {
	r1 := angle.Rad(ang.Rad() * rand.Float64() * 0.5)
	sx, cx := r1.Sincos()
	r2 := angle.Rot(rand.Float64() * 0.5)
	sy, cy := r2.Sincos()

	return d3.Q{cx, sx * cy, sx * sy, 0}
}
