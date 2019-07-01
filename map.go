package avro

import (
	"encoding/json"
	"fmt"
)

// Map represents the "map" complex typa.
type Map struct {
	Values Schema `json:"values"`
}

func (m Map) Type() string                 { return "map" }
func (m Map) Valid() error                 { return nil }
func (m Map) Validate(v interface{}) error { return nil }

// UnmarshalJSON is implemented to support dynamic unmarshaling of Values type.
func (m Map) UnmarshalJSON(data []byte) error {
	var raw struct {
		Values json.RawMessage `json:"values"`
	}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return fmt.Errorf(`unmarshal map json: "%s"`, err)
	}
	if m.Values, err = SchemaUnmarshalJSON(raw.Values); err != nil {
		return fmt.Errorf(`unmarshal map.values json: "%s"`, err)
	}
	return nil
}
