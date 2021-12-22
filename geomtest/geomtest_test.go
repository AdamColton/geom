package geomtest

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	tt := map[string]struct {
		msg      []interface{}
		expected string
	}{
		"empty": {},
		"string": {
			msg:      []interface{}{"just a string"},
			expected: "just a string",
		},
		"int": {
			msg:      []interface{}{31415},
			expected: "31415",
		},
		"format-string": {
			msg:      []interface{}{"%0.2f-ish", 3.141592653},
			expected: "3.14-ish",
		},
		"many-args": {
			msg:      []interface{}{314, "test", false},
			expected: "314testfalse",
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc.expected, Message(tc.msg...))
		})
	}
}

type mockT struct {
	buf *bytes.Buffer
}

func newMock() *mockT {
	return &mockT{
		buf: bytes.NewBuffer(nil),
	}
}

func (m *mockT) Errorf(format string, args ...interface{}) {
	fmt.Fprintf(m.buf, format, args...)
}

func TestWithMessage(t *testing.T) {
	m := newMock()
	b := Equal(m, 1, 1.0, "Type-Mismatch")
	assert.False(t, b)
	assert.Equal(t, "unsupported_type: int: Type-Mismatch", m.buf.String())
}
