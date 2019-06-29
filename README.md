# avro
Modern Avro implementation for Go

## In development

This library is in development and is not production ready.

## Why another Avro implementation

There are already several Avro libraries for Go:

- [linkedin/goavro](https://github.com/linkedin/goavro)
- [elodina/go-avro](https://github.com/elodina/go-avro)
- [actgardner/gogen-avro](https://github.com/actgardner/gogen-avro)

However I could not find one which had all of the following features:

- Encode/Decode to custom types
- Schema Registry interface
- Full Avro specification implemented (including IDL/RPC)
- Native Schema type and expressive type semantics
- Fully dynamic at runtime (does not use code generation)
- Pluggable logical type decoding with sane defaults
- Avro HTTP RPC implementation for `net/rpc`
- Validation of arbitrary values against schemas

This library attempts to implement all of the above.
