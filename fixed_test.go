package avro

import (
	"encoding/json"
	"testing"

	"github.com/matryer/is"
)

func TestFixed_Type(t *testing.T) {
	is := is.New(t)

	f := Fixed{}
	is.Equal(f.Type(), "fixed") // enum type is "fixed"
}

func TestFixed_Valid(t *testing.T) {
	is := is.New(t)

	f := Fixed{}
	is.True(f.Valid() != nil) // nameless fixed should be invalid

	f.Name = "Test"
	is.True(f.Valid() != nil) // fixed with 0 size should be invalid

	f.Size = 8
	is.NoErr(f.Valid()) // named fixed with non-zero size should be valid
}

func TestFixedValidate(t *testing.T) {
	is := is.New(t)
	vBytes := []byte{0xFF, 0xFF}
	vString := "\xFF\xFF"
	nilBytes := []byte(nil)

	f := Fixed{
		Size: 2,
	}
	is.True(f.Validate(vBytes) != nil) // invalid schema cannot validate

	f = Fixed{}
	f.Name = "Test"
	f.Size = 2
	is.NoErr(f.Validate(vBytes))             // []byte with length == size is valid
	is.NoErr(f.Validate(vString))            // string with byte length == size is valid
	is.NoErr(f.Validate(&vBytes))            // *[]byte with length == size is valid
	is.NoErr(f.Validate(&vString))           // *string with byte length == size is valid
	is.True(f.Validate([]byte{0xFF}) != nil) // []byte with length != size is invalid
	is.True(f.Validate(nil) != nil)          // nil is invalid
	is.True(f.Validate([]rune{'a'}) != nil)  // unsupported value type should error
	is.True(f.Validate(5) != nil)            // unsupported value type should error
	is.True(f.Validate(nilBytes) != nil)     // nil []byte is invalid
	is.True(f.Validate(&nilBytes) != nil)    // ptr to nil []byte is invalid

	type customBytes []byte
	customBytesVal := customBytes([]byte{0xFF, 0xFF})
	type customStr string
	customStrVal := customStr("\xFF\xFF")
	is.NoErr(f.Validate(customBytesVal))  // only concrete type should matter
	is.NoErr(f.Validate(&customBytesVal)) // only concrete type should matter
	is.NoErr(f.Validate(customStrVal))    // only concrete type should matter
	is.NoErr(f.Validate(&customStrVal))   // only concrete type should matter
}

func TestFixed_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		Name  string
		Spec  []byte
		Error bool
		Size  uint
	}{
		{"valid spec", []byte(`{"type":"fixed","size":2}`), false, 2},
		{"empty json", []byte(``), true, 0},
		{"invalid json", []byte(`{`), true, 0},
		{"invalid type", []byte(`{"type":"__WRONG__","size":2}`), true, 0},
		{"invalid size", []byte(`{"type":"fixed","size":"__WRONG__"}`), true, 0},
	}

	for _, test := range tests {
		test := test // capture range variable
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			f := Fixed{}
			err := f.UnmarshalJSON(test.Spec)
			if test.Error {
				is.True(err != nil) // should return error
			} else {
				is.NoErr(err)               // should not return error
				is.Equal(f.Size, test.Size) // unexpected size
			}
		})
	}
}

func TestFixed_MarshalJSON(t *testing.T) {
	is := is.New(t)

	f := Fixed{}
	b, err := json.Marshal(f)
	is.True(err != nil) // invalid fixed should not marshal

	f.Name = "Test"
	f.Size = 2
	spec := []byte(`{"type":"fixed","name":"Test","size":2}`)
	b, err = json.Marshal(&f)
	is.NoErr(err)                     // marshal valid schema no error
	is.Equal(string(b), string(spec)) // marshal valid spec
}
