package table

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAsInt(t *testing.T) {
	check := func(data interface{}, isInt bool, expectedValue int64) {
		val, ok := asInt(data)
		assert.Equal(t, isInt, ok)
		assert.Equal(t, expectedValue, val)
	}

	check(3, true, 3)
	check(3.3, false, 0)
	check(int8(3), true, 3)
	check(int16(3), true, 3)
	check(int32(3), true, 3)
	check(int64(3), true, 3)
	check(uint(3), true, 3)
	check(uint8(3), true, 3)
	check(uint16(3), true, 3)
	check(uint32(3), true, 3)
	check(uint64(3), true, 3)
	check(StyledCell{Data: 3}, true, 3)
	check(time.Duration(3), true, 3)
}

func TestAsNumber(t *testing.T) {
	check := func(data interface{}, isFloat bool, expectedValue float64) {
		val, ok := asNumber(data)
		assert.Equal(t, isFloat, ok)
		assert.InDelta(t, expectedValue, val, 0.001)
	}

	check(uint32(3), true, 3)
	check(3.3, true, 3.3)
	check(float32(3.3), true, 3.3)
	check(StyledCell{Data: 3.3}, true, 3.3)
}
