package avro

import (
	"testing"

	"github.com/matryer/is"
)

func TestRecord_Type(t *testing.T) {
	is := is.New(t)

	r := Record{}
	is.Equal(r.Type(), "record") // record type is "record"
}

func TestRecord_Valid(t *testing.T) {
	is := is.New(t)

	r := Record{}
	is.True(r.Valid() != nil) // nameless record should be invalid

	r.Name = "Test"
	is.NoErr(r.Valid()) // named record with no fields should be valid

	r.Fields = []Field{
		{Name: ""},
	}
	is.True(r.Valid() != nil) // having a nameless field should be invalid

	r.Fields[0].Name = "Test"
	is.True(r.Valid() != nil) // having a typeless field should be invalid

	r.Fields[0].Type = mockValidNamedSchema
	is.NoErr(r.Valid()) // having a named, typed field should be valid

	invalidDefault := interface{}(0)
	r.Fields[0].Default = &invalidDefault
	is.True(r.Valid() != nil) // field with invalid default should be invalid

	r.Fields[0].Default = nil
	r.Fields[0].Order = "__WRONG__"
	is.True(r.Valid() != nil) // field with invalid order should be invalid

	r.Fields[0].Order = "ascending"
	is.NoErr(r.Valid()) // field with valid order should be valid

	r.Fields[0].Type = mockInvalidNamedSchema
	is.True(r.Valid() != nil) // having a field with invalid schema should be invalid
}

func TestRecord_Validate(t *testing.T) {
	is := is.New(t)

	r := Record{}
	is.True(r.Validate(map[string]interface{}{}) != nil) // invalid schema cannot validate

	r = Record{
		NameFields: NameFields{
			Name: "Test",
		},
		Fields: []Field{
			{Name: "X", Type: mockValidNamedSchema},
		},
	}

	type custom struct {
		X int
	}
	is.NoErr(r.Validate(custom{mockValue}))  // valid struct
	is.NoErr(r.Validate(&custom{mockValue})) // valid pointer to struct
	is.True(r.Validate(custom{0}) != nil)    // invalid struct
	is.True(r.Validate(nil) != nil)          // nil should be invalid

	type customTag struct {
		Y int `avro:"X"`
	}
	is.NoErr(r.Validate(customTag{mockValue})) // should get name from struct tag

	type invalid struct {
		Z int
	}
	is.True(r.Validate(invalid{mockValue}) != nil) // missing field should be invalid

	validMSI := map[string]interface{}{"X": mockValue}
	is.NoErr(r.Validate(validMSI))                       // valid msi
	is.NoErr(r.Validate(&validMSI))                      // pointer to valid msi
	is.NoErr(r.Validate(map[string]int{"X": mockValue})) // valid custom map
	is.True(r.Validate(map[string]int{"X": 0}) != nil)   // custom map with invalid field value

	is.True(r.Validate(map[string]int{
		"X": mockValue,
		"Y": mockValue,
	}) != nil) // custom map with non-existing field

	invalidMSI := map[string]interface{}{"__WRONG__": mockValue}
	is.True(r.Validate(invalidMSI) != nil) // invalid field name should be invalid

	invalidMSI = map[string]interface{}{"X": 0}
	is.True(r.Validate(invalidMSI) != nil) // invalid field value should be invalid

	nilMSI := map[string]interface{}(nil)
	is.True(r.Validate(nilMSI) != nil)  // nil msi should be invalid
	is.True(r.Validate(&nilMSI) != nil) // nil msi pointer should be invalid

	badMap := map[int]interface{}{}
	is.True(r.Validate(badMap) != nil) // map without string key should be invalid

	is.True(r.Validate(0) != nil) // invalid type should be invalid
}
