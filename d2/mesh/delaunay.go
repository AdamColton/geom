package mesh

import (
	"github.com/adamcolton/geom/d2"
	"github.com/adamcolton/geom/d2/shape/ellipse"
	"github.com/adamcolton/geom/d2/shape/polygon"
	"github.com/adamcolton/geom/d2/shape/triangle"
)

type Delaunay []triangle.Triangle

func NewDelaunay(pts []d2.Pt) Delaunay {
	if len(pts) < 3 {
		return nil
	}

	baseTriangle := triangle.Super(pts)
	delaunay := Delaunay{baseTriangle}
	delaunay = delaunay.Add(pts)
	delaunay = delaunay.Remove(baseTriangle[:])
	return delaunay
}

func (d Delaunay) Add(pts []d2.Pt) Delaunay {
	for _, pt := range pts {
		d = d.add(pt)
	}
	return d
}

func (d Delaunay) add(pt d2.Pt) Delaunay {
	var badTriangles, goodTriangles []triangle.Triangle

	for _, t := range d {
		if ellipse.CircumscribeCircle(t).Contains(pt) {
			badTriangles = append(badTriangles, t)
		} else {
			goodTriangles = append(goodTriangles, t)
		}
	}

	es := NewEdgeMesh()
	for _, t := range badTriangles {
		es.Flip(t[:]...)
	}

	for e, v := range es {
		if !v {
			continue
		}
		goodTriangles = append(goodTriangles, triangle.Triangle{e[0], e[1], pt})
	}
	return goodTriangles
}

func (d Delaunay) Remove(pts []d2.Pt) Delaunay {
	var out Delaunay
outer:
	for _, t := range d {
		for _, p1 := range pts {
			for _, p2 := range t {
				if p1 == p2 {
					continue outer
				}
			}
		}
		out = append(out, t)
	}
	return out
}

func (d Delaunay) EdgeMesh() EdgeMesh {
	return EdgeMeshFromTrianges(d)
}

func (d Delaunay) Voronoi() EdgeMesh {
	m := NewEdgeMesh()
	dEdges := make(map[Edge]d2.Pt)
	for _, t := range d {
		c1 := t.CircumCenter()
		for _, e := range Edges(t[:]...) {
			if c2, ok := dEdges[e]; ok {
				m.Add(c1, c2)
				delete(dEdges, e)
			} else {
				dEdges[e] = c1
			}
		}
	}
	return m
}

func (d Delaunay) VoronoiPolygons(skip []d2.Pt) []polygon.Polygon {
	pps := make(map[d2.Pt]polygon.PolarPolygon)
	for _, t := range d {
		c := t.CircumCenter()
		for _, v := range t {
			pps[v] = append(pps[v], c.Subtract(v).Polar())
		}
	}

	var out []polygon.Polygon
outer:
	for c, pp := range pps {
		for _, s := range skip {
			if c == s {
				continue outer
			}
		}
		if len(pps) > 2 {
			out = append(out, pp.Polygon(c))
		}
	}
	return out
}
