package example2

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestBlockedImport(t *testing.T) { //nolint
	assert.Equal(t, 1, 1)
}
