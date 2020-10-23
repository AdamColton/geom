package raytrace

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/adamcolton/geom/angle"

	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/curve/line"
	"github.com/adamcolton/geom/d3/render/scene"
	"github.com/adamcolton/geom/d3/shape/triangle"
)

type Context struct {
	*Intersection
}

type Shader func(*Context) *Material

type RayShader interface {
	RayShader(*Context) *Material
}

type SceneFrame struct {
	*scene.SceneFrame
	Background Shader
	Depth      int
	RayMult    int
	Shaders    []Shader
}

func (sf *SceneFrame) PopulateShaders() {
	sf.Shaders = make([]Shader, len(sf.SceneFrame.Meshes))
	for i, m := range sf.SceneFrame.Meshes {
		shader, _ := m.Shader.(RayShader)
		sf.Shaders[i] = shader.RayShader
	}
}

type Intersection struct {
	Ray line.Line
	*scene.SceneFrame
	*triangle.Triangle
	MeshIdx       int
	PolygonIndex  int
	TriangleIndex int
	T             float64
}

func (sf *SceneFrame) intersect(ray line.Line) *Intersection {
	out := &Intersection{
		Ray:        ray,
		SceneFrame: sf.SceneFrame,
		MeshIdx:    -1,
	}
	for mIdx, m := range sf.Meshes {
		for pIdx, p := range m.Original.Polygons {
			for tIdx, ptIdxs := range p {
				t := &triangle.Triangle{
					m.Space[ptIdxs[0]],
					m.Space[ptIdxs[1]],
					m.Space[ptIdxs[2]],
				}
				t0, ok := t.Intersection(ray)
				if ok && t0 > 1e-5 && (out.Triangle == nil || out.T > t0) {
					out.MeshIdx = mIdx
					out.PolygonIndex = pIdx
					out.TriangleIndex = tIdx
					out.Triangle = t
					out.T = t0
				}
			}
		}
	}
	return out
}

func (sf *SceneFrame) trace(ray line.Line, depth int) *Color {
	ctx := &Context{
		Intersection: sf.intersect(ray),
	}
	if ctx.MeshIdx == -1 {
		return sf.Background(ctx).Color
	}
	m := sf.Shaders[ctx.MeshIdx](ctx)
	//return m.Color
	if m.Luminous > 0 {
		return m.Color.Scale(m.Luminous)
	}

	pt := ctx.Ray.Pt1(ctx.T)
	rv := reflect(ctx.Ray.D, ctx.Triangle.Normal().Normal())
	colors := make([]*Color, sf.RayMult*depth)
	for i := range colors {
		q := randomAngle(m.Diffuse)
		t := q.T()
		ray := line.Line{
			T0: pt,
			D:  t.V(rv),
		}
		colors[i] = sf.trace(ray, depth-1)
	}
	c := Avg(colors...)
	c = Reflect(m.Color.Scale(m.Reflective), c)
	return c
}

func (sf *SceneFrame) Image() *image.RGBA {
	w, h := sf.Camera.Width, sf.Camera.Height
	z := (0.5 * float64(w)) / math.Tan(float64(sf.Camera.Angle)*0.5)
	t := sf.Camera.Q.TInv().T(d3.Translate(sf.Camera.Pt).TInv())
	img := image.NewRGBA(image.Rect(0, 0, sf.Camera.Width, sf.Camera.Height))
	ln := w * h
	//w64, h64 := float64(sf.Width), float64(sf.Height)
	dx, dy := w/2, h/2
	var idx32 int32 = -1
	wg := &sync.WaitGroup{}
	cpus := runtime.NumCPU()
	wg.Add(cpus)
	fn := func() {
		for {
			idx := int(atomic.AddInt32(&idx32, 1))
			if idx >= ln {
				break
			}
			ptX, ptY := (idx%w)-dx, (idx/w)-dy

			v := t.V(d3.V{float64(ptX), float64(ptY), z})
			l := line.Line{
				T0: sf.Camera.Pt,
				D:  v,
			}
			ix, iy := ptX+dx, h-ptY-dy-1
			c := sf.trace(l, sf.Depth)
			img.SetRGBA(ix, iy, color.RGBA{uint8(c.R * 255), uint8(c.G * 255), uint8(c.B * 255), 255})
		}
		wg.Add(-1)
	}
	for i := 0; i < cpus; i++ {
		go fn()
	}
	wg.Wait()
	return img
}

type Material struct {
	Color      *Color
	Luminous   float64
	Reflective float64
	Diffuse    angle.Rad
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
