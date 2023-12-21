package flect

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"unsafe"
)

type Model[T any] struct {
	attrs *attrsMap[T]
}

func NewModel[T any](deserializer Deserializer) Model[T] {
	if deserializer == nil {
		deserializer = new(NameDeserializer)
	}

	var zero [0]T
	typ := reflect.TypeOf(zero).Elem()
	if typ.Kind() != reflect.Struct {
		// TODO: accept maps as well, it's gonna be easy to implement as they don't need
		//  any special treatment
		panic("not a struct")
	}

	attrs := new(attrsMap[T])
	deserializeFields(0, attrs, deserializer, typ)

	return Model[T]{
		attrs: attrs,
	}
}

func deserializeFields[T any](offset uintptr, attrs *attrsMap[T], deserializer Deserializer, typ reflect.Type) {
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		if field.Type.Kind() == reflect.Struct {
			deserializer.Descend(field)
			deserializeFields(offset+field.Offset, attrs, deserializer, field.Type)
			deserializer.Ascend()
			continue
		}

		name := deserializer.Visit(field)
		baseType := baseTypeOf(field.Type)
		if baseType == Unknown {
			panic("cannot determine the base type")
		}

		attrs.Insert(name, Field[T]{
			Type: baseType,
			meta: fieldMeta{
				Size:   field.Type.Size(),
				Offset: offset + field.Offset,
			},
		})
	}
}

func baseTypeOf(typ reflect.Type) BaseType {
	switch typ.Kind() {
	case reflect.String:
		return String
	case reflect.Int:
		return Int
	case reflect.Int8:
		return I8
	case reflect.Int16:
		return I16
	case reflect.Int32:
		return I32
	case reflect.Int64:
		return I64
	case reflect.Uint:
		return Uint
	case reflect.Uint8:
		return U8
	case reflect.Uint16:
		return U16
	case reflect.Uint32:
		return U32
	case reflect.Uint64:
		return U64
	case reflect.Bool:
		return Bool
	}

	return Unknown
}

func (m Model[T]) String() string {
	entries := m.attrs.Entries()

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Field.meta.Offset < entries[j].Field.meta.Offset
	})

	var fields []string

	for _, entry := range entries {
		fields = append(fields, fmt.Sprintf(
			"%s: size=%d offset=%d", entry.Key, entry.Field.meta.Size, entry.Field.meta.Offset,
		))
	}

	return fmt.Sprintf(
		"Model[%s]{%s}", reflect.TypeOf([0]T{}).Elem().Name(), strings.Join(fields, ", "),
	)
}

func (m Model[T]) Field(key string) (Field[T], bool) {
	return m.attrs.Lookup(key)
}

type Param struct {
	Key   string
	Value unsafe.Pointer
}

func Instantiate[T any](model Model[T], params ...Param) T {
	var zero T

	for _, param := range params {
		field, found := model.Field(param.Key)
		if !found {
			continue
		}

		zero = field.WriteUPtr(zero, param.Value)
	}

	return zero
}

type Field[T any] struct {
	Type BaseType
	meta fieldMeta
}

func (f Field[T]) Size() uint {
	return uint(f.meta.Size)
}

// WriteUPtr takes unsafe pointer and writes as much data from it, as the type
// consumes
func (f Field[T]) WriteUPtr(into T, src unsafe.Pointer) T {
	dst := unsafe.Add(unsafe.Pointer(&into), f.meta.Offset)
	memcpy(dst, src, f.meta.Size)

	return into
}

func (f Field[T]) WriteBool(into T, b bool) T {
	return f.WriteUPtr(into, unsafe.Pointer(&b))
}

func (f Field[T]) WriteUInt(into T, num uint) T {
	return f.WriteUPtr(into, unsafe.Pointer(&num))
}

func (f Field[T]) WriteUInt8(into T, num uint8) T {
	return f.WriteUPtr(into, unsafe.Pointer(&num))
}

func (f Field[T]) WriteUInt16(into T, num uint16) T {
	return f.WriteUPtr(into, unsafe.Pointer(&num))
}

func (f Field[T]) WriteUInt32(into T, num uint32) T {
	return f.WriteUPtr(into, unsafe.Pointer(&num))
}

func (f Field[T]) WriteUInt64(into T, num uint64) T {
	return f.WriteUPtr(into, unsafe.Pointer(&num))
}

func (f Field[T]) WriteInt(into T, num int) T {
	return f.WriteUPtr(into, unsafe.Pointer(&num))
}

func (f Field[T]) WriteInt8(into T, num int8) T {
	return f.WriteUPtr(into, unsafe.Pointer(&num))
}

func (f Field[T]) WriteInt16(into T, num int16) T {
	return f.WriteUPtr(into, unsafe.Pointer(&num))
}

func (f Field[T]) WriteInt32(into T, num int32) T {
	return f.WriteUPtr(into, unsafe.Pointer(&num))
}

func (f Field[T]) WriteInt64(into T, num int64) T {
	return f.WriteUPtr(into, unsafe.Pointer(&num))
}

func (f Field[T]) WriteString(into T, value string) T {
	return f.WriteUPtr(into, unsafe.Pointer(&value))
}

func memcpy(dst, src unsafe.Pointer, size uintptr) {
	copy(unsafe.Slice((*byte)(dst), size), unsafe.Slice((*byte)(src), size))
}

type fieldMeta struct {
	Size, Offset uintptr
}

type Deserializer interface {
	Descend(field reflect.StructField)
	Ascend()
	Visit(field reflect.StructField) string
}

type NameDeserializer struct {
	stack []string
}

func (n *NameDeserializer) Descend(field reflect.StructField) {
	n.stack = append(n.stack, field.Name)
}

func (n *NameDeserializer) Ascend() {
	n.stack = n.stack[:len(n.stack)-1]
}

func (n *NameDeserializer) Visit(field reflect.StructField) (name string) {
	if len(n.stack) > 0 {
		name = strings.Join(n.stack, ".") + "."
	}

	return name + field.Name
}
