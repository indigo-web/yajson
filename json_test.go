package yajson

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJSON(t *testing.T) {
	j := `
	{
		"Boah": "any text here",
		"Something": "some text inside of the string",
		"Nothing": "Hello, world!",
		"any string": "This must never appear"
	}
	`

	t.Run("myJSONModel", func(t *testing.T) {
		parser := New[basicStringModel]()
		model, err := parser.Parse(j)
		require.NoError(t, err)
		assert.Equal(t, "some text inside of the string", model.Something)
		assert.Equal(t, "Hello, world!", model.Nothing)
		assert.Equal(t, "any text here", model.Boah)
	})
}
