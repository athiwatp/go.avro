package avro

import (
	"errors"
	"testing"

	"github.com/matryer/is"
)

var invalidName = "0_INVALID_NAME"

func TestNameRegex(t *testing.T) {
	is := is.New(t)
	is.True(!nameRegex.MatchString(invalidName)) // invalidName does not match nameRegex
}

type mockPrimitiveSchema struct {
	typeFunc     func() string
	validFunc    func() error
	validateFunc func(interface{}) error
}

func (m mockPrimitiveSchema) Type() string {
	return m.typeFunc()
}

func (m mockPrimitiveSchema) Valid() error {
	return m.validFunc()
}

func (m mockPrimitiveSchema) Validate(v interface{}) error {
	return m.validateFunc(v)
}

type mockNamedSchema struct {
	NameFields
	mockPrimitiveSchema
}

func (m mockNamedSchema) Type() string {
	return m.mockPrimitiveSchema.typeFunc()
}

func (m mockNamedSchema) Valid() error {
	return m.mockPrimitiveSchema.validFunc()
}

func (m mockNamedSchema) Validate(v interface{}) error {
	return m.mockPrimitiveSchema.validateFunc(v)
}

const mockValue = 0xDEADBEEF

var errMockInvalid = errors.New("value is not 0xDEADBEEF")

func buildMockValidPrimitiveSchema(schemaType string) mockPrimitiveSchema {
	return mockPrimitiveSchema{
		typeFunc: func() string {
			return schemaType
		},
		validFunc: func() error {
			return nil
		},
		validateFunc: func(v interface{}) error {
			if i, ok := v.(int); ok && i == mockValue {
				return nil
			}
			return errMockInvalid
		},
	}
}

func buildMockValidNamedSchema(name, schemaType string) mockNamedSchema {
	return mockNamedSchema{
		NameFields: NameFields{
			Name: name,
		},
		mockPrimitiveSchema: buildMockValidPrimitiveSchema(schemaType),
	}
}

// mockValidPrimitiveSchema mocks a valid primitive schema that checks if value == mockValue.
var mockValidPrimitiveSchema = buildMockValidPrimitiveSchema("mock-valid")

// mockValidNamedSchema mocks a valid named schema that checks if value == mockValue.
var mockValidNamedSchema = buildMockValidNamedSchema("MockValid", "mock-named-valid")

// mockInvalidNamedSchema mocks an invalid named schema that cannot validate any values.
var mockInvalidNamedSchema = mockNamedSchema{
	NameFields: NameFields{
		Name: "MockInvalid",
	},
	mockPrimitiveSchema: mockPrimitiveSchema{
		typeFunc: func() string {
			return "mock-named-invalid"
		},
		validFunc: func() error {
			return errors.New("mock-invalid is always invalid")
		},
		validateFunc: func(v interface{}) error {
			return errors.New("mock-invalid cannot validate any value")
		},
	},
}

func TestMock(t *testing.T) {
	is := is.New(t)
	is.True(mockValidNamedSchema.Valid() == nil)                  // mock schema is always valid
	is.True(mockValidNamedSchema.Validate(nil) == errMockInvalid) // non-mockValue is invalid
	is.True(mockValidNamedSchema.Validate(mockValue) == nil)      // mockValue is valid
}

func TestSchemaUnmarshalJSON(t *testing.T) {
	is := is.NewRelaxed(t)

	t.Run("json string to primitive", func(t *testing.T) {
		t.Parallel()

		s, err := SchemaUnmarshalJSON([]byte(`"int"`))
		is.NoErr(err)     // unmarshal valid json string without error
		is.True(s == Int) // returns correct primitive schema

		s, err = SchemaUnmarshalJSON([]byte(`"string"`))
		is.NoErr(err)        // unmarshal valid json string without error
		is.True(s == String) // returns correct primitive schema

		s, err = SchemaUnmarshalJSON([]byte(`""`))
		is.True(err != nil) // unmarshal invalid json string returns error
		is.True(s == nil)   // error returns nil schema
	})

	t.Run("json array to union", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("json object to primitive", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("json object to complex", func(t *testing.T) {
		t.Parallel()
	})
}
