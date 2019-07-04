package avro

import (
	"encoding/json"
	"testing"

	"github.com/matryer/is"
)

func TestUnion_Type(t *testing.T) {
	is := is.New(t)

	u := Union{}
	is.Equal(u.Type(), "union") // union type is "union"
}

func TestUnion_Valid(t *testing.T) {
	is := is.New(t)

	u := Union{}
	is.True(u.Valid() != nil) // empty union is invalid

	u = Union{nil}
	is.True(u.Valid() != nil) // nil schema is invalid

	u = Union{mockValidNamedSchema, mockValidNamedSchema}
	is.True(u.Valid() != nil) // duplicate primitive types are invalid

	u = Union{mockInvalidNamedSchema}
	is.True(u.Valid() != nil) // invalid contained type is invalid

	u = Union{
		mockValidNamedSchema,
		buildMockValidPrimitiveSchema(mockValidNamedSchema.Name),
	}
	is.NoErr(u.Valid()) // named type name does not clash with primitive type

	u = Union{
		mockValidNamedSchema,
		buildMockValidNamedSchema("OtherMock", "mock-other"),
	}
	is.NoErr(u.Valid()) // differently typed named types with same name is valid

	u = Union{
		buildMockValidNamedSchema("SameName", "enum"),
		buildMockValidNamedSchema("SameName", "enum"),
	}
	is.True(u.Valid() != nil) // duplicate named type is invalid

	u = Union{Union{mockValidNamedSchema}}
	is.True(u.Valid() != nil) // directly contained union is invalid
}

func TestUnion_Validate(t *testing.T) {
	is := is.New(t)

	u := Union{}
	is.True(u.Validate(nil) != nil)       // empty union should not validate
	is.True(u.Validate(mockValue) != nil) // empty union should not validate

	u = Union{mockValidPrimitiveSchema}
	is.NoErr(u.Validate(mockValue)) // valid union validates value
	is.True(u.Validate(0) != nil)   // valid union invalidates value
}

func TestUnion_UnmarshalJSON(t *testing.T) {
	is := is.New(t)

	u := Union{}
	is.True(u.UnmarshalJSON([]byte(`[`)) != nil)             // invalid spec should error
	is.True(u.UnmarshalJSON([]byte(`["__WRONG__"]`)) != nil) // contained invalid spec should error

	spec := []byte(`["string", "int"]`)
	is.NoErr(json.Unmarshal(spec, &u)) // unmarshals without error
	is.Equal(len(u), 2)                // all items are unmarshalled
	is.Equal(u[0], String)             // unmarshals primitive
	is.Equal(u[1], Int)                // unmarshals primitive
}

func TestUnion_MarshalJSON(t *testing.T) {
	is := is.New(t)

	u := Union{String}
	b, err := u.MarshalJSON()
	is.NoErr(err)                     // marshals valid union without error
	is.Equal(b, []byte(`["string"]`)) // marshals into array with children

	u = Union{mockInvalidNamedSchema}
	_, err = u.MarshalJSON()
	is.True(err != nil) // does not marshal with invalid child schema
}
