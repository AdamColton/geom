package scene

import (
	"runtime"
	"sync"
	"sync/atomic"
)

var cpus = runtime.NumCPU()

// Run takes a worker that will operate on some workload until complete and
// runs it on each core.
func Run(worker func(coreIdx int)) {
	wg := &sync.WaitGroup{}
	wg.Add(cpus)
	wrap := func(coreIdx int) {
		worker(coreIdx)
		wg.Add(-1)
	}
	for i := 0; i < cpus; i++ {
		go wrap(i)
	}
	wg.Wait()
}

// RunRange will call worker in parallel with every int from 0 to max (0
// inclusive, max exclusive).
func RunRange(max int, worker func(rangeIdx, coreIdx int)) {
	var idx32 int32 = -1
	Run(func(coreIdx int) {
		for {
			idx := int(atomic.AddInt32(&idx32, 1))
			if idx >= max {
				return
			}
			worker(idx, coreIdx)
		}
	})
}
