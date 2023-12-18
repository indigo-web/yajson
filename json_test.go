package yajson

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type basicStringModel struct {
	Something string `json:"some_thing"`
	Nothing   string
	Foo       string `json:"foo"`
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
		"foo": "any text here",
		"some_thing": "some text inside of the string",
		"Nothing": "Hello, world!",
		"any string": "This must never appear"
	}
	`

	t.Run("plain model", func(t *testing.T) {
		parser := New[basicStringModel]()
		model, err := parser.Parse(plain)
		require.NoError(t, err)
		assert.Equal(t, "some text inside of the string", model.Something)
		assert.Equal(t, "Hello, world!", model.Nothing)
		assert.Equal(t, "any text here", model.Foo)
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
}
