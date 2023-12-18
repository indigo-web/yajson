package flect

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"unsafe"
)

func TestFlect(t *testing.T) {
	type myStruct struct {
		A byte
		B uint16
		C int32
	}

	model := NewModel[myStruct](nil)

	t.Run("instantiate", func(t *testing.T) {
		m := Instantiate(model,
			Attr{"A", uptr(5)},
			Attr{"B", uptr(32769)},
			Attr{"C", uptr(-67108864)},
		)

		assert.Equal(t, uint8(5), m.A)
		assert.Equal(t, uint16(32769), m.B)
		assert.Equal(t, int32(-67108864), m.C)
	})

	t.Run("fill partially", func(t *testing.T) {
		m := Instantiate(model,
			Attr{"A", uptr(5)},
			Attr{"C", uptr(-67108864)},
		)

		assert.Equal(t, uint8(5), m.A)
		assert.Equal(t, uint16(0), m.B)
		assert.Equal(t, int32(-67108864), m.C)
	})

	t.Run("fill with unknown field", func(t *testing.T) {
		m := Instantiate(model,
			Attr{"A", uptr(5)},
			Attr{"G", uptr(123)},
			Attr{"C", uptr(-67108864)},
			Attr{"M", uptr(123)},
		)

		assert.Equal(t, uint8(5), m.A)
		assert.Equal(t, uint16(0), m.B)
		assert.Equal(t, int32(-67108864), m.C)
	})

	type nestedStructs struct {
		A string
		B struct {
			C string
			D struct {
				E string
			}
		}
	}

	t.Run("nested structs", func(t *testing.T) {
		model := NewModel[nestedStructs](new(NameDeserializer))
		m := Instantiate(model, Attr{
			Key:   "A",
			Value: uptr("foo"),
		}, Attr{
			Key:   "B.C",
			Value: uptr("bar"),
		}, Attr{
			Key:   "B.D.E",
			Value: uptr("spam"),
		})

		require.Equal(t, "foo", m.A)
		require.Equal(t, "bar", m.B.C)
		require.Equal(t, "spam", m.B.D.E)
	})
}

func uptr[T any](val T) unsafe.Pointer {
	// I'm not pretty sure, whether taking a pointer just here is safe enough. So
	// let it leak, anyway used in tests only
	return unsafe.Pointer(ptr(val))
}

func ptr[T any](val T) *T {
	return &val
}
