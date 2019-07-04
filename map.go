package avro

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// Map represents the "map" complex typa.
type Map struct {
	Values Schema `json:"values"`
}

func (m Map) Type() string {
	return "map"
}

func (m Map) Valid() error {
	if m.Values == nil {
		return ErrValidation{
			Children: map[string]error{
				"values": errors.New("cannot be nil"),
			},
		}
	}
	if err := m.Values.Valid(); err != nil {
		return ErrValidation{
			Children: map[string]error{
				"values": err,
			},
		}
	}
	return nil
}

func (m Map) Validate(v interface{}) error {
	if v == nil {
		return errors.New(`nil is not a valid map`)
	}
	if err := m.Valid(); err != nil {
		return fmt.Errorf(`validation aborted, map schema is invalid: %s`, err)
	}

	// Static check for map[string]interface{} and pointer.
	var msi map[string]interface{}
	if p, ok := v.(*map[string]interface{}); ok {
		msi = *p
		if msi == nil {
			return errors.New(`pointer to nil map is not a valid map`)
		}
	}
	if mv, ok := v.(map[string]interface{}); ok {
		msi = mv
		if msi == nil {
			return errors.New(`nil map is not a valid map`)
		}
	}
	if msi != nil {
		errs := map[string]error{}
		for k, val := range msi {
			if err := m.Values.Validate(val); err != nil {
				errs[k] = err
			}
		}
		if len(errs) > 0 {
			return ErrValidation{
				Children: errs,
			}
		}
	}

	// Check for other map types and custom types with reflect.
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if kind := rv.Kind(); kind != reflect.Map {
		return fmt.Errorf(`value with type "%s" is not a valid Map`, kind)
	}
	if keyKind := rv.Type().Key().Kind(); keyKind != reflect.String {
		return fmt.Errorf(`map key has type "%s" but it must be string`, keyKind)
	}
	errs := map[string]error{}
	for _, k := range rv.MapKeys() {
		val := rv.MapIndex(k).Interface()
		if err := m.Values.Validate(val); err != nil {
			errs[k.String()] = err
		}
	}
	if len(errs) > 0 {
		return ErrValidation{
			Children: errs,
		}
	}
	return nil
}

// UnmarshalJSON is implemented to check the "type" field and to support
// dynamic unmarshaling of the Values type.
func (m *Map) UnmarshalJSON(data []byte) error {
	var raw struct {
		Type   string          `json:"type"`
		Values json.RawMessage `json:"values"`
	}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return fmt.Errorf(`unmarshal map json: "%s"`, err)
	}
	if raw.Type != m.Type() {
		return fmt.Errorf(`cannot read type "%s" into %s`, raw.Type, m.Type())
	}
	if m.Values, err = SchemaUnmarshalJSON(raw.Values); err != nil {
		return fmt.Errorf(`unmarshal map.values json: "%s"`, err)
	}
	return nil
}

// MarshalJSON adds the "type" field and validates before marshaling.
func (m Map) MarshalJSON() ([]byte, error) {
	if err := m.Valid(); err != nil {
		return nil, err
	}
	raw := struct {
		Type   string `json:"type"`
		Values Schema `json:"values"`
	}{
		Type:   m.Type(),
		Values: m.Values,
	}
	return json.Marshal(raw)
}
