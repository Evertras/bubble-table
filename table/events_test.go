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
}
