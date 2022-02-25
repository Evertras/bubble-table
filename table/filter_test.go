package table

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsRowMatched(t *testing.T) {
	columns := []Column{
		NewColumn("title", "title", 10).WithFiltered(true),
		NewColumn("description", "description", 10)}

	assert.True(t, isRowMatched(columns,
		NewRow(map[string]interface{}{
			"title":       "AAA",
			"description": "",
		}), ""))

	assert.True(t, isRowMatched(columns,
		NewRow(map[string]interface{}{
			"title":       "AAA",
			"description": "",
		}), "AA"))

	assert.True(t, isRowMatched(columns,
		NewRow(map[string]interface{}{
			"title":       "AAA",
			"description": "",
		}), "A"))

	assert.True(t, isRowMatched(columns,
		NewRow(map[string]interface{}{
			"title":       "AAA",
			"description": "",
		}), "a"))

	assert.False(t, isRowMatched(columns,
		NewRow(map[string]interface{}{
			"title":       "AAA",
			"description": "",
		}), "B"))

	assert.False(t, isRowMatched(columns,
		NewRow(map[string]interface{}{
			"title":       "AAA",
			"description": "BBB",
		}), "BBB"))
}
