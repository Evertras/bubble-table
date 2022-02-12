package table

import (
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

	headers := []Header{
		NewHeader(firstKey, firstTitle, firstWidth),
		NewHeader(secondKey, secondTitle, secondWidth),
	}

	table := New(headers)

	rendered := table.View()

	assert.Contains(t, rendered, firstTitle)
	assert.Contains(t, rendered, secondTitle)
}
