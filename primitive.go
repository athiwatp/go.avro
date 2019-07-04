package avro

import "fmt"

const (
	Null    Primitive = "null"
	Boolean Primitive = "boolean"
	Int     Primitive = "int"
	Long    Primitive = "long"
	Float   Primitive = "float"
	Double  Primitive = "double"
	Bytes   Primitive = "bytes"
	String  Primitive = "string"
)

var primitives = []Primitive{
	Null,
	Boolean,
	Int,
	Long,
	Float,
	Double,
	Bytes,
	String,
}

// Primitive represents any primitive type: null, boolean, int, long, float,
// double, bytes, and string.
type Primitive string

func (p Primitive) Type() string {
	return string(p)
}

func (p Primitive) Valid() error {
	switch p {
	case Null, Boolean, Int, Long, Float, Double, Bytes, String:
		return nil
	}
	return fmt.Errorf(`"%s" is not a valid primitive type`, p)
}

func (p Primitive) Validate(v interface{}) error {
	if p != Null && v == nil {
		return p.errInvalid()
	}
	if ptrIface, ok := v.(*interface{}); ok {
		v = *ptrIface
	}
	switch p {
	case Null:
		if v == nil {
			return nil
		}
	case Boolean:
		switch v.(type) {
		case bool, *bool:
			return nil
		}
	case Int:
		switch v.(type) {
		case int, int32, int64, uint, uint32, uint64,
			*int, *int32, *int64, *uint, *uint32, *uint64:
			return nil
		}
	case Long:
		switch v.(type) {
		case int64, uint64, *int64, *uint64:
			return nil
		case int, *int:
			if uint64(^uint(0)) == ^uint64(0) {
				// int is 64-bit
				return nil
			}
		}
	case Float:
		switch v.(type) {
		case float32, *float32:
			return nil
		}
	case Double:
		switch v.(type) {
		case float64, *float64:
			return nil
		}
	case Bytes:
		switch v.(type) {
		case []byte, *[]byte:
			return nil
		}
	case String:
		switch v.(type) {
		case string, *string:
			return nil
		}
	}
	return p.errInvalid()
}

func (p Primitive) errInvalid() error {
	return fmt.Errorf(`value is not a valid "%s"`, p)
}
