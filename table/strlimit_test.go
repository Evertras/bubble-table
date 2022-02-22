package table

import (
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func TestLimitStr(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		max      int
		expected string
	}{
		{
			name:     "Short",
			input:    "Hello",
			max:      50,
			expected: "Hello",
		},
		{
			name:     "Close",
			input:    "Hello",
			max:      6,
			expected: "Hello",
		},
		{
			name:     "Equal",
			input:    "Hello",
			max:      5,
			expected: "Hello",
		},
		{
			name:     "Shorter",
			input:    "Hello this is a really long string",
			max:      8,
			expected: "Hello t…",
		},
		{
			name:     "Zero max",
			input:    "Hello",
			max:      0,
			expected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := limitStr(test.input, test.max)

			assert.Equal(t, test.expected, output)
			assert.LessOrEqual(t, utf8.RuneCountInString(output), test.max)
		})
	}
}
