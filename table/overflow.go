package table

import "github.com/charmbracelet/lipgloss"

func genOverflowStyle(base lipgloss.Style, width int) lipgloss.Style {
	style := lipgloss.NewStyle().Width(width).Align(lipgloss.Right)

	style.Inherit(base)

	return style
}

func genOverflowColumn(width int) Column {
	return NewColumn(columnKeyOverflow, ">", width)
}
