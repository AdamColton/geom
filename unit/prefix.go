package unit

import (
	"fmt"
	"math"
	"strings"
)

type Prefix float64

const (
	Tera  Prefix = 1e12
	Giga  Prefix = 1e9
	Mega  Prefix = 1e6
	Kilo  Prefix = 1e3
	Hecto Prefix = 1e2
	Deka  Prefix = 1e1
	Deci  Prefix = 1e-1
	Centi Prefix = 1e-2
	Mili  Prefix = 1e-3
	Micro Prefix = 1e-6
	Nano  Prefix = 1e-9
	Pico  Prefix = 1e-12
	Femto Prefix = 1e-15
)

var symbols = map[Prefix]string{
	Tera:  "T",
	Giga:  "G",
	Mega:  "M",
	Kilo:  "k",
	Hecto: "h",
	Deka:  "da",
	Deci:  "d",
	Centi: "c",
	Mili:  "m",
	Micro: "µ",
	Nano:  "n",
	Pico:  "p",
	Femto: "f",
}

func (p Prefix) String() string {
	m := p.Mag()
	s := symbols[m]
	mult := p / m
	if mult == 1.0 {
		return s
	}
	fs := fmt.Sprintf("%0.4f", mult)
	fs = strings.TrimRight(strings.TrimRight(fs, "0"), ".")
	return fmt.Sprintf("%s%s", fs, s)
}

func (p Prefix) Mag() Prefix {
	f := float64(p)
	mag := math.Pow(10, math.Floor(math.Log10(f)))
	if mag > float64(Tera) {
		mag = float64(Tera)
	} else if mag < float64(Femto) {
		mag = float64(Femto)
	} else {
		for {
			_, found := symbols[Prefix(mag)]
			if found {
				break
			}
			mag /= 10
		}
	}
	return Prefix(mag)
}
