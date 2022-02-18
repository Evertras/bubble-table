package table

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicTableShowsAllHeaders(t *testing.T) {
	const (
		firstKey   = "first-key"
		firstTitle = "First Title"
		firstWidth = 10

		secondKey   = "second-key"
		secondTitle = "Second Title"
		secondWidth = 20
	)

	columns := []Column{
		NewColumn(firstKey, firstTitle, firstWidth),
		NewColumn(secondKey, secondTitle, secondWidth),
	}

	table := New(columns)

	rendered := table.View()

	assert.Contains(t, rendered, firstTitle)
	assert.Contains(t, rendered, secondTitle)

	assert.False(t, strings.HasSuffix(rendered, "\n"), "Should not end in newline")
}

func TestNilColumnsSafelyReturnsEmptyView(t *testing.T) {
	table := New(nil)

	assert.Equal(t, "", table.View())
}
