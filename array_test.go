package avro

import (
	"encoding/json"
	"testing"

	"github.com/matryer/is"
)

func TestArray_Type(t *testing.T) {
	is := is.New(t)

	a := Array{}
	is.Equal(a.Type(), "array") // array type is "array"
}

func TestArray_Valid(t *testing.T) {
	is := is.New(t)

	a := Array{}
	is.True(a.Valid() != nil) // invalid when Items is nil

	a.Items = Array{}
	is.True(a.Valid() != nil) // invalid when Items is invalid

	a.Items = String
	is.NoErr(a.Valid()) // valid when Items is valid Schema
}

func TestArray_Validate(t *testing.T) {
	is := is.New(t)

	a := Array{}
	is.True(a.Validate([]interface{}{}) != nil) // invalid schema cannot validate
	is.True(a.Validate([]string{}) != nil)      // invalid schema cannot validate
	is.True(a.Validate(nil) != nil)             // invalid schema cannot validate

	a.Items = mockValidNamedSchema
	is.True(a.Validate([]int{0, mockValue}) != nil)   // invalid if contains invalid items
	is.True(a.Validate([]int{mockValue, 0}) != nil)   // invalid if contains invalid items
	is.NoErr(a.Validate([]int{mockValue, mockValue})) // valid if contains only valid items

	val := []int{mockValue}
	is.NoErr(a.Validate(&val)) // pointer to valid value is valid
	ptr := &val
	is.True(a.Validate(&ptr) != nil) // double pointer to valid value is invalid

	// []interface{} static validation
	siVal := []interface{}{mockValue}
	is.NoErr(a.Validate(siVal))  // []interface{} with valid value is valid
	is.NoErr(a.Validate(&siVal)) // *[]interface{} with valid value is valid

	siVal = append(siVal, 0)
	is.True(a.Validate(siVal) != nil)  // []interface{} with invalid value is invalid
	is.True(a.Validate(&siVal) != nil) // []interface{} with invalid value is invalid
}

func TestArray_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		Name  string
		Spec  []byte
		Error bool
		Items Schema
	}{
		{"valid spec", []byte(`{"type":"array","items":"string"}`), false, String},
		{"empty json", []byte(``), true, nil},
		{"invalid json", []byte(`{`), true, nil},
		{"invalid type", []byte(`{"type":"__WRONG__","items":"string"}`), true, nil},
		{"invalid items", []byte(`{"type":"array","items":"__WRONG__"}`), true, nil},
	}

	for _, test := range tests {
		test := test // capture range variable
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			a := Array{}
			err := a.UnmarshalJSON(test.Spec)
			if test.Error {
				is.True(err != nil) // should return error
			} else {
				is.NoErr(err)                 // should not return error
				is.Equal(a.Items, test.Items) // unexpected items schema
			}
		})
	}
}

func TestArray_MarshalJSON(t *testing.T) {
	is := is.New(t)

	a := Array{}
	b, err := json.Marshal(a)
	is.True(err != nil) // invalid array should not marshal

	spec := []byte(`{"type":"array","items":"string"}`)
	b, err = json.Marshal(Array{String})
	is.NoErr(err)                     // marshal valid schema no error
	is.Equal(string(b), string(spec)) // marshal valid spec
}
