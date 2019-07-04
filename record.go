package avro

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// Record represents the "record" complex type.
type Record struct {
	NameFields
	Doc    string  `json:"doc,omitempty"`
	Fields []Field `json:"fields"`
}

type Field struct {
	Name    string       `json:"name"`
	Doc     string       `json:"doc,omitempty"`
	Type    Schema       `json:"type"`
	Default *interface{} `json:"default,omitempty"`
	Order   string       `json:"order,omitempty"`
	Aliases []string     `json:"aliases,omitempty"`
}

func (r Record) Type() string { return "record" }

func (r Record) Valid() error {
	if err := r.NameFields.Valid(); err != nil {
		return err
	}
	errs := map[string]error{}
	for i, f := range r.Fields {
		name := fmt.Sprintf(`field #%d ("%s")`, i, f.Name)
		if !nameRegex.MatchString(f.Name) {
			errs[name] = errors.New("invalid name")
		}
		// Check if the Type is missing.
		if f.Type == nil {
			errs[fmt.Sprintf(`%s type`, name)] = errors.New("missing type")
			continue
		}
		// Check if the Type is valid.
		if err := f.Type.Valid(); err != nil {
			errs[fmt.Sprintf(`%s type`, name)] = err
			continue
		}
		// Check if the Default value is valid.
		if f.Default != nil {
			if err := f.Type.Validate(f.Default); err != nil {
				errs[fmt.Sprintf(`%s default`, name)] = err
			}
		}
		// Check if the Order is valid.
		switch f.Order {
		case "", "ascending", "descending", "ignore":
		default:
			errs[fmt.Sprintf(`%s order`, name)] = fmt.Errorf(`"%s" is not a valid value`, f.Order)
		}
	}
	if len(errs) > 0 {
		return ErrValidation{
			Children: errs,
		}
	}
	return nil
}

func (r Record) Validate(v interface{}) error {
	if v == nil {
		return errors.New(`nil is not a valid record`)
	}
	if err := r.Valid(); err != nil {
		return fmt.Errorf(`validation aborted, record schema is invalid: %s`, err)
	}

	// Static check for map[string]interface{} and pointer.
	var msi map[string]interface{}
	if p, ok := v.(*map[string]interface{}); ok {
		msi = *p
		if msi == nil {
			return errors.New(`pointer to nil map is not a valid record`)
		}
	}
	if mv, ok := v.(map[string]interface{}); ok {
		msi = mv
		if msi == nil {
			return errors.New(`nil map is not a valid record`)
		}
	}
	if msi != nil {
		errs := map[string]error{}
		for k, val := range msi {
			f, ok := r.GetField(k)
			if !ok {
				errs[k] = errors.New("record does not have a field with this name")
				continue
			}
			if err := f.Type.Validate(val); err != nil {
				errs[k] = err
			}
		}
		if len(errs) > 0 {
			return ErrValidation{
				Children: errs,
			}
		}
		return nil
	}

	// Reflect check for structs and custom map types.
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Struct:
		errs := map[string]error{}
		numField := rv.NumField()
		t := rv.Type()
		for i := 0; i < numField; i++ {
			rf := t.Field(i)
			// Check that name or tag exists in record type.
			var name string
			if tag := rf.Tag.Get("avro"); tag != "" {
				name = tag
			} else {
				name = rf.Name
			}
			field, ok := r.GetField(name)
			if !ok {
				errs[name] = errors.New("record does not have a field with this name")
				continue
			}
			// Check if field value is valid.
			if err := field.Type.Validate(rv.Field(i).Interface()); err != nil {
				errs[name] = err
			}
		}
		if len(errs) > 0 {
			return ErrValidation{
				Children: errs,
			}
		}
		return nil
	case reflect.Map:
		if keyKind := rv.Type().Key().Kind(); keyKind != reflect.String {
			return fmt.Errorf(`map key has type "%s" but it must be string`, keyKind)
		}
		errs := map[string]error{}
		for _, k := range rv.MapKeys() {
			name := k.String()
			field, ok := r.GetField(name)
			if !ok {
				errs[name] = errors.New("record does not have a field with this name")
				continue
			}
			val := rv.MapIndex(k).Interface()
			if err := field.Type.Validate(val); err != nil {
				errs[name] = err
			}
		}
		if len(errs) > 0 {
			return ErrValidation{
				Children: errs,
			}
		}
		return nil
	}

	return fmt.Errorf(`value with type "%s" is not a valid record`, rv.Kind())
}

func (r Record) GetField(name string) (*Field, bool) {
	for _, f := range r.Fields {
		if f.Name == name {
			return &f, true
		}
	}
	return nil, false
}

// UnmarshalJSON is implemented to support dynamic unmarshaling of Field Types.
func (f Field) UnmarshalJSON(data []byte) error {
	var raw struct {
		Name    string          `json:"name"`
		Doc     string          `json:"doc,omitempty"`
		Type    json.RawMessage `json:"type"`
		Default *interface{}    `json:"default,omitempty"`
		Order   string          `json:"order,omitempty"`
		Aliases []string        `json:"aliases,omitempty"`
	}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return fmt.Errorf(`unmarshal field json: "%s"`, err)
	}
	if f.Type, err = SchemaUnmarshalJSON(raw.Type); err != nil {
		return fmt.Errorf(`unmarshal field.type json: "%s"`, err)
	}
	f.Name = raw.Name
	f.Doc = raw.Doc
	f.Default = raw.Default
	f.Order = raw.Order
	f.Aliases = raw.Aliases
	return nil
}

// MarshalJSON adds the "type" field and validates before marshaling.
func (r Record) MarshalJSON() ([]byte, error) {
	return nil, nil
}
