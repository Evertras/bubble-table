package table

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// visibleDataLines returns the total terminal lines occupied by visible rows in
// the range [startRowIndex, endRowIndex], including any separator lines between
// rows when WithRowBorder is active.
func (m *Model) visibleDataLines(startRowIndex, endRowIndex int) int {
	numRows := endRowIndex - startRowIndex + 1

	if numRows <= 0 {
		return numRows
	}

	if !m.multiline {
		if m.rowSeparator {
			// Each pair of adjacent rows has one separator line between them.
			return 2*numRows - 1
		}

		return numRows
	}

	total := 0
	visibleRows := m.GetVisibleRows()

	for rowIdx := startRowIndex; rowIdx <= endRowIndex; rowIdx++ {
		total += m.rowLineCount(visibleRows[rowIdx])

		if m.rowSeparator && rowIdx > startRowIndex {
			total++
		}
	}

	return total
}

// View renders the table. It does not end in a newline, so that it can be
// composed with other elements more consistently.
//
//nolint:cyclop
func (m Model) View() string {
	// Safety valve for empty tables
	if len(m.columns) == 0 {
		return ""
	}

	body := strings.Builder{}

	rowStrs := make([]string, 0, 1)

	headers := m.renderHeaders()

	startRowIndex, endRowIndex := m.VisibleIndices()

	padding := m.calculatePadding(m.visibleDataLines(startRowIndex, endRowIndex))

	numRows := endRowIndex - startRowIndex + 1

	if m.headerVisible {
		rowStrs = append(rowStrs, headers)
	} else if numRows > 0 || padding > 0 {
		//nolint: mnd // This is just getting the first newlined substring
		split := strings.SplitN(headers, "\n", 2)
		rowStrs = append(rowStrs, split[0])
	}

	for i := startRowIndex; i <= endRowIndex; i++ {
		isLast := padding == 0 && i == endRowIndex
		separatorBelow := m.rowSeparator && i < endRowIndex
		rowStrs = append(rowStrs, m.renderRow(i, isLast, separatorBelow))
	}

	for i := 1; i <= padding; i++ {
		rowStrs = append(rowStrs, m.renderBlankRow(i == padding))
	}

	var footer string

	if len(rowStrs) > 0 {
		footer = m.renderFooter(lipgloss.Width(rowStrs[0]), false)
	} else {
		footer = m.renderFooter(lipgloss.Width(headers), true)
	}

	if footer != "" {
		rowStrs = append(rowStrs, footer)
	}

	if len(rowStrs) == 0 {
		return ""
	}

	body.WriteString(lipgloss.JoinVertical(lipgloss.Left, rowStrs...))

	return body.String()
}
