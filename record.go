package avro

import (
	"encoding/json"
	"fmt"
)

// Record represents the "record" complex type.
type Record struct {
	NamedSchema
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

func (r Record) Type() string                 { return "record" }
func (r Record) Valid() error                 { return nil }
func (r Record) Validate(v interface{}) error { return nil }

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
