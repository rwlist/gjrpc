package gen

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/rwlist/gjrpc/internal/gen/argparse"

	"github.com/stretchr/testify/assert"
)

var (
	_, b, _, _ = runtime.Caller(0)

	// examples directory
	examplesDir = filepath.Join(filepath.Dir(b), "../../examples")
)

func TestExample1(t *testing.T) {
	genFile := filepath.Join(examplesDir, "1_simple_service/gen_server/gen_router.go")

	assertEqualFile(t, genFile, func() error {
		assert.NoError(t, os.Chdir(filepath.Dir(genFile)))

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
	routerGenGo := filepath.Join(examplesDir, "2_cookie_oauth/gen_server/gen_router.go")
	protoTs := filepath.Join(examplesDir, "2_cookie_oauth/gen_client/proto.ts")

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
