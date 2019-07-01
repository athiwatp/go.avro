package avro

import (
	"fmt"
	"reflect"
)

// ErrMissingRequiredAttribute is returned by Encode or Decode when a schema is
// missing an attribute required by the specification.
type ErrMissingRequiredAttribute struct {
	Attribute string
}

func (e ErrMissingRequiredAttribute) Error() string {
	return fmt.Sprintf(`missing required attribute "%s"`, e.Attribute)
}

// ErrInvalidAttributeType is returned by Decode when a schema has an attribute
// whose type is different than the type required by the specification.
type ErrInvalidAttributeType struct {
	Attribute string
	Expected  string
	Actual    string
}

func (e ErrInvalidAttributeType) Error() string {
	return fmt.Sprintf(
		`expected attribute "%s" to have type "%s" but it was "%s"`,
		e.Attribute, e.Expected, e.Actual,
	)
}

func invalidAttributeType(attr, expected string, actual interface{}) ErrInvalidAttributeType {
	return ErrInvalidAttributeType{
		Attribute: attr,
		Expected:  expected,
		Actual:    reflect.TypeOf(actual).Name(),
	}
}

// ErrInvalidValue is returned by Encode or Decode when a schema has an
// attribute whose value is not permitted by the specification.
type ErrInvalidValue struct {
	Field  string
	Actual string // Caller must convert actual to human readable string
}

func (e ErrInvalidValue) Error() string {
	return fmt.Sprintf(`"%s" is an invalid value for field "%s"`, e.Actual, e.Field)
}

type ErrValidation struct {
	FieldErrors map[string]ErrValidation
}

func (e ErrValidation) Error() string { return "" }
