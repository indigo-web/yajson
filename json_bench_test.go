package yajson

import (
	"encoding/json"
	"github.com/indigo-web/utils/uf"
	"github.com/romshark/jscan/v2"
	"testing"
)

func BenchmarkJSON(b *testing.B) {
	j := `
	{
		"Boah": "any text here",
		"Something": "okay, let it be",
		"Nothing": "Hello, world!",
		"any string": "This must never appear"
	}
	`

	parser := New[basicStringModel]()
	m := basicStringModel{}

	b.Run("my parser", func(b *testing.B) {
		b.SetBytes(int64(len(j)))
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			m, _ = parser.Parse(j)
		}
	})

	b.Run("without my parser", func(b *testing.B) {
		b.SetBytes(int64(len(j)))
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_ = jscan.Scan(j, func(i *jscan.Iterator[string]) (err bool) {
				return false
			})
		}
	})

	b.Run("std", func(b *testing.B) {
		model := any(&basicStringModel{})

		b.SetBytes(int64(len(j)))
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_ = json.Unmarshal(uf.S2B(j), model)
		}
	})

	keepalive(m)
}

func keepalive(basicStringModel) {}
