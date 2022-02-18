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
	hasRows := len(m.rows) > 0

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

	if numColumns == 1 {
		if hasRows {
			headerStyleLeft = m.border.styleSingleColumnTop
		} else {
			headerStyleLeft = m.border.styleSingleCell
		}

		headerStyleLeft = headerStyleLeft.Copy().Inherit(m.headerStyle)
	} else {
		if hasRows {
			headerStyleLeft = m.border.styleMultiTopLeft
			headerStyleInner = m.border.styleMultiTop
			headerStyleRight = m.border.styleMultiTopRight
		} else {
			headerStyleLeft = m.border.styleSingleRowLeft
			headerStyleInner = m.border.styleSingleRowInner
			headerStyleRight = m.border.styleSingleRowRight
		}

		headerStyleLeft = headerStyleLeft.Copy().Inherit(m.headerStyle)
		headerStyleInner = headerStyleInner.Copy().Inherit(m.headerStyle)
		headerStyleRight = headerStyleRight.Copy().Inherit(m.headerStyle)
	}

	for i, header := range m.columns {
		headerSection := fmt.Sprintf(header.fmtString, header.Title)
		var borderStyle lipgloss.Style

		if i == 0 {
			borderStyle = headerStyleLeft
		} else if i < len(m.columns)-1 {
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
