package avro

import (
	"fmt"
	"testing"

	"github.com/matryer/is"
)

func TestPrimitive_Type(t *testing.T) {
	is := is.NewRelaxed(t)

	for _, p := range primitives {
		t.Run(string(p), func(t *testing.T) {
			is.True(len(p) > 0)           // primitive is not empty string
			is.Equal(p.Type(), string(p)) // primitive returns itself as type
		})
	}
}

func TestPrimitive_Valid(t *testing.T) {
	is := is.NewRelaxed(t)

	for _, p := range primitives {
		t.Run(string(p), func(t *testing.T) {
			is.NoErr(p.Valid()) // built in primitive is valid
		})
	}

	is.True(Primitive("__WRONG__").Valid() != nil) // custom primitive is invalid
}

func TestPrimitive_Validate(t *testing.T) {
	tests := []struct {
		Primitive
		Value interface{}
		Error bool
	}{
		{Null, nil, false},
		{Null, "hello", true},
		{Null, 0, true},
		{Boolean, true, false},
		{Boolean, false, false},
		{Boolean, 0, true},
		{Boolean, nil, true},
		{Int, 0, false},
		{Int, 1, false},
		{Int, 1.5235, true},
		{Int, nil, true},
		{Long, int64(0), false},
		{Long, int64(1), false},
		{Long, int(1), uint64(^uint(0)) != ^uint64(0)}, // Only error on non-64-bit
		{Long, float32(1.5), true},
		{Long, uint8(0x0), true},
		{Long, nil, true},
		{Float, float32(0.0), false},
		{Float, float32(1.5), false},
		{Float, float64(1.5), true},
		{Float, int(0), true},
		{Float, uint8(0x0), true},
		{Float, "hello", true},
		{Float, 0, true},
		{Float, nil, true},
		{Double, float64(0.0), false},
		{Double, float64(1.5), false},
		{Double, 0, true},
		{Double, 0x0, true},
		{Double, "hello", true},
		{Double, 0, true},
		{Double, nil, true},
		{Bytes, []byte{0x0}, false},
		{Bytes, []byte{}, false},
		{Bytes, "hello", true},
		{Bytes, 0, true},
		{Bytes, nil, true},
		{String, "", false},
		{String, "hello", false},
		{String, 0, true},
		{String, nil, true},
	}

	for _, test := range tests {
		test := test
		var name string
		if test.Error {
			name = fmt.Sprintf(`"%s" invalidates %v`, test.Type(), test.Value)
		} else {
			name = fmt.Sprintf(`"%s" validates %v`, test.Type(), test.Value)
		}
		t.Run(name, func(t *testing.T) {
			is := is.New(t)
			if test.Error {
				is.True(test.Validate(test.Value) != nil) // primitive invalidates value
			} else {
				t.Log(fmt.Sprintf("%T\n", test.Value))
				is.NoErr(test.Validate(test.Value))  // primitive validates value
				is.NoErr(test.Validate(&test.Value)) // primitive validates *interface{} with value
			}
		})
	}
}
