package main

import (
	"image/png"
	"os"

	"github.com/adamcolton/geom/angle"
	"github.com/nfnt/resize"

	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/render/material"
	"github.com/adamcolton/geom/d3/render/raytrace"
	"github.com/adamcolton/geom/d3/render/scene"
	"github.com/adamcolton/geom/d3/solid/mesh"
)

func main() {
	s := &raytrace.Frame{
		Frame: &scene.Frame{
			Camera: scene.NewCamera(d3.Pt{0, 0, 0}, angle.Deg(30)).
				SetRot(d3.Q{1, 0, 0, 0}).
				Widescreen(500),
			Meshes: make([]*scene.FrameMesh, 0, 3),
		},
		RayFrame: &raytrace.RayFrame{
			Background: backgroundShader{}.RayShader,
			Depth:      4,
			RayMult:    2,
			ImageScale: 5,
		},
	}
	s.AddMesh(getArrow(), d3.Identity(), arrowShader{})
	s.AddMesh(getLight(), d3.Identity(), lightShader{})
	s.AddMesh(getFloor(), d3.Identity(), floorShader{})
	s.PopulateShaders()
	img := s.Image(nil)

	f, _ := os.Create("test.png")
	png.Encode(f, resize.Resize(uint(s.Camera.Size.X), 0, img, resize.Bilinear))
	f.Close()
}

func getArrow() *mesh.TriangleMesh {
	f := []d3.Pt{
		{0, 2, -10},
		{1.5, 3.5, -10},
		{3, 2, -10},
		{2, 2, -10},
		{2, 0, -10},
		{1, 0, -10},
		{1, 2, -10},
	}
	f = d3.Translate(d3.V{-1.5, -1.0, 0}).T().T(
		d3.Rotation{
			Angle: angle.Deg(90),
			Plane: d3.XY,
		}.T(),
	).Pts(f)
	m, err := mesh.NewExtrusion(f).
		Extrude(d3.Translate(d3.V{0, 0, -1}).T()).
		Close().
		TriangleMesh()
	if err != nil {
		panic(err)
	}
	return &m
}

func getLight() *mesh.TriangleMesh {
	z := -1.0
	size := 0.3
	return &mesh.TriangleMesh{
		Pts: []d3.Pt{
			{-size, -size, z},
			{size, -size, z},
			{size, size, z},
			{-size, size, z},
		},
		Polygons: [][][3]uint32{
			{
				{0, 1, 2},
				{0, 2, 3},
			},
		},
	}
}

func getFloor() *mesh.TriangleMesh {
	y := -1.5
	size := 100.0
	return &mesh.TriangleMesh{
		Pts: []d3.Pt{
			{-size, y, -size},
			{size, y, -size},
			{size, y, size},
			{-size, y, size},
		},
		Polygons: [][][3]uint32{
			{
				{0, 1, 2},
				{0, 2, 3},
			},
		},
	}
}

type backgroundShader struct{}

func (backgroundShader) RayShader(ctx *raytrace.Context) *raytrace.Material {
	return &raytrace.Material{
		Color:    &material.Color{0.6, 0.6, 1.0},
		Luminous: 1.0,
		Diffuse:  angle.Deg(90),
	}
}

type arrowShader struct{}

func (arrowShader) RayShader(ctx *raytrace.Context) *raytrace.Material {
	return &raytrace.Material{
		Color:    &material.Color{0.8, 0.5, 0.5},
		Luminous: 0,
		Diffuse:  angle.Deg(2),
	}
}

var lightMaterial = &raytrace.Material{
	Color:    &material.Color{1.0, 1.0, 1.0},
	Luminous: 1.0,
	Diffuse:  angle.Deg(90),
}

type lightShader struct{}

func (lightShader) RayShader(ctx *raytrace.Context) *raytrace.Material {
	return lightMaterial
}

var (
	c1 = &material.Color{0.9, 0.9, 0.9}
	c2 = &material.Color{0.1, 0.1, 0.1}
)

type floorShader struct{}

func (floorShader) RayShader(ctx *raytrace.Context) *raytrace.Material {
	pt := ctx.Ray.Pt1(ctx.T)
	x, z := int(pt.X), int(pt.Z)
	c := c1
	if (x^z)&1 == 1 {
		c = c2
	}

	return &raytrace.Material{
		Color:   c,
		Diffuse: angle.Deg(45),
	}
}
