package avro

import (
	"encoding/json"
	"testing"

	"github.com/matryer/is"
)

func TestEnum_Type(t *testing.T) {
	is := is.New(t)

	e := Enum{}
	is.Equal(e.Type(), "enum") // enum type is "enum"
}

func TestEnum_Valid(t *testing.T) {
	is := is.New(t)

	e := Enum{}
	is.True(e.Valid() != nil) // nameless enum should be invalid

	e.Name = invalidName
	is.True(e.Valid() != nil) // invalid name should be invalid

	e.Name = "Test"
	is.NoErr(e.Valid()) // name-only enum should be valid

	e.Symbols = []string{"one", "two"}
	is.NoErr(e.Valid()) // valid name and symbols should be valid

	e.Symbols = []string{"two", "two"}
	is.True(e.Valid() != nil) // duplicate symbols should be invalid

	e.Symbols = []string{"one", invalidName}
	is.True(e.Valid() != nil) // invalid symbol name should be invalid
}

func TestEnum_Validate(t *testing.T) {
	is := is.New(t)

	e := Enum{
		Symbols: []string{"x"},
	}
	is.True(e.Validate("x") != nil) // invalid schema cannot validate

	e = Enum{
		NameFields: NameFields{
			Name: "Test",
		},
		Symbols: []string{"x", "y"},
	}
	is.NoErr(e.Validate("x"))       // should validate matching symbol
	is.NoErr(e.Validate("y"))       // should validate matching symbol
	is.True(e.Validate("z") != nil) // should reject invalid string symbol
	is.True(e.Validate(0) != nil)   // should reject unsupported type
	is.True(e.Validate(nil) != nil) // should reject nil
	v := "y"
	is.NoErr(e.Validate(&v)) // should validate pointer to valid symbol

	type custom string
	cv := custom("x")
	is.NoErr(e.Validate(cv))  // should validate when concrete type is string
	is.NoErr(e.Validate(&cv)) // should validate when concrete type is string
}

func TestEnum_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		Name    string
		Spec    []byte
		Error   bool
		Symbols []string
	}{
		{
			"valid spec",
			[]byte(`{"type":"enum","name":"Test","symbols":["a","b"]}`),
			false,
			[]string{"a", "b"},
		},
		{"empty json", []byte(``), true, nil},
		{"invalid json", []byte(`{`), true, nil},
		{
			"invalid type",
			[]byte(`{"type":"__WRONG__","name":"Test","symbols":["a","b"]}`),
			true,
			[]string{"a", "b"},
		},
		{
			"invalid symbols",
			[]byte(`{"type":"enum","name":"Test","symbols":"__WRONG__"}`),
			true,
			nil,
		},
	}

	for _, test := range tests {
		test := test // capture range variable
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			e := Enum{}
			err := json.Unmarshal(test.Spec, &e)
			if test.Error {
				is.True(err != nil) // should return error
			} else {
				is.NoErr(err)                               // should not return error
				is.Equal(len(e.Symbols), len(test.Symbols)) // symbols should be same length
				for i, sym := range e.Symbols {
					is.Equal(sym, test.Symbols[i]) // symbol should match expected value
				}
			}
		})
	}
}

func TestEnum_MarshalJSON(t *testing.T) {
	is := is.New(t)

	e := Enum{}
	data, err := json.Marshal(e)
	is.True(err != nil) // invalid enum should not marshal

	spec := []byte(`{"type":"enum","name":"Test","symbols":["a","b"]}`)
	e.Name = "Test"
	e.Symbols = []string{"a", "b"}
	data, err = json.Marshal(e)
	is.NoErr(err)                        // valid enum should marshal
	is.Equal(string(data), string(spec)) // spec should match expected
}
