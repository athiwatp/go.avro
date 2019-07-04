package avro

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// Array represents the "array" complex type.
type Array struct {
	Items Schema `json:"items"`
}

// Type returns the Avro type name "array".
func (a Array) Type() string { return "array" }

// Valid checks the array is valid.
// The Items schema must be present and valid.
func (a Array) Valid() error {
	if a.Items == nil {
		return ErrValidation{
			Children: map[string]error{
				"items": errors.New("cannot be nil"),
			},
		}
	}
	if err := a.Items.Valid(); err != nil {
		return ErrValidation{
			Children: map[string]error{
				"items": err,
			},
		}
	}
	return nil
}

// Validate checks that a Go value is an array or slice whose values all
// conform to the Items schema. If the value is a pointer, it will be
// dereferenced once before checking against the schema.
func (a Array) Validate(v interface{}) error {
	if v == nil {
		return errors.New(`nil is not a valid array`)
	}
	if err := a.Valid(); err != nil {
		return fmt.Errorf(`validation aborted, array schema is invalid: %s`, err)
	}

	// Statically check for []interface{}.
	if s, ok := v.([]interface{}); ok {
		errs := map[string]error{}
		for i, sv := range s {
			if err := a.Items.Validate(sv); err != nil {
				errs[fmt.Sprintf("item at index %d", i)] = err
			}
		}
		if len(errs) > 0 {
			return ErrValidation{
				Children: errs,
			}
		}
		return nil
	}

	// Use reflect to inspect all other types.
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		errs := map[string]error{}
		for i, l := 0, rv.Len(); i < l; i++ {
			item := rv.Index(i)
			if err := a.Items.Validate(item.Interface()); err != nil {
				errs[fmt.Sprintf("item at index %d", l)] = err
			}
		}
		if len(errs) > 0 {
			return ErrValidation{
				Children: errs,
			}
		}
		return nil
	}
	return ErrValidation{
		error: fmt.Errorf(`value has type "%s" but must be slice or array`, rv.Kind()),
	}
}

// UnmarshalJSON is implemented to check the "type" field and to support
// dynamic unmarshaling of the Items field.
func (a *Array) UnmarshalJSON(data []byte) error {
	var raw struct {
		Type  string          `json:"type"`
		Items json.RawMessage `json:"items"`
	}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return fmt.Errorf(`unmarshal array json: "%s"`, err)
	}
	if raw.Type != a.Type() {
		return fmt.Errorf(`cannot read type "%s" into %s`, raw.Type, a.Type())
	}
	if a.Items, err = SchemaUnmarshalJSON(raw.Items); err != nil {
		return fmt.Errorf(`unmarshal array.items json: "%s"`, err)
	}
	return nil
}

// MarshalJSON adds the "type" field and validates before marshaling.
func (a Array) MarshalJSON() ([]byte, error) {
	if err := a.Valid(); err != nil {
		return nil, err
	}
	raw := struct {
		Type  string `json:"type"`
		Items Schema `json:"items"`
	}{
		Type:  a.Type(),
		Items: a.Items,
	}
	return json.Marshal(raw)
}
