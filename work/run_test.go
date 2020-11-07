package work

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testRecord struct {
	core, idx int
}

func TestRunRange(t *testing.T) {
	out := make([]*testRecord, 1000)
	RunRange(len(out), func(rangeIdx, coreIdx int) {
		assert.Nil(t, out[rangeIdx])
		out[rangeIdx] = &testRecord{
			core: coreIdx,
			idx:  rangeIdx,
		}
		time.Sleep(time.Millisecond)
	})

	for i, r := range out {
		assert.Equal(t, r.idx, i)
	}
}
