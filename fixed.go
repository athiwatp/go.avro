package avro

// Fixed represents the "fixed" complex type.
type Fixed struct {
	NamedSchema
	Size uint `json:"size"`
}

func (f Fixed) Type() string                 { return "fixed" }
func (f Fixed) Valid() error                 { return nil }
func (f Fixed) Validate(v interface{}) error { return nil }
