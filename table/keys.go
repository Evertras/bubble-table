package table

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines the keybindings for the table when it's focused.
type KeyMap struct {
	RowDown key.Binding
	RowUp   key.Binding

	RowSelectToggle key.Binding
}

// DefaultKeyMap returns a set of sensible defaults for controlling a focused table.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		RowDown: key.NewBinding(
			key.WithKeys("down", "j"),
		),
		RowUp: key.NewBinding(
			key.WithKeys("up", "k"),
		),
		RowSelectToggle: key.NewBinding(
			key.WithKeys(" ", "enter"),
		),
	}
}
