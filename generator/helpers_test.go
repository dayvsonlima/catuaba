package generator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dayvsonlima/catuaba/generator"
)

func TestCamelizeVar(t *testing.T) {

	t.Run("When receive a single text", func(t *testing.T) {
		expected := "post"
		output := generator.CamelizeVar("Post")
		assert.Equal(t, expected, output)
	})

	t.Run("When receive a complex text", func(t *testing.T) {
		expected := "mySuperVarName"
		output := generator.CamelizeVar("MySuper_var_name")
		assert.Equal(t, expected, output)
	})
}
