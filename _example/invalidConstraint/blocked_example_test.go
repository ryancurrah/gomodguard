package example2

import (
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
	"gotest.tools/v3/assert"
)

func TestBlockedImport(t *testing.T) { //nolint
	assert.Equal(t, errors.New("test"), gin.ErrorTypeBind)
}
