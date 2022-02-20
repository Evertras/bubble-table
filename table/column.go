package table

import (
	"github.com/charmbracelet/lipgloss"
)

// Column is a column in the table.
type Column struct {
	Title string
	Key   string
	Width int

	style lipgloss.Style
}

// NewColumn creates a new column with the given information.
func NewColumn(key, title string, width int) Column {
	return Column{
		Key:   key,
		Title: title,
		Width: width,
		style: lipgloss.NewStyle().Width(width).Align(lipgloss.Right),
	}
}

// WithStyle applies a style to the column as a whole.
func (c Column) WithStyle(style lipgloss.Style) Column {
	//c.style = lipgloss.NewStyle().Width(c.Width).Align(lipgloss.Right).Inherit(style)
	c.style = style.Copy().Inherit(lipgloss.NewStyle().Width(c.Width).Align(lipgloss.Right))

	return c
}
