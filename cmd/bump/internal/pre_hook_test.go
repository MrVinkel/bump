package internal_test

import (
	"bytes"
	"testing"

	"github.com/mrvinkel/bump/cmd/bump/internal"
	"github.com/stretchr/testify/assert"
)

func TestRunSuccess(t *testing.T) {
	shell := "sh -c -"
	commands := []string{"echo hello", "echo world", "echo $FOO"}
	var writer bytes.Buffer
	env := map[string]string{"FOO": "bar"}

	err := internal.Run(shell, commands, &writer, env)

	assert.NoError(t, err)
	assert.Equal(t, "hello\nworld\nbar\n", writer.String())
}

func TestRunFailure(t *testing.T) {
	shell := "sh -c -"
	commands := []string{"exit 1"}
	var writer bytes.Buffer
	var env map[string]string

	err := internal.Run(shell, commands, &writer, env)

	assert.Error(t, err)
	assert.Empty(t, writer.String())
}
