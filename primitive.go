package avro

import "fmt"

const (
	Null    Primitive = "null"
	Boolean           = "boolean"
	Int               = "int"
	Long              = "long"
	Float             = "float"
	Double            = "double"
	Bytes             = "bytes"
	String            = "string"
)

// Primitive represents any primitive type: null, boolean, int, long, float,
// double, bytes, and string.
type Primitive string

func (p Primitive) Type() string { return string(p) }
func (p Primitive) Valid() error {
	switch p {
	case Null, Boolean, Int, Long, Float, Double, Bytes, String:
		return nil
	}
	return fmt.Errorf(`"%s" is not a valid primitive type`, p)
}
func (p Primitive) Validate(v interface{}) error { return nil }
