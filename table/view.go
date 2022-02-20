package table

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the table.  It does not end in a newline, so that it can be
// composed with other elements more consistently.
func (m Model) View() string {
	// Safety valve for empty tables
	if len(m.columns) == 0 {
		return ""
	}

	body := strings.Builder{}

	headerStrings := []string{}

	headerStyles := m.styleHeaders()

	for columnIndex, column := range m.columns {
		headerSection := limitStr(column.Title, column.Width)
		var borderStyle lipgloss.Style

		if columnIndex == 0 {
			borderStyle = headerStyles.left.Copy()
		} else if columnIndex < len(m.columns)-1 {
			borderStyle = headerStyles.inner.Copy()
		} else {
			borderStyle = headerStyles.right.Copy()
		}

		borderStyle = borderStyle.Inherit(column.style)

		headerStrings = append(headerStrings, borderStyle.Render(headerSection))
	}

	headerBlock := lipgloss.JoinHorizontal(lipgloss.Bottom, headerStrings...)

	rowStrs := []string{headerBlock}

	startRowIndex, endRowIndex := m.VisibleIndices()
	for i := startRowIndex; i <= endRowIndex; i++ {
		rowStrs = append(rowStrs, m.renderRow(i, i == endRowIndex))
	}

	footer := m.renderFooter()

	if footer != "" {
		rowStrs = append(rowStrs, footer)
	}

	body.WriteString(lipgloss.JoinVertical(lipgloss.Left, rowStrs...))

	return body.String()
}
