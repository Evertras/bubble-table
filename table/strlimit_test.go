package table

import (
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

// This function is only long because of repetitive test definitions, this is fine
// nolint: funlen
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
		{
			name:     "Unicode width",
			input:    "✓",
			max:      1,
			expected: "✓",
		},
		{
			name:     "Unicode truncated",
			input:    "✓✓✓",
			max:      2,
			expected: "✓…",
		},
		{
			name:     "Unicode japenese equal",
			input:    "直立",
			max:      5,
			expected: "直立",
		},
		{
			name:     "Unicode japenese truncated",
			input:    "直立した恐",
			max:      5,
			expected: "直立…",
		},
		{
			name:     "Multiline truncated",
			input:    "hi\nall",
			max:      5,
			expected: "hi…",
		},
		{
			name:     "Multiline with exact max width",
			input:    "hello\nall",
			max:      5,
			expected: "hell…",
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
