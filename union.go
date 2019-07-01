package avro

import (
	"encoding/json"
	"fmt"
)

// Union represents the "union" complex type.
type Union []Schema

func (u Union) Type() string                 { return "union" }
func (u Union) Valid() error                 { return nil }
func (u Union) Validate(v interface{}) error { return nil }

// UnmarshalJSON is implemented to support dynamic unmarshaling of contained types.
func (u Union) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return fmt.Errorf(`unmarshal union json: "%s"`, err)
	}
	for i, rawContained := range raw {
		var contained Schema
		if contained, err = SchemaUnmarshalJSON(rawContained); err != nil {
			return fmt.Errorf(`unmarshal union contained json: "%s"`, err)
		}
		u[i] = contained
	}
	return nil
}
