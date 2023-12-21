package flect

//go:generate stringer -type=BaseType
type BaseType uint8

const (
	Unknown BaseType = iota
	String
	Bool
	Array
	U8
	U16
	U32
	U64
	Uint
	I8
	I16
	I32
	I64
	Int
)
