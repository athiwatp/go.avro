package avro

import (
	"encoding/json"
	"testing"

	"github.com/matryer/is"
)

func TestMap_Type(t *testing.T) {
	is := is.New(t)

	m := Map{}
	is.Equal(m.Type(), "map") // map type is "map"
}

func TestMap_Valid(t *testing.T) {
	is := is.New(t)

	m := Map{}
	is.True(m.Valid() != nil) // invalid when Values is nil

	m.Values = Map{}
	is.True(m.Valid() != nil) // invalid when Values is invalid

	m.Values = String
	is.NoErr(m.Valid()) // valid when Values is valid Schema
}

func TestMap_Validate(t *testing.T) {
	is := is.New(t)

	m := Map{}
	is.True(m.Validate(map[string]interface{}{}) != nil) // invalid schema cannot validate
	is.True(m.Validate([]interface{}{}) != nil)          // invalid schema cannot validate
	is.True(m.Validate(nil) != nil)                      // invalid schema cannot validate

	m.Values = mockValidNamedSchema
	is.True(m.Validate(map[string]int{"one": 0, "two": mockValue}) != nil)   // invalid if contains invalid items
	is.NoErr(m.Validate(map[string]int{"one": mockValue, "two": mockValue})) // valid if contains only valid items

	val := map[string]int{"x": mockValue}
	is.NoErr(m.Validate(&val)) // pointer to valid value is valid
	ptr := &val
	is.True(m.Validate(&ptr) != nil) // double pointer to valid value is invalid

	valBadKey := map[int]int{mockValue: mockValue}
	is.True(m.Validate(valBadKey) != nil)  // non-string key is invalid
	is.True(m.Validate(&valBadKey) != nil) // non-string key is invalid

	// static map[string]interface{} validation
	msi := map[string]interface{}{"x": mockValue}
	is.NoErr(m.Validate(msi))  // valid msi is valid
	is.NoErr(m.Validate(&msi)) // pointer to valid msi is valid

	msi["y"] = "invalid"
	is.True(m.Validate(msi) != nil)  // invalid msi is invalid
	is.True(m.Validate(&msi) != nil) // pointer to invalid msi is invalid

	msi = nil
	is.True(m.Validate(msi) != nil)  // nil msi is invalid
	is.True(m.Validate(&msi) != nil) // pointer to nil msi is invalid
}

func TestMap_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		Name   string
		Spec   []byte
		Error  bool
		Values Schema
	}{
		{"valid spec", []byte(`{"type":"map","values":"string"}`), false, String},
		{"empty json", []byte(``), true, nil},
		{"invalid json", []byte(`{`), true, nil},
		{"invalid type", []byte(`{"type":"__WRONG__","values":"string"}`), true, nil},
		{"invalid items", []byte(`{"type":"map","values":"__WRONG__"}`), true, nil},
	}

	for _, test := range tests {
		test := test // capture range variable
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			m := Map{}
			err := m.UnmarshalJSON(test.Spec)
			if test.Error {
				is.True(err != nil) // should return error
			} else {
				is.NoErr(err)                   // should not return error
				is.Equal(m.Values, test.Values) // unexpected values schema
			}
		})
	}
}

func TestMap_MarshalJSON(t *testing.T) {
	is := is.New(t)

	m := Map{}
	b, err := json.Marshal(m)
	is.True(err != nil) // invalid map should not marshal

	spec := []byte(`{"type":"map","values":"string"}`)
	b, err = json.Marshal(Map{String})
	is.NoErr(err)                     // marshal valid schema no error
	is.Equal(string(b), string(spec)) // marshal valid spec
}
