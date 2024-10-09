package table

import (
	"testing"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/stretchr/testify/assert"
)

func TestKeyMapShortHelp(t *testing.T) {

	columns := []Column{NewColumn("c1", "Column1", 10)}
	model := New(columns)
	km := DefaultKeyMap()
	model.WithKeyMap(km)
	assert.Nil(t, model.AdditionalShortHelpKeys)
	assert.Equal(t, model.ShortHelp(), []key.Binding{
		model.keyMap.RowDown, model.keyMap.RowUp, model.keyMap.RowSelectToggle, model.keyMap.PageDown, model.keyMap.PageUp, model.keyMap.Filter, model.keyMap.FilterBlur, model.keyMap.FilterClear})

	// Testing if the 'adding of keys' works too.
	model.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "Testing additional keybinds")),
		}
	}
	assert.NotNil(t, model.AdditionalShortHelpKeys)
	assert.Equal(t, model.ShortHelp(), []key.Binding{
		model.keyMap.RowDown, model.keyMap.RowUp, model.keyMap.RowSelectToggle, model.keyMap.PageDown, model.keyMap.PageUp, model.keyMap.Filter, model.keyMap.FilterBlur, model.keyMap.FilterClear, key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "Testing additional keybinds"))})

}

func TestKeyMapFullHelp(t *testing.T) {

	columns := []Column{NewColumn("c1", "Column1", 10)}
	model := New(columns)
	km := DefaultKeyMap()
	model.WithKeyMap(km)
	assert.Nil(t, model.AdditionalFullHelpKeys)
	assert.Equal(t,
		model.FullHelp(),
		[][]key.Binding{
			{model.keyMap.RowDown, model.keyMap.RowUp, model.keyMap.RowSelectToggle},
			{model.keyMap.PageDown, model.keyMap.PageUp, model.keyMap.PageFirst, model.keyMap.PageLast},
			{model.keyMap.Filter, model.keyMap.FilterBlur, model.keyMap.FilterClear, model.keyMap.ScrollRight, model.keyMap.ScrollLeft}},
	)
	// Testing if the 'adding of keys' works too.
	model.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "Testing additional keybinds")),
		}
	}
	assert.NotNil(t, model.AdditionalFullHelpKeys)
	assert.Equal(t,
		model.FullHelp(),
		[][]key.Binding{
			{model.keyMap.RowDown, model.keyMap.RowUp, model.keyMap.RowSelectToggle},
			{model.keyMap.PageDown, model.keyMap.PageUp, model.keyMap.PageFirst, model.keyMap.PageLast},
			{model.keyMap.Filter, model.keyMap.FilterBlur, model.keyMap.FilterClear, model.keyMap.ScrollRight, model.keyMap.ScrollLeft},
			{key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "Testing additional keybinds"))}},
	)
}

// Testing if Model actually implements the 'help.KeyMap' interface.
func TestKeyMapInterface(t *testing.T) {
	model := New(nil)
	assert.Implements(t, (*help.KeyMap)(nil), model)
}
