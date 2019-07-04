package avro

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Union represents the "union" complex type.
type Union []Schema

func (u Union) Type() string {
	return "union"
}

func (u Union) Valid() error {
	if len(u) == 0 {
		return errors.New("union may not be empty")
	}
	typeMap := map[string]struct{}{} // for checking duplicates
	for i, s := range u {
		if s == nil {
			return fmt.Errorf(`schema #%d in union is nil`, i)
		}
		if err := s.Valid(); err != nil {
			return fmt.Errorf(`schema #%d in union is invalid: %s`, i, err)
		}
		// Duplicates are not allowed except for named types, which must have unique names.
		switch t := s.Type(); t {
		case "union":
			return fmt.Errorf(`schema #%d in union is nested union`, i)
		case "fixed", "enum", "record":
			key := typeKey(s)
			if _, ok := typeMap[key]; ok {
				return fmt.Errorf(`schema #%d in union is duplicate "%s"`, i, t)
			}
			typeMap[key] = struct{}{}
		default:
			key := typeKey(s)
			if _, ok := typeMap[key]; ok {
				return fmt.Errorf(`schema #%d in union is duplicate "%s"`, i, t)
			}
			typeMap[key] = struct{}{}
		}
	}
	return nil
}

func typeKey(s Schema) string {
	t := s.Type()
	switch t {
	case "fixed", "enum", "record":
		return fmt.Sprintf(`%s %s`, t, s.(NamedSchema).Fullname())
	}
	return t
}

func (u Union) Validate(v interface{}) error {
	if err := u.Valid(); err != nil {
		return err
	}
	errs := map[string]error{}
	for _, s := range u {
		if err := s.Validate(v); err != nil {
			errs[typeKey(s)] = err
		}
	}
	if len(errs) > 0 {
		return ErrValidation{
			error:    errors.New("value does not match any type in the union; here is a breakdown"),
			Children: errs,
		}
	}
	return nil
}

// UnmarshalJSON is implemented to support dynamic unmarshaling of contained types.
func (u *Union) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return fmt.Errorf(`unmarshal union json: "%s"`, err)
	}
	for _, rawContained := range raw {
		var contained Schema
		if contained, err = SchemaUnmarshalJSON(rawContained); err != nil {
			return fmt.Errorf(`unmarshal union contained json: "%s"`, err)
		}
		*u = append(*u, contained)
	}
	return nil
}

// MarshalJSON validates before marshaling.
func (u Union) MarshalJSON() ([]byte, error) {
	if err := u.Valid(); err != nil {
		return nil, err
	}
	return json.Marshal([]Schema(u))
}
