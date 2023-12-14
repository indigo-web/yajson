package yajson

import (
	"fmt"
	"github.com/romshark/jscan/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type basicStringModel struct {
	Something string `json:"Something"`
	Nothing   string `json:"Nothing"`
	Boah      string `json:"Boah"`
}

type nestedModel struct {
	Hello        string
	Alternatives struct {
		FirstAlternative string `json:"first_alternative"`
	}
}

func TestJSON(t *testing.T) {
	plain := `
	{
		"Boah": "any text here",
		"Something": "some text inside of the string",
		"Nothing": "Hello, world!",
		"any string": "This must never appear"
	}
	`

	t.Run("myJSONModel", func(t *testing.T) {
		parser := New[basicStringModel]()
		model, err := parser.Parse(plain)
		require.NoError(t, err)
		assert.Equal(t, "some text inside of the string", model.Something)
		assert.Equal(t, "Hello, world!", model.Nothing)
		assert.Equal(t, "any text here", model.Boah)
	})

	nested := `
	{
		"Hello": "World",
		"Alternatives": {
			"first_alternative": "no"
		}
	}
	`

	t.Run("nested model", func(t *testing.T) {
		parser := New[nestedModel]()
		model, err := parser.Parse(nested)
		require.NoError(t, err)
		assert.Equal(t, "World", model.Hello)
		assert.Equal(t, "no", model.Alternatives.FirstAlternative)
	})

	t.Run("test smth", func(t *testing.T) {
		err := jscan.Scan(nested, func(i *jscan.Iterator[string]) (err bool) {
			fmt.Printf("%q:\n", i.Pointer())
			fmt.Printf("├─ valueType:  %s\n", i.ValueType().String())
			if k := i.Key(); k != "" {
				fmt.Printf("├─ key:        %q\n", k[1:len(k)-1])
			}
			if ai := i.ArrayIndex(); ai != -1 {
				fmt.Printf("├─ arrayIndex: %d\n", ai)
			}
			if v := i.Value(); v != "" {
				fmt.Printf("├─ value:      %q\n", v)
			}
			fmt.Printf("└─ level:      %d\n", i.Level())
			return false // Resume scanning
		})

		require.False(t, err.IsErr(), err.Error())
	})
}
