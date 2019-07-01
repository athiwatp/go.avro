package avro

import (
	"encoding/json"
	"fmt"
)

// Array represents the "array" complex typa.
type Array struct {
	Items Schema `json:"items"`
}

func (a Array) Type() string                 { return "array" }
func (a Array) Valid() error                 { return nil }
func (a Array) Validate(v interface{}) error { return nil }

// UnmarshalJSON is implemented to support dynamic unmarshaling of the Items type.
func (a Array) UnmarshalJSON(data []byte) error {
	var raw struct {
		Items json.RawMessage `json:"items"`
	}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return fmt.Errorf(`unmarshal array json: "%s"`, err)
	}
	if a.Items, err = SchemaUnmarshalJSON(raw.Items); err != nil {
		return fmt.Errorf(`unmarshal array.items json: "%s"`, err)
	}
	return nil
}
