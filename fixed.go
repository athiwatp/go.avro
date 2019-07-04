package avro

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// Fixed represents the "fixed" complex type.
type Fixed struct {
	NameFields
	Size uint `json:"size"`
}

func (f Fixed) Type() string { return "fixed" }

func (f Fixed) Valid() error {
	if err := f.NameFields.Valid(); err != nil {
		return err
	}
	if f.Size == 0 {
		return errors.New(`fixed size cannot be 0`)
	}
	return nil
}

func (f Fixed) Validate(v interface{}) error {
	if v == nil {
		return errors.New(`nil is not a valid fixed`)
	}
	if err := f.Valid(); err != nil {
		return fmt.Errorf(`validation aborted, fixed schema is invalid: %s`, err)
	}

	// Type switch for primitive types.
	switch s := v.(type) {
	case []byte:
		return f.checkBytes(s)
	case string:
		return f.checkBytes([]byte(s))
	case *[]byte:
		return f.checkBytes(*s)
	case *string:
		return f.checkBytes([]byte(*s))
	}

	// Check for concrete types with reflect.
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Slice:
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			return f.checkBytes(rv.Bytes())
		}
	case reflect.String:
		return f.checkBytes([]byte(rv.String()))
	}

	return fmt.Errorf(`value of type "%T" is not a valid fixed`, v)
}

func (f Fixed) checkBytes(b []byte) error {
	if b == nil {
		return errors.New(`nil is not a valid fixed`)
	}
	if len(b) != int(f.Size) {
		return fmt.Errorf(`value has %d bytes, but should have %d`, len(b), f.Size)
	}
	return nil
}

type jsonFixed struct {
	Type string `json:"type"`
	NameFields
	Size uint `json:"size"`
}

// UnmarshalJSON is implemented to check the "type" field.
func (f *Fixed) UnmarshalJSON(data []byte) error {
	raw := jsonFixed{}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return fmt.Errorf(`unmarshal fixed json: "%s"`, err)
	}
	if raw.Type != f.Type() {
		return fmt.Errorf(`cannot read type "%s" into %s`, raw.Type, f.Type())
	}
	f.NameFields = raw.NameFields
	f.Size = raw.Size
	return nil
}

// MarshalJSON adds the "type" field and validates before marshaling.
func (f Fixed) MarshalJSON() ([]byte, error) {
	if err := f.Valid(); err != nil {
		return nil, err
	}
	raw := jsonFixed{
		Type:       f.Type(),
		NameFields: f.NameFields,
		Size:       f.Size,
	}
	return json.Marshal(&raw)
}
