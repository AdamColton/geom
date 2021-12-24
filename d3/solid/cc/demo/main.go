package main

import (
	"os"
	"strconv"

	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/solid/cc"
	"github.com/adamcolton/geom/d3/solid/mesh"
)

func main() {
	f := []d3.Pt{
		{0, 0, 0},
		{1, 0, 0},
		{1, 1, 0},
		{0, 1, 0},
	}

	m := mesh.NewExtrusion(f).
		Extrude(
			d3.Translate(d3.V{0, 0, 1}).T(),
		).
		Close()

	file, _ := os.Create("cube_0.obj")
	m.WriteObj(file)
	file.Close()
	for i := 0; i < 3; i++ {
		m = cc.Subdivide(m, 1)
		file, _ := os.Create("cube_" + strconv.Itoa(i+1) + ".obj")
		m.WriteObj(file)
		file.Close()
	}

	f = []d3.Pt{
		{0, 0, 0},
		{0.5, -0.25, 0},
		{1, 0, 0},
		{1, 1, 0},
		{0.5, 1.25, 0},
		{0, 1, 0},
	}

	m = mesh.NewExtrusion(f).
		EdgeExtrude(d3.Scale(d3.V{2, 2, 1}).T()).
		Extrude(
			d3.Translate(d3.V{0, 0, 1}).T(),
			d3.Translate(d3.V{0, 0, 2}).T(),
			d3.Translate(d3.V{0, 0, 1}).T(),
		).
		EdgeMerge(d3.Scale(d3.V{0.5, 0.5, 1}).T()).
		Close()

	file, _ = os.Create("hex_0.obj")
	m.WriteObj(file)
	file.Close()
	for i := 0; i < 3; i++ {
		m = cc.Subdivide(m, 1)
		file, _ := os.Create("hex_" + strconv.Itoa(i+1) + ".obj")
		m.WriteObj(file)
		file.Close()
	}
}
