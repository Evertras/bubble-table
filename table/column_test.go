package table

import (
	"fmt"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
)

func TestColumnTitle(t *testing.T) {
	tests := []struct {
		title    string
		expected string
	}{
		{
			title:    "foo",
			expected: "foo",
		},
		{
			title:    "bar",
			expected: "bar",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("title %s gives %s", test.title, test.expected), func(t *testing.T) {
			col := NewColumn("key", test.title, 10)
			assert.Equal(t, test.expected, col.Title())
		})
	}
}

func TestColumnKey(t *testing.T) {
	tests := []struct {
		key      string
		expected string
	}{
		{
			key:      "foo",
			expected: "foo",
		},
		{
			key:      "bar",
			expected: "bar",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("key %s gives %s", test.key, test.expected), func(t *testing.T) {
			col := NewColumn(test.key, "title", 10)
			assert.Equal(t, test.expected, col.Key())
		})
	}
}

func TestColumnWidth(t *testing.T) {
	tests := []struct {
		width    int
		expected int
	}{
		{
			width:    3,
			expected: 3,
		},
		{
			width:    4,
			expected: 4,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("width %d gives %d", test.width, test.expected), func(t *testing.T) {
			col := NewColumn("key", "title", test.width)
			assert.Equal(t, test.expected, col.Width())
		})
	}
}

func TestColumnFlexFactor(t *testing.T) {
	tests := []struct {
		flexFactor int
		expected   int
	}{
		{
			flexFactor: 3,
			expected:   3,
		},
		{
			flexFactor: 4,
			expected:   4,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("flexFactor %d gives %d", test.flexFactor, test.expected), func(t *testing.T) {
			col := NewFlexColumn("key", "title", test.flexFactor)
			assert.Equal(t, test.expected, col.FlexFactor())
		})
	}
}

func TestColumnIsFlex(t *testing.T) {
	testsFlexColumn := []struct {
		flexFactor int
		expected   bool
	}{
		{
			flexFactor: 3,
			expected:   true,
		},
		{
			flexFactor: 0,
			expected:   true,
		},
	}

	for _, test := range testsFlexColumn {
		t.Run(fmt.Sprintf("flexFactor %d gives %t", test.flexFactor, test.expected), func(t *testing.T) {
			col := NewFlexColumn("key", "title", test.flexFactor)
			assert.Equal(t, test.expected, col.IsFlex())
		})
	}

	testsRegularColumn := []struct {
		width    int
		expected bool
	}{
		{
			width:    3,
			expected: false,
		},
		{
			width:    0,
			expected: false,
		},
	}

	for _, test := range testsRegularColumn {
		t.Run(fmt.Sprintf("width %d gives %t", test.width, test.expected), func(t *testing.T) {
			col := NewColumn("key", "title", test.width)
			assert.Equal(t, test.expected, col.IsFlex())
		})
	}
}

func TestColumnFilterable(t *testing.T) {
	tests := []struct {
		filterable bool
		expected   bool
	}{
		{
			filterable: true,
			expected:   true,
		},
		{
			filterable: false,
			expected:   false,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("filterable %t gives %t", test.filterable, test.expected), func(t *testing.T) {
			col := NewColumn("key", "title", 10)
			col = col.WithFiltered(test.filterable)
			assert.Equal(t, test.expected, col.Filterable())
		})
	}
}

func TestColumnStyle(t *testing.T) {
	width := 10
	tests := []struct {
		style    lipgloss.Style
		expected lipgloss.Style
	}{
		{
			style:    lipgloss.NewStyle(),
			expected: lipgloss.NewStyle().Width(width),
		},
		{
			style:    lipgloss.NewStyle().Bold(true),
			expected: lipgloss.NewStyle().Bold(true).Width(width),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("style %v gives %v", test.style, test.expected), func(t *testing.T) {
			col := NewColumn("key", "title", width).WithStyle(test.style)
			assert.Equal(t, test.expected, col.Style())
		})
	}
}

func TestColumnFormatString(t *testing.T) {
	tests := []struct {
		fmtString string
		expected  string
	}{
		{
			fmtString: "%v",
			expected:  "%v",
		},
		{
			fmtString: "%.2f",
			expected:  "%.2f",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("fmtString %s gives %s", test.fmtString, test.expected), func(t *testing.T) {
			col := NewColumn("key", "title", 10)
			col = col.WithFormatString(test.fmtString)
			assert.Equal(t, test.expected, col.FmtString())
		})
	}
}
