package gen

import (
	"github.com/rwlist/gjrpc/internal/gen/argparse"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExample1(t *testing.T) {
	exampleDir, err := filepath.Abs("../../examples/1_simple_service/gen_server/gen_router.go")
	assert.NoError(t, err)

	assertEqualFile(t, exampleDir, func() error {
		assert.NoError(t, os.Chdir("../../examples/1_simple_service/gen_server"))

		args := &argparse.CliArgs{
			Command: "gen:server:router",
			Args: map[string]string{
				"protoPkg":       "../proto",
				"handlersStruct": "Handlers",
				"out":            "gen_router.go",
			},
		}
		FromCmdline(args)
		return nil
	})
}

func assertEqualFile(t *testing.T, path string, genFunc func() error) {
	before, err := os.ReadFile(path)
	assert.NoError(t, err)

	err = genFunc()
	assert.NoError(t, err)

	after, err := os.ReadFile(path)
	assert.NoError(t, err)

	assert.Equal(t, string(before), string(after))
}
