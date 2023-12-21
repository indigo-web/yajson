package flect

import (
	"testing"
	"unsafe"
)

func BenchmarkModelFiller(b *testing.B) {
	type myStruct struct {
		A byte
		B uint16
		C int32
	}

	b.Run("full", func(b *testing.B) {
		model := NewModel[myStruct](nil)
		fields := []Param{
			{"A", uptr(5)},
			{"B", uptr(32769)},
			{"C", uptr(67108864)},
		}
		b.ReportAllocs()
		b.SetBytes(int64(unsafe.Sizeof(myStruct{})))
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_ = Instantiate(model, fields...)
		}
	})

	b.Run("single field", func(b *testing.B) {
		model := NewModel[myStruct](nil)

		b.ReportAllocs()
		b.SetBytes(int64(unsafe.Sizeof(myStruct{}.C)))
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			field, _ := model.Field("C")
			_ = field.WriteUInt32(myStruct{}, 67108864)
		}
	})
}
