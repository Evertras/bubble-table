package table

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModelInitReturnsNil(t *testing.T) {
	model := New(nil)

	cmd := model.Init()

	assert.Nil(t, cmd)
}
