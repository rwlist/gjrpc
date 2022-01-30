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

func TestExample2(t *testing.T) {
	routerGenGo, err := filepath.Abs("../../examples/2_cookie_oauth/gen_server/gen_router.go")
	assert.NoError(t, err)

	protoTs, err := filepath.Abs("../../examples/2_cookie_oauth/gen_client/proto.ts")
	assert.NoError(t, err)

	assertEqualFile(t, routerGenGo, func() error {
		assert.NoError(t, os.Chdir(filepath.Dir(routerGenGo)))

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

	assertEqualFile(t, protoTs, func() error {
		assert.NoError(t, os.Chdir(filepath.Dir(protoTs)))

		args := &argparse.CliArgs{
			Command: "gen:client:ts-proto",
			Args: map[string]string{
				"protoPkg": "../proto",
				"out":      "proto.ts",
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
