package avro

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"
)

var nameRegex = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// Schema represents the functionality of all Avro Schema types.
type Schema interface {
	// Type returns the base type's name ("string", "union", "record", etc.).
	Type() string

	// Valid checks if the Schema itself is valid.
	// If there is an error, it will have type ErrValidation.
	Valid() error

	// Validate checks if value conforms to this Schema.
	// If value is a pointer it will be dereferenced once then checked.
	// If there is an error, it will have type ErrValidation.
	Validate(value interface{}) error
}

// NameFields embeds data unique to named Schema types (Record, Enum, Fixed).
type NameFields struct {
	Name      string   `json:"name"`
	Namespace string   `json:"namespace,omitempty"`
	Aliases   []string `json:"aliases,omitempty"`
}

// NamedSchema represents the functionality of any Record, Enum, or Fixed Schema.
type NamedSchema interface {
	GetNameFields() NameFields
	Fullname() string
	Valid() error
}

// Fullname returns the full namespaced name of the Schema.
func (n NameFields) Fullname() string {
	if n.Namespace != "" {
		return fmt.Sprintf("%s.%s", n.Namespace, n.Name)
	}
	return n.Name
}

func (n NameFields) GetNameFields() NameFields {
	return n
}

func (n NameFields) Valid() error {
	if n.Name == "" {
		return errors.New(`name cannot be empty`)
	}
	if !nameRegex.MatchString(n.Name) {
		return fmt.Errorf(`"%s" is an invalid name`, n.Name)
	}
	return nil
}

type Factories map[Reference]func() interface{} // TODO: maybe Schema as key?

var DefaultFactories = Factories{
	"date": func() interface{} { return new(time.Time) },
}

// SchemaUnmarshalJSON creates a Schema from an Avro schema declaration.
func SchemaUnmarshalJSON(spec []byte) (Schema, error) {
	var i interface{}
	if err := json.Unmarshal(spec, &i); err != nil {
		return nil, fmt.Errorf("unmarshal schema json: %s", err)
	}
	switch s := i.(type) {
	case string:
		switch s {
		case "null", "boolean", "int", "long", "float", "double", "bytes", "string":
			return Primitive(s), nil
		}
		// TODO assume it's a named reference? Use a Registry?
		return nil, fmt.Errorf(`unsupported type: "%s"`, s)
	case []interface{}:
		return unmarshalAndValidate(spec, Union(nil))
	case map[string]interface{}:
		// Decode based on "type" field.
		t, ok := s["type"].(string)
		if !ok {
			return nil, invalidAttributeType("type", "string", s["type"])
		}
		if t == "" {
			return nil, ErrMissingRequiredAttribute{"type"}
		}
		switch t {
		case "null", "boolean", "int", "long", "float", "double", "bytes", "string":
			return Primitive(t), nil
		case "record":
			return unmarshalAndValidate(spec, Record{})
		case "enum":
			return unmarshalAndValidate(spec, Enum{})
		case "array":
			return unmarshalAndValidate(spec, Array{})
		case "map":
			return unmarshalAndValidate(spec, Map{})
		case "fixed":
			return unmarshalAndValidate(spec, Fixed{})
		}
		return nil, ErrInvalidValue{"type", t}
	}
	return nil, errors.New("the provided avro spec was not valid json")
}

func unmarshalAndValidate(spec []byte, s Schema) (Schema, error) {
	if err := json.Unmarshal(spec, s); err != nil {
		return nil, fmt.Errorf("unmarshal %s json: %s", s.Type(), err)
	}
	return s, s.Valid()
}
