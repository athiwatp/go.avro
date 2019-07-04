package avro

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// Enum represents the "enum" complex type.
type Enum struct {
	NameFields
	Doc     string   `json:"doc,omitempty"`
	Symbols []string `json:"symbols"`
}

func (e Enum) Type() string { return "enum" }

func (e Enum) Valid() error {
	errs := map[string]error{}
	if err := e.NameFields.Valid(); err != nil {
		errs["name"] = err
	}
	symMap := map[string]struct{}{}
	for _, sym := range e.Symbols {
		errKey := fmt.Sprintf(`symbol "%s"`, sym)
		if !nameRegex.MatchString(sym) {
			errs[errKey] = errors.New(`invalid symbol name`)
		}
		if _, ok := symMap[sym]; ok {
			errs[errKey] = errors.New(`duplicate symbol`)
		}
		symMap[sym] = struct{}{}
	}
	if len(errs) > 0 {
		return ErrValidation{
			Children: errs,
		}
	}
	return nil
}

func (e Enum) Validate(v interface{}) error {
	if v == nil {
		return errors.New(`nil is not a valid enum`)
	}
	if err := e.Valid(); err != nil {
		return fmt.Errorf(`validation aborted, enum schema is invalid: %s`, err)
	}

	// Static check for strings.
	switch s := v.(type) {
	case *string:
		return e.exists(*s)
	case string:
		return e.exists(s)
	}

	// Reflect for custom types with string concrete type.
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if k := rv.Kind(); k != reflect.String {
		return fmt.Errorf(`value of type "%s" is not a valid enum`, k)
	}
	return e.exists(rv.String())
}

func (e Enum) exists(s string) error {
	for _, sym := range e.Symbols {
		if s == sym {
			return nil
		}
	}
	return fmt.Errorf(`symbol "%s" does not exist in the enum`, s)
}

type jsonEnum struct {
	Type string `json:"type"`
	NameFields
	Doc     string   `json:"doc,omitempty"`
	Symbols []string `json:"symbols"`
}

// UnmarshalJSON is implemented to check the "type" field.
func (e *Enum) UnmarshalJSON(data []byte) error {
	raw := jsonEnum{}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return fmt.Errorf(`unmarshal enum json: "%s"`, err)
	}
	if raw.Type != e.Type() {
		return fmt.Errorf(`cannot read type "%s" into %s`, raw.Type, e.Type())
	}
	e.NameFields = raw.NameFields
	e.Doc = raw.Doc
	e.Symbols = raw.Symbols
	return nil
}

// MarshalJSON adds the "type" field and validates before marshaling.
func (e Enum) MarshalJSON() ([]byte, error) {
	if err := e.Valid(); err != nil {
		return nil, err
	}
	raw := jsonEnum{
		Type:       e.Type(),
		NameFields: e.NameFields,
		Doc:        e.Doc,
		Symbols:    e.Symbols,
	}
	return json.Marshal(&raw)
}
