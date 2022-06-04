package table

import "github.com/charmbracelet/lipgloss"

func (m Model) renderHeaders() string {
	headerStrings := []string{}

	totalRenderedWidth := 0

	headerStyles := m.styleHeaders()

	renderHeader := func(column Column, borderStyle lipgloss.Style) string {
		borderStyle = borderStyle.Inherit(column.style).Inherit(m.baseStyle)

		headerSection := limitStr(column.title, column.width)

		return borderStyle.Render(headerSection)
	}

	if m.horizontalScrollOffsetCol > 0 {
		borderStyle := headerStyles.left.Copy()

		rendered := renderHeader(genOverflowColumnLeft(1), borderStyle)

		totalRenderedWidth += lipgloss.Width(rendered)

		headerStrings = append(headerStrings, rendered)
	}

	for columnIndex, column := range m.columns {
		if columnIndex < m.horizontalScrollOffsetCol {
			continue
		}

		var borderStyle lipgloss.Style

		if len(headerStrings) == 0 {
			borderStyle = headerStyles.left.Copy()
		} else if columnIndex < len(m.columns)-1 {
			borderStyle = headerStyles.inner.Copy()
		} else {
			borderStyle = headerStyles.right.Copy()
		}

		rendered := renderHeader(column, borderStyle)

		if m.maxTotalWidth != 0 {
			renderedWidth := lipgloss.Width(rendered)

			const borderAdjustment = 1

			if totalRenderedWidth+renderedWidth > m.maxTotalWidth-borderAdjustment*2 {
				overflowWidth := m.maxTotalWidth - totalRenderedWidth - borderAdjustment
				overflowStyle := genOverflowStyle(headerStyles.right, overflowWidth)
				overflowColumn := genOverflowColumnRight(overflowWidth)

				overflowStr := renderHeader(overflowColumn, overflowStyle)

				headerStrings = append(headerStrings, overflowStr)

				break
			}

			totalRenderedWidth += renderedWidth
		}

		headerStrings = append(headerStrings, rendered)
	}

	headerBlock := lipgloss.JoinHorizontal(lipgloss.Bottom, headerStrings...)

	return headerBlock
}
