package scene

import (
	"github.com/adamcolton/geom/d3"
	"github.com/adamcolton/geom/d3/solid/mesh"
)

type TriangleIndex struct {
	MeshIdx     int
	PolygonIdx  int
	TriangleIdx int
}

type Mesh struct {
	Original *mesh.TriangleMesh
	TransformFactory
	Shader interface{}
}

type FrameMesh struct {
	Original *mesh.TriangleMesh
	Shader   interface{}
	Space    []d3.Pt
}
