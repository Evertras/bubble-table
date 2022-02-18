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

	m := New(cols).WithRows([]Row{
		NewRow(RowData{
			"id": "first",
		}),
		NewRow(RowData{
			"id": "second",
		}),
	})

	m, _ = m.Update(tea.KeyMsg{
		Type: tea.KeyUp,
	})

	highlighted := m.HighlightedRow()

	assert.Equal(t, "first", highlighted.Data["id"].(string), "Should still be on first row")
}

func TestFocusedMovesWhenArrowsPressed(t *testing.T) {
	cols := []Column{
		NewColumn("id", "ID", 3),
	}

	m := New(cols).WithRows([]Row{
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

	// Note that this is assuming default keymap
	keyUp := tea.KeyMsg{Type: tea.KeyUp}
	keyDown := tea.KeyMsg{Type: tea.KeyDown}

	curID := func() string {
		return m.HighlightedRow().Data["id"].(string)
	}

	assert.Equal(t, "first", curID(), "Should start on first row")

	m, _ = m.Update(keyDown)
	assert.Equal(t, "second", curID(), "Default key down should move down a row")

	m, _ = m.Update(keyUp)
	assert.Equal(t, "first", curID(), "Should move back up")

	m, _ = m.Update(keyUp)
	assert.Equal(t, "third", curID(), "Moving up from top should wrap to bottom")

	m, _ = m.Update(keyDown)
	assert.Equal(t, "first", curID(), "Moving down from bottom should wrap to top")
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

	m := New(cols).WithRows([]Row{
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
		return m.HighlightedRow().Data["id"].(string)
	}

	assert.Equal(t, "first", curID(), "Should start on first row")

	m, _ = m.Update(keyDown)
	assert.Equal(t, "first", curID(), "Down arrow should do nothing")

	m, _ = m.Update(keyCtrlB)
	assert.Equal(t, "second", curID(), "Custom key map for down failed")

	m, _ = m.Update(keyUp)
	assert.Equal(t, "second", curID(), "Up arrow should do nothing")

	m, _ = m.Update(keyCtrlA)
	assert.Equal(t, "first", curID(), "Custom key map for up failed")
}

func TestSelectingRowWhenTableUnselectableDoesNothing(t *testing.T) {
	cols := []Column{
		NewColumn("id", "ID", 3),
	}

	m := New(cols).WithRows([]Row{
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

	assert.False(t, m.rows[0].selected, "Row shouldn't be selected to start")

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	assert.False(t, m.rows[0].selected, "Row shouldn't be selected after key press")
}

func TestSelectingRowToggles(t *testing.T) {
	cols := []Column{
		NewColumn("id", "ID", 3),
	}

	m := New(cols).WithRows([]Row{
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

	assert.False(t, m.rows[0].selected, "Row shouldn't be selected to start")

	m, _ = m.Update(keyEnter)
	assert.True(t, m.rows[0].selected, "Row should be selected after first toggle")

	m, _ = m.Update(keyEnter)
	assert.False(t, m.rows[0].selected, "Row should not be selected after second toggle")

	m, _ = m.Update(keyDown)
	m, _ = m.Update(keyEnter)
	assert.True(t, m.rows[1].selected, "Second row should be selected after moving and toggling")
}
