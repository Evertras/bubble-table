package table

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// A bit overkill but let's be thorough!

func TestMin(t *testing.T) {
	tests := []struct {
		x        int
		y        int
		expected int
	}{
		{
			x:        3,
			y:        4,
			expected: 3,
		},
		{
			x:        3,
			y:        3,
			expected: 3,
		},
		{
			x:        -4,
			y:        3,
			expected: -4,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%d and %d gives %d", test.x, test.y, test.expected), func(t *testing.T) {
			result := min(test.x, test.y)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		x        int
		y        int
		expected int
	}{
		{
			x:        3,
			y:        4,
			expected: 4,
		},
		{
			x:        3,
			y:        3,
			expected: 3,
		},
		{
			x:        -4,
			y:        3,
			expected: 3,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%d and %d gives %d", test.x, test.y, test.expected), func(t *testing.T) {
			result := max(test.x, test.y)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestGCD(t *testing.T) {
	tests := []struct {
		x        int
		y        int
		expected int
	}{
		{
			x:        3,
			y:        4,
			expected: 1,
		},
		{
			x:        3,
			y:        6,
			expected: 3,
		},
		{
			x:        4,
			y:        6,
			expected: 2,
		},
		{
			x:        1000,
			y:        100000,
			expected: 1000,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%d and %d has GCD %d", test.x, test.y, test.expected), func(t *testing.T) {
			result := gcd(test.x, test.y)
			assert.Equal(t, test.expected, result)
		})
	}
}
