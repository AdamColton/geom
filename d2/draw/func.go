package draw

import (
	"image/color"
	"reflect"
	"runtime"
	"strings"
	"sync"
)

// Call is a helper function that uses relection for generating images. It needs
// a Context generator. For each function passed in, a go routine starts that
// generates a context and passes into the function, then the image is saved
// with the name of the function.
func Call(gen func() *Context, funcs ...func(*Context)) {
	var wg sync.WaitGroup
	wrap := func(fn func(*Context)) {
		name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
		if trim := strings.LastIndex(name, "."); trim > -1 {
			name = name[trim+1:]
		}
		ctx := gen()
		fn(ctx)
		ctx.SavePNG(name + ".png")
		wg.Add(-1)
	}
	wg.Add(len(funcs))
	for _, fn := range funcs {
		go wrap(fn)
	}
	wg.Wait()
}

const (
	max16     = ^uint16(0)
	max16As64 = float64(max16)
)

// Color is a helper for defining colors using percentages.
func Color(r, g, b float64) color.Color {
	r16 := uint16(r * max16As64)
	g16 := uint16(g * max16As64)
	b16 := uint16(b * max16As64)
	return color.RGBA64{r16, g16, b16, max16}
}
