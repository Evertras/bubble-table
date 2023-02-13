package table

import (
	"testing"

	"github.com/muesli/reflow/ansi"
	"github.com/stretchr/testify/assert"
)

// This function is only long because of repetitive test definitions, this is fine
//
//nolint:funlen
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
		{
			name:     "Embedded ANSI control sequences with exact max width",
			input:    "\x1b[31;41mtest\x1b[0m",
			max:      4,
			expected: "\x1b[31;41mtest\x1b[0m",
		},
		{
			name:     "Embedded ANSI control sequences with truncation",
			input:    "\x1b[31;41mte\x1b[0m\x1b[0m\x1b[0mst",
			max:      3,
			expected: "\x1b[31;41mte\x1b[0m\x1b[0m\x1b[0m…",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := limitStr(test.input, test.max)

			assert.Equal(t, test.expected, output)
			assert.LessOrEqual(t, ansi.PrintableRuneWidth(output), test.max)
		})
	}
}
