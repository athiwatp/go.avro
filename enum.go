package avro

// Enum represents the "enum" complex type.
type Enum struct {
	NamedSchema
	Doc     string       `json:"doc,omitempty"`
	Symbols []string     `json:"symbols"`
	Default *interface{} `json:"default,omitempty"`
}

func (e Enum) Type() string                 { return "enum" }
func (e Enum) Valid() error                 { return nil }
func (e Enum) Validate(v interface{}) error { return nil }
