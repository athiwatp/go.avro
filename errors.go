package avro

import (
	"fmt"
	"reflect"
	"strings"
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
	error
	Children map[string]error
}

func (e ErrValidation) Error() string {
	return e.errIndent(0)
}

func (e ErrValidation) errIndent(idt int) string {
	errs := make([]string, 0, len(e.Children)+1)
	pad := strings.Repeat("  ", idt)
	if e.error != nil {
		errs = append(errs, pad+e.error.Error())
	}
	for key, err := range e.Children {
		var errStr string
		if err, ok := err.(ErrValidation); ok {
			errStr = err.errIndent(idt + 1)
		} else {
			errStr = err.Error()
		}
		errs = append(errs, fmt.Sprintf("%s%s: %s", pad, key, errStr))
	}
	return strings.Join(errs, "\n")
}
