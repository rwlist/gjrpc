package astinfo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindPackagePath(t *testing.T) {
	wd, err := os.Getwd()
	assert.NoError(t, err)

	path, err := findPackagePath(wd)
	assert.NoError(t, err)

	assert.Equal(t, "github.com/rwlist/gjrpc/internal/gen/astinfo", path)
}
