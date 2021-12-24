package main

import (
	"github.com/adamcolton/geom/d2/grid"
	"github.com/adamcolton/geom/d2/grid/drawgrid"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	d := grid.Diamond(grid.Pt{5, 5}, 8, 0.5)
	img := drawgrid.Image(d, drawgrid.FloatToGrey)
	drawgrid.SavePNG(img, "diamond.png")
}
