package table

import (
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/stretchr/testify/assert"
)

func TestModelInitReturnsNil(t *testing.T) {
	model := New(nil)

	cmd := model.Init()

	assert.Nil(t, cmd)
}

func TestGetVisibleRows(t *testing.T) {
	input := textinput.Model{}
	input.SetValue("AAA")
	columns := []Column{NewColumn("title", "title", 10).WithFiltered(true)}
	rows := []Row{
		NewRow(map[string]interface{}{
			"title":       "AAA",
			"description": "",
		}),
		NewRow(map[string]interface{}{
			"title":       "BBB",
			"description": "",
		}),
		NewRow(map[string]interface{}{
			"title":       "CCC",
			"description": "",
		}),
	}
	m := Model{filtered: true, filterTextInput: input, columns: columns, rows: rows}
	visibleRows := m.GetVisibleRows()
	assert.Len(t, visibleRows, 1)
}
