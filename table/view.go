package table

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the table.  It does not end in a newline, so that it can be
// composed with other elements more consistently.
func (m Model) View() string {
	numColumns := len(m.columns)

	// Safety valve for empty tables
	if numColumns == 0 {
		return ""
	}

	body := strings.Builder{}

	headerStrings := []string{}

	var (
		headerStyleLeft  lipgloss.Style
		headerStyleInner lipgloss.Style
		headerStyleRight lipgloss.Style
	)

	headerStyleLeft, headerStyleInner, headerStyleRight = m.styleHeaders()

	for columnIndex, column := range m.columns {
		headerSection := fmt.Sprintf(column.fmtString, column.Title)
		var borderStyle lipgloss.Style

		if columnIndex == 0 {
			borderStyle = headerStyleLeft
		} else if columnIndex < len(m.columns)-1 {
			borderStyle = headerStyleInner
		} else {
			borderStyle = headerStyleRight
		}

		headerStrings = append(headerStrings, borderStyle.Render(headerSection))
	}

	headerBlock := lipgloss.JoinHorizontal(lipgloss.Bottom, headerStrings...)

	rowStrs := []string{headerBlock}
	for i := range m.rows {
		rowStrs = append(rowStrs, m.renderRow(i))
	}

	body.WriteString(lipgloss.JoinVertical(lipgloss.Left, rowStrs...))

	return body.String()
}
