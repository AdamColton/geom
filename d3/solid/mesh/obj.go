package mesh

import (
	"bufio"
	"github.com/adamcolton/geom/d3"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// https://en.wikipedia.org/wiki/Wavefront_.obj_file

var hdr_v = []byte("v ")
var hdr_f = []byte("f")
var spc = []byte(" ")
var nl = []byte("\n")

var Prec = 10

type ww struct {
	w   io.Writer
	err error
}

func (w ww) write(b []byte) {
	if w.err != nil {
		return
	}
	_, w.err = w.w.Write(b)
}

func (w ww) writeStr(str string) {
	if w.err != nil {
		return
	}
	_, w.err = w.w.Write([]byte(str))
}

func (m *Mesh) WriteObj(writer io.Writer) error {
	w := ww{w: writer}
	for _, pt := range m.Pts {
		w.write(hdr_v)
		w.writeStr(strconv.FormatFloat(pt.X, 'g', Prec, 64))
		w.write(spc)
		w.writeStr(strconv.FormatFloat(pt.Y, 'g', Prec, 64))
		w.write(spc)
		w.writeStr(strconv.FormatFloat(pt.Z, 'g', Prec, 64))
		w.write(nl)
	}

	for _, t := range m.Polygons {
		w.write(hdr_f)
		for _, idx := range t {
			w.write(spc)
			w.writeStr(strconv.Itoa(int(idx + 1)))
		}
		w.write(nl)
	}

	return w.err
}

var vRe = regexp.MustCompile(`v (\d+(?:\.\d+)? ?) (\d+(?:\.\d+)?) (\d+(?:\.\d+)?)`)
var fRe = regexp.MustCompile(`f (\d+ ?)+`)

func ReadObj(reader io.Reader) (*Mesh, error) {
	mesh := &Mesh{}
	r := bufio.NewReader(reader)
	for line, err := r.ReadString('\n'); err == nil; line, err = r.ReadString('\n') {
		if m := vRe.FindStringSubmatch(line); len(m) > 0 {
			var pt d3.Pt
			pt.X, _ = strconv.ParseFloat(m[1], 64)
			pt.Y, _ = strconv.ParseFloat(m[2], 64)
			pt.Z, _ = strconv.ParseFloat(m[3], 64)
			mesh.Pts = append(mesh.Pts, pt)
		} else if m := fRe.FindStringSubmatch(line); len(m) > 0 {
			intStrs := strings.Split(m[1], " ")
			poly := make([]uint32, 0, len(intStrs))
			for _, s := range intStrs {
				if i, b := toUint32(s); b {
					poly = append(poly, i)
				}
			}
			mesh.Polygons = append(mesh.Polygons, poly)
		}
	}

	return mesh, nil
}

func toUint32(str string) (uint32, bool) {
	i, err := strconv.Atoi(str)
	return uint32(i), err == nil
}
