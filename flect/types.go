package flect

type BasicType uint8

const (
	String BasicType = iota + 1
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
