package table

import "github.com/charmbracelet/lipgloss"

type StyledCell struct {
	Data  interface{}
	Style lipgloss.Style
}

func NewStyledCell(data interface{}, style lipgloss.Style) StyledCell {
	return StyledCell{data, style}
}
