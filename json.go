package yajson

import (
	"errors"
	"github.com/indigo-web/utils/buffer"
	"github.com/indigo-web/utils/uf"
	"github.com/indigo-web/yajson/flect"
	"github.com/romshark/jscan/v2"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

const (
	BufferSpaceMin = 1024
	BufferSpaceMax = 127 * 1024
)

var (
	ErrNoSpace      = errors.New("no space for values")
	ErrTypeMismatch = errors.New("JSON carries an entry with an incompatible type")
	ErrUnknownType  = errors.New("the type is not supported by yajson")
)

type JSON[T any] struct {
	model  flect.Model[T]
	buffer *buffer.Buffer
}

// New returns a JSON parser instance for one specific model, defined via the generic
func New[T any]() *JSON[T] {
	return &JSON[T]{
		model:  flect.NewModel[T](new(pathDeserializer)),
		buffer: buffer.New(BufferSpaceMin, BufferSpaceMax),
	}
}

// Parse parses the passed JSON and stores the values into a new instance of the model
func (j *JSON[T]) Parse(input string) (result T, err error) {
	jsonErr := jscan.Scan(input, func(i *jscan.Iterator[string]) (exit bool) {
		key := i.Pointer()
		if len(key) == 0 {
			return false
		}

		field, found := j.model.Field(key)
		if !found {
			return false
		}

		switch field.Type {
		case flect.String:
			if i.ValueType() != jscan.ValueTypeString {
				err = ErrTypeMismatch
				return true
			}

			if !j.buffer.Append(uf.S2B(i.Value())) {
				err = ErrNoSpace
				return true
			}

			value := uf.B2S(j.buffer.Finish())
			result = field.WriteString(result, value[1:len(value)-1])
		case flect.Bool:
			switch i.ValueType() {
			case jscan.ValueTypeTrue, jscan.ValueTypeFalse:
			default:
				err = ErrTypeMismatch
				return true
			}

			result = field.WriteBool(result, i.ValueType() == jscan.ValueTypeTrue)
		case flect.Int, flect.I8, flect.I16, flect.I32, flect.I64:
			if i.ValueType() != jscan.ValueTypeNumber {
				err = ErrTypeMismatch
				return true
			}

			num, er := strconv.ParseInt(i.Value(), 10, int(field.Size()))
			if er != nil {
				err = ErrTypeMismatch
				return true
			}

			result = field.WriteUPtr(result, unsafe.Pointer(&num))
		case flect.Uint, flect.U8, flect.U16, flect.U32, flect.U64:
			if i.ValueType() != jscan.ValueTypeNumber {
				err = ErrTypeMismatch
				return true
			}

			num, er := strconv.ParseUint(i.Value(), 10, int(field.Size()))
			if er != nil {
				err = ErrTypeMismatch
				return true
			}

			result = field.WriteUPtr(result, unsafe.Pointer(&num))
		default:
			err = ErrUnknownType
			return true
		}

		return false
	})

	if jsonErr.IsErr() {
		err = jsonErr
	}

	// this doesn't destroy data we've written, but will override them on the next call
	j.buffer.Clear()

	return result, err
}

type pathDeserializer struct {
	stack []string
}

func (p *pathDeserializer) Descend(field reflect.StructField) {
	p.stack = append(p.stack, field.Name)
}

func (p *pathDeserializer) Ascend() {
	p.stack = p.stack[:len(p.stack)-1]
}

func (p *pathDeserializer) Visit(field reflect.StructField) (name string) {
	name = "/"

	if len(p.stack) > 0 {
		name += strings.Join(p.stack, "/") + "/"
	}

	if tag, found := field.Tag.Lookup("json"); found {
		return name + tag
	}

	return name + field.Name
}
