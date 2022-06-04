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

	rowStrs := []string{m.renderHeaders()}

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
