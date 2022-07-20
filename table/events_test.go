package table

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestUserEventsEmptyWhenNothingHappens(t *testing.T) {
	model := New([]Column{})

	events := model.GetLastUpdateUserEvents()

	assert.Len(t, events, 0, "Should be empty when nothing has happened")

	model, _ = model.Update(nil)

	events = model.GetLastUpdateUserEvents()

	assert.Len(t, events, 0, "Should be empty when no changes made in Update")
}

// nolint: funlen // This is a bunch of checks in a row, this is fine
func TestUserEventHighlightedIndexChanged(t *testing.T) {
	// Don't need any actual row data for this
	empty := RowData{}

	model := New([]Column{}).
		Focused(true).
		WithRows(
			[]Row{
				NewRow(empty),
				NewRow(empty),
				NewRow(empty),
				NewRow(empty),
			},
		)

	hitDown := func() {
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	}

	hitUp := func() {
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	}

	checkEvent := func(events []UserEvent, expectedPreviousIndex, expectedCurrentIndex int) {
		if len(events) != 1 {
			assert.FailNow(t, "Asked to check events with len of not 1, test is bad")
		}

		switch event := events[0].(type) {
		case UserEventHighlightedIndexChanged:
			assert.Equal(t, expectedPreviousIndex, event.PreviousRowIndex)
			assert.Equal(t, expectedCurrentIndex, event.SelectedRowIndex)

		default:
			assert.Failf(t, "Event is not expected type UserEventHighlightedIndexChanged", "%+v", event)
		}
	}

	events := model.GetLastUpdateUserEvents()
	assert.Len(t, events, 0, "Should be empty when nothing has happened")

	// Hit down to change row down by one
	hitDown()
	events = model.GetLastUpdateUserEvents()
	assert.Len(t, events, 1, "Missing event for scrolling down")
	checkEvent(events, 0, 1)

	// Do some no-op
	model, _ = model.Update(nil)
	events = model.GetLastUpdateUserEvents()
	assert.Len(t, events, 0, "Events not cleared between Updates")

	// Hit up to go back to top
	hitUp()
	events = model.GetLastUpdateUserEvents()
	assert.Len(t, events, 1, "Missing event to scroll back up")
	checkEvent(events, 1, 0)

	// Hit up to scroll around to bottom
	hitUp()
	events = model.GetLastUpdateUserEvents()
	assert.Len(t, events, 1, "Missing event to scroll up with wrap")
	checkEvent(events, 0, 3)

	// Now check to make sure it doesn't trigger when row index doesn't change
	model = model.WithRows([]Row{NewRow(empty)})
	hitDown()
	events = model.GetLastUpdateUserEvents()
	assert.Len(t, events, 0, "There's no row to change to for single row table, event shouldn't exist")

	model = model.WithRows([]Row{})
	hitDown()
	events = model.GetLastUpdateUserEvents()
	assert.Len(t, events, 0, "There's no row to change to for an empty table, event shouldn't exist")
}

// nolint: funlen // This is a bunch of checks in a row, this is fine
func TestUserEventRowSelectToggled(t *testing.T) {
	// Don't need any actual row data for this
	empty := RowData{}

	model := New([]Column{}).
		Focused(true).
		WithRows(
			[]Row{
				NewRow(empty),
				NewRow(empty),
				NewRow(empty),
				NewRow(empty),
			},
		).
		SelectableRows(true)

	hitDown := func() {
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	}

	hitSelectToggle := func() {
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeySpace})
	}

	checkEvent := func(events []UserEvent, expectedRowIndex int, expectedSelectionState bool) {
		if len(events) != 1 {
			assert.FailNow(t, "Asked to check events with len of not 1, test is bad")
		}

		switch event := events[0].(type) {
		case UserEventRowSelectToggled:
			assert.Equal(t, expectedRowIndex, event.RowIndex, "Row index wrong")
			assert.Equal(t, expectedSelectionState, event.IsSelected, "Selection state wrong")

		default:
			assert.Failf(t, "Event is not expected type UserEventRowSelectToggled", "%+v", event)
		}
	}

	events := model.GetLastUpdateUserEvents()
	assert.Len(t, events, 0, "Should be empty when nothing has happened")

	// Try initial selection
	hitSelectToggle()
	events = model.GetLastUpdateUserEvents()
	assert.Len(t, events, 1, "Missing event for selection toggle")
	checkEvent(events, 0, true)

	// Do some no-op
	model, _ = model.Update(nil)
	events = model.GetLastUpdateUserEvents()
	assert.Len(t, events, 0, "Events not cleared between Updates")

	// Check deselection
	hitSelectToggle()
	events = model.GetLastUpdateUserEvents()
	assert.Len(t, events, 1, "Missing event to toggle select for second time")
	checkEvent(events, 0, false)

	// Try one row down... note that the row change event should clear after the
	// first keypress
	hitDown()
	hitSelectToggle()
	events = model.GetLastUpdateUserEvents()
	assert.Len(t, events, 1, "Missing event after scrolling down")
	checkEvent(events, 1, true)

	// Check edge case of empty table
	model = model.WithRows([]Row{})
	hitSelectToggle()
	events = model.GetLastUpdateUserEvents()
	assert.Len(t, events, 0, "There's no row to select for an empty table, event shouldn't exist")
}

func TestFilterFocusEvents(t *testing.T) {
	model := New([]Column{}).Filtered(true).Focused(true)

	events := model.GetLastUpdateUserEvents()

	assert.Empty(t, events, "Unexpected events to start")

	// Start filter
	model, _ = model.Update(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'/'},
	})
	events = model.GetLastUpdateUserEvents()
	assert.Len(t, events, 1, "Only expected one event")
	switch events[0].(type) {
	case UserEventFilterInputFocused:
	default:
		assert.FailNow(t, "Unexpected event type")
	}

	// Stop filter
	model, _ = model.Update(tea.KeyMsg{
		Type: tea.KeyEnter,
	})
	events = model.GetLastUpdateUserEvents()
	assert.Len(t, events, 1, "Only expected one event")
	switch events[0].(type) {
	case UserEventFilterInputUnfocused:
	default:
		assert.FailNow(t, "Unexpected event type")
	}
}
