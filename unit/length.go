package unit

import (
	"fmt"
	"strings"
)

type Length float64

const (
	Meter = 1.0
	Foot  = 0.3048
	Inch  = 0.0254
)

var lengthConversion = map[string]float64{
	"m":     Meter,
	"meter": Meter,
	"foot":  Foot,
	"ft":    Foot,
	"'":     Foot,
	"inch":  Inch,
	"in":    Inch,
	"\"":    Inch,
}

func NewLenUnit(unit string) (Length, error) {
	unit = strings.TrimSpace(unit)
	uln := len(unit)
	for s, c := range lengthConversion {
		sln := len(s)
		if sln <= uln && unit[uln-sln:] == s {
			p, err := Pre(unit[:uln-sln])
			if err != nil {
				return 1.0, err
			}
			return Length(float64(p) * c), nil
		}
	}
	return 1.0, fmt.Errorf("counld not parse")
}
