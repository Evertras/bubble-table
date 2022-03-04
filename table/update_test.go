package table

import (
	"testing"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestUnfocusedDoesntMove(t *testing.T) {
	cols := []Column{
		NewColumn("id", "ID", 3),
	}

	model := New(cols).WithRows([]Row{
		NewRow(RowData{
			"id": "first",
		}),
		NewRow(RowData{
			"id": "second",
		}),
	})

	model, _ = model.Update(tea.KeyMsg{
		Type: tea.KeyUp,
	})

	highlighted := model.HighlightedRow()

	id, ok := highlighted.Data["id"].(string)

	assert.True(t, ok, "Failed to convert to string")

	assert.Equal(t, "first", id, "Should still be on first row")
}

func TestFocusedMovesWhenArrowsPressed(t *testing.T) {
	cols := []Column{
		NewColumn("id", "ID", 3),
	}

	model := New(cols).WithRows([]Row{
		NewRow(RowData{
			"id": "first",
		}),
		NewRow(RowData{
			"id": "second",
		}),
		NewRow(RowData{
			"id": "third",
		}),
	}).Focused(true).WithPageSize(2)

	// Note that this is assuming default keymap
	keyUp := tea.KeyMsg{Type: tea.KeyUp}
	keyDown := tea.KeyMsg{Type: tea.KeyDown}
	keyLeft := tea.KeyMsg{Type: tea.KeyLeft}
	keyRight := tea.KeyMsg{Type: tea.KeyRight}

	curID := func() string {
		str, ok := model.HighlightedRow().Data["id"].(string)

		assert.True(t, ok, "Failed to convert to string")

		return str
	}

	assert.Equal(t, "first", curID(), "Should start on first row")

	model, _ = model.Update(keyDown)
	assert.Equal(t, "second", curID(), "Default key down should move down a row")

	model, _ = model.Update(keyUp)
	assert.Equal(t, "first", curID(), "Should move back up")

	model, _ = model.Update(keyUp)
	assert.Equal(t, "third", curID(), "Moving up from top should wrap to bottom")

	model, _ = model.Update(keyDown)
	assert.Equal(t, "first", curID(), "Moving down from bottom should wrap to top")

	model, _ = model.Update(keyRight)
	assert.Equal(t, "third", curID(), "Moving right should move to second page")

	model, _ = model.Update(keyRight)
	assert.Equal(t, "first", curID(), "Moving right again should move to first page")

	model, _ = model.Update(keyLeft)
	assert.Equal(t, "third", curID(), "Moving left should move to last page")

	model, _ = model.Update(keyLeft)
	assert.Equal(t, "first", curID(), "Moving left should move back to first page")
}

func TestFocusedMovesWithCustomKeyMap(t *testing.T) {
	cols := []Column{
		NewColumn("id", "ID", 3),
	}

	customKeys := KeyMap{
		RowUp:   key.NewBinding(key.WithKeys("ctrl+a")),
		RowDown: key.NewBinding(key.WithKeys("ctrl+b")),

		RowSelectToggle: key.NewBinding(key.WithKeys("ctrl+c")),
	}

	model := New(cols).WithRows([]Row{
		NewRow(RowData{
			"id": "first",
		}),
		NewRow(RowData{
			"id": "second",
		}),
		NewRow(RowData{
			"id": "third",
		}),
	}).Focused(true).WithKeyMap(customKeys)

	keyUp := tea.KeyMsg{Type: tea.KeyUp}
	keyDown := tea.KeyMsg{Type: tea.KeyDown}
	keyCtrlA := tea.KeyMsg{Type: tea.KeyCtrlA}
	keyCtrlB := tea.KeyMsg{Type: tea.KeyCtrlB}

	assert.Equal(t, "ctrl+a", keyCtrlA.String(), "Test sanity check failed for ctrl+a")
	assert.Equal(t, "ctrl+b", keyCtrlB.String(), "Test sanity check failed for ctrl+b")

	curID := func() string {
		str, ok := model.HighlightedRow().Data["id"].(string)

		assert.True(t, ok, "Failed to convert to string")

		return str
	}

	assert.Equal(t, "first", curID(), "Should start on first row")

	model, _ = model.Update(keyDown)
	assert.Equal(t, "first", curID(), "Down arrow should do nothing")

	model, _ = model.Update(keyCtrlB)
	assert.Equal(t, "second", curID(), "Custom key map for down failed")

	model, _ = model.Update(keyUp)
	assert.Equal(t, "second", curID(), "Up arrow should do nothing")

	model, _ = model.Update(keyCtrlA)
	assert.Equal(t, "first", curID(), "Custom key map for up failed")
}

func TestSelectingRowWhenTableUnselectableDoesNothing(t *testing.T) {
	cols := []Column{
		NewColumn("id", "ID", 3),
	}

	model := New(cols).WithRows([]Row{
		NewRow(RowData{
			"id": "first",
		}),
		NewRow(RowData{
			"id": "second",
		}),
		NewRow(RowData{
			"id": "third",
		}),
	}).Focused(true)

	assert.False(t, model.GetVisibleRows()[0].selected, "Row shouldn't be selected to start")

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})

	assert.False(t, model.GetVisibleRows()[0].selected, "Row shouldn't be selected after key press")
}

func TestSelectingRowToggles(t *testing.T) {
	cols := []Column{
		NewColumn("id", "ID", 3),
	}

	model := New(cols).WithRows([]Row{
		NewRow(RowData{
			"id": "first",
		}),
		NewRow(RowData{
			"id": "second",
		}),
		NewRow(RowData{
			"id": "third",
		}),
	}).Focused(true).SelectableRows(true)

	keyEnter := tea.KeyMsg{Type: tea.KeyEnter}
	keyDown := tea.KeyMsg{Type: tea.KeyDown}

	assert.False(t, model.GetVisibleRows()[0].selected, "Row shouldn't be selected to start")
	assert.Len(t, model.SelectedRows(), 0)

	model, _ = model.Update(keyEnter)
	assert.True(t, model.GetVisibleRows()[0].selected, "Row should be selected after first toggle")
	assert.Len(t, model.SelectedRows(), 1)

	model, _ = model.Update(keyEnter)
	assert.False(t, model.GetVisibleRows()[0].selected, "Row should not be selected after second toggle")
	assert.Len(t, model.SelectedRows(), 0)

	model, _ = model.Update(keyDown)
	model, _ = model.Update(keyEnter)
	assert.True(t, model.GetVisibleRows()[1].selected, "Second row should be selected after moving and toggling")
}

func TestFilterWithKeypresses(t *testing.T) {
	cols := []Column{
		NewColumn("name", "Name", 10).WithFiltered(true),
	}

	model := New(cols).WithRows([]Row{
		NewRow(RowData{"name": "Pikachu"}),
		NewRow(RowData{"name": "Charmander"}),
	}).Focused(true).Filtered(true)

	hitKey := func(key rune) {
		model, _ = model.Update(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{key},
		})
	}

	hitEnter := func() {
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	}

	visible := model.GetVisibleRows()

	assert.Len(t, visible, 2)
	hitKey(rune(model.KeyMap().Filter.Keys()[0][0]))
	assert.Len(t, visible, 2)
	hitKey('p')
	hitKey('i')
	hitKey('k')

	visible = model.GetVisibleRows()

	assert.Len(t, visible, 1)

	hitEnter()

	hitKey('x')

	assert.Len(t, visible, 1)
}
