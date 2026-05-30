package table

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// This is long and could use some refactoring in the future, but unsure of how
// to pick it apart right now.
//
//nolint:funlen,cyclop
func (m Model) renderHeaders() string {
	headerStrings := []string{}

	// Track column widths (content-only, not including border chars) so we can
	// build the top border line and separator line as plain strings.
	renderedColWidths := []int{}

	totalRenderedWidth := 0

	headerStyles := m.styleHeaders()

	renderHeader := func(column Column, borderStyle lipgloss.Style) string {
		// In lipgloss v2, Width(n) sets the *total* outer width including borders.
		// Compute the actual border char overhead using the border size accessors so
		// that any Width already set on borderStyle (e.g. from genOverflowStyle) does
		// not pollute the overhead calculation.
		borderOverhead := borderStyle.GetBorderLeftSize() + borderStyle.GetBorderRightSize()
		style := borderStyle.Inherit(column.style).Inherit(m.baseStyle).Width(column.width + borderOverhead)

		headerSection := limitStr(column.title, column.width)

		return style.Render(headerSection)
	}

	for columnIndex, column := range m.columns {
		var borderStyle lipgloss.Style

		if m.horizontalScrollOffsetCol > 0 && columnIndex == m.horizontalScrollFreezeColumnsCount {
			if columnIndex == 0 {
				borderStyle = headerStyles.left
			} else {
				borderStyle = headerStyles.inner
			}

			overflowCol := genOverflowColumnLeft(1)
			rendered := renderHeader(overflowCol, borderStyle)

			totalRenderedWidth += lipgloss.Width(rendered)

			headerStrings = append(headerStrings, rendered)
			renderedColWidths = append(renderedColWidths, overflowCol.width)
		}

		if columnIndex >= m.horizontalScrollFreezeColumnsCount &&
			columnIndex < m.horizontalScrollOffsetCol+m.horizontalScrollFreezeColumnsCount {
			continue
		}

		if len(headerStrings) == 0 {
			borderStyle = headerStyles.left
		} else if columnIndex < len(m.columns)-1 {
			borderStyle = headerStyles.inner
		} else {
			borderStyle = headerStyles.right
		}

		rendered := renderHeader(column, borderStyle)

		if m.maxTotalWidth != 0 {
			renderedWidth := lipgloss.Width(rendered)

			const (
				borderAdjustment = 1
				overflowColWidth = 2
			)

			targetWidth := m.maxTotalWidth - overflowColWidth

			if columnIndex == len(m.columns)-1 {
				// If this is the last header, we don't need to account for the
				// overflow arrow column
				targetWidth = m.maxTotalWidth
			}

			if totalRenderedWidth+renderedWidth > targetWidth {
				overflowWidth := m.maxTotalWidth - totalRenderedWidth - borderAdjustment
				overflowStyle := genOverflowStyle(headerStyles.right, overflowWidth)
				overflowColumn := genOverflowColumnRight(overflowWidth)

				overflowStr := renderHeader(overflowColumn, overflowStyle)

				headerStrings = append(headerStrings, overflowStr)
				renderedColWidths = append(renderedColWidths, overflowColumn.width)

				break
			}

			totalRenderedWidth += renderedWidth
		}

		headerStrings = append(headerStrings, rendered)
		renderedColWidths = append(renderedColWidths, column.width)
	}

	headerBlock := lipgloss.JoinHorizontal(lipgloss.Bottom, headerStrings...)

	// Build the top border line as a plain string using the recorded content widths.
	topLine := m.border.buildTopBorderLine(renderedColWidths)

	// Determine the last line of the header block:
	//  - If data rows follow, use a separator (junction characters).
	//  - If no rows follow and no footer, use the bottom border (corner characters).
	//  - If no rows follow but a footer follows, use a separator so the footer
	//    can attach cleanly.
	hasRows := len(m.GetVisibleRows()) > 0 || m.calculatePadding(0) > 0

	var lastLine string

	if hasRows {
		lastLine = m.border.buildSeparatorLine(renderedColWidths)
	} else {
		lastLine = m.border.buildBottomBorderLine(renderedColWidths, m.hasFooter())
	}

	return strings.Join([]string{topLine, headerBlock, lastLine}, "\n")
}
