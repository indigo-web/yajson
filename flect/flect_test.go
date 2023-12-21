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

	t.Run("instantiate", func(t *testing.T) {
		model := NewModel[myStruct](nil)
		m := Instantiate(model,
			Param{"A", uptr(5)},
			Param{"B", uptr(32769)},
			Param{"C", uptr(-67108864)},
		)

		assert.Equal(t, uint8(5), m.A)
		assert.Equal(t, uint16(32769), m.B)
		assert.Equal(t, int32(-67108864), m.C)
	})

	t.Run("fill partially", func(t *testing.T) {
		model := NewModel[myStruct](nil)
		m := Instantiate(model,
			Param{"A", uptr(5)},
			Param{"C", uptr(-67108864)},
		)

		assert.Equal(t, uint8(5), m.A)
		assert.Equal(t, uint16(0), m.B)
		assert.Equal(t, int32(-67108864), m.C)
	})

	t.Run("fill with unknown field", func(t *testing.T) {
		model := NewModel[myStruct](nil)
		m := Instantiate(model,
			Param{"A", uptr(5)},
			Param{"G", uptr(123)},
			Param{"C", uptr(-67108864)},
			Param{"M", uptr(123)},
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
		m := Instantiate(model, Param{
			Key:   "A",
			Value: uptr("foo"),
		}, Param{
			Key:   "B.C",
			Value: uptr("bar"),
		}, Param{
			Key:   "B.D.E",
			Value: uptr("spam"),
		})

		require.Equal(t, "foo", m.A)
		require.Equal(t, "bar", m.B.C)
		require.Equal(t, "spam", m.B.D.E)
	})

	type (
		token string
		maybe bool
	)

	type differentFieldTypes struct {
		A token
		B int
		C maybe
		D string
		E uint8
	}

	t.Run("different field types", func(t *testing.T) {
		model := NewModel[differentFieldTypes](nil)
		assert.Equal(t, String.String(), field(t, model, "A").Type.String())
		assert.Equal(t, Int.String(), field(t, model, "B").Type.String())
		assert.Equal(t, Bool.String(), field(t, model, "C").Type.String())
		assert.Equal(t, String.String(), field(t, model, "D").Type.String())
		assert.Equal(t, U8.String(), field(t, model, "E").Type.String())

	})
}

func field[T any](t *testing.T, model Model[T], key string) Field[T] {
	f, found := model.Field(key)
	require.Truef(t, found, "key not found: %s", key)

	return f
}

func uptr[T any](val T) unsafe.Pointer {
	// I'm not pretty sure, whether taking a pointer just here is safe enough. So
	// let it leak, anyway used in tests only
	return unsafe.Pointer(ptr(val))
}

func ptr[T any](val T) *T {
	return &val
}
