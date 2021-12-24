package geomtest

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/adamcolton/geom/calc/cmpr"
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
	buf    *bytes.Buffer
	helper bool
}

func newMock() *mockT {
	return &mockT{
		buf: bytes.NewBuffer(nil),
	}
}

func (m *mockT) Errorf(format string, args ...interface{}) {
	fmt.Fprintf(m.buf, format, args...)
}

func (m *mockT) Helper() {
	m.helper = true
}

func TestWithMessage(tt *testing.T) {
	t := assert.New(tt)
	m := newMock()
	b := Equal(m, 1, 1.0, "Type-Mismatch")
	t.False(b)
	t.True(m.helper)
	t.Equal("unsupported_type: int: Type-Mismatch", m.buf.String())
}

func TestEqualInDelta(tt *testing.T) {
	t := assert.New(tt)

	m := newMock()
	EqualInDelta(m, 1.0, 1.0, Small)
	t.Equal("", m.buf.String())

	EqualInDelta(m, 1.0, 2.0, Small)
	t.Equal("Expected 1 got 2", m.buf.String())
}

type mockAssert int

func (m *mockAssert) AssertEqual(to interface{}, t cmpr.Tolerance) error {
	if *m == *(to.(*mockAssert)) {
		return nil
	}
	return fmt.Errorf("not equal")
}

func newMockAssert(i int) *mockAssert {
	m := mockAssert(i)
	return &m
}

func TestAssertEqual(tt *testing.T) {
	t := assert.New(tt)

	err := AssertEqual([]float64{1.0, 2.0}, []float64{1.0, 2.0}, Small)
	t.NoError(err)

	err = AssertEqual([]float64{1.0, 2.0}, []interface{}{1.0, 2.0}, Small)
	t.NoError(err)

	err = AssertEqual(newMockAssert(1), newMockAssert(1), Small)
	t.NoError(err)

	err = AssertEqual(newMockAssert(1), newMockAssert(2), Small)
	t.Error(err)

	err = AssertEqual(mockAssert(1), mockAssert(2), Small)
	t.Equal("unsupported_type: geomtest.mockAssert (*geomtest.mockAssert fulfills AssertEqualizer)", err.Error())

	err = AssertEqual([]float64{1.0, 2.0}, 1.0, Small)
	t.Equal("Types do not match: expected \"[]float64\", got \"float64\"", err.Error())

	err = AssertEqual([]float64{1.0, 2.0}, []interface{}{1.0, 2}, Small)
	t.Equal("\t1: Types do not match: expected \"float64\", got \"int\"", err.Error())
}

func TestGeomAssert(tt *testing.T) {
	t := assert.New(tt)

	m := newMock()
	g := New(m)

	b := g.Equal(1, 2)
	t.False(b)

	b = g.Equal(1.0, 1.0+(float64(Small)/10.0))
	t.True(b)

	b = g.Equal(newMockAssert(1), newMockAssert(2))
	t.False(b)
}
