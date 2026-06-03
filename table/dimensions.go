package table

import (
	"charm.land/lipgloss/v2"
)

func (m *Model) recalculateWidth() {
	if m.targetTotalWidth != 0 {
		m.totalWidth = m.targetTotalWidth
	} else {
		total := 0

		for _, column := range m.columns {
			total += column.width
		}

		m.totalWidth = total + len(m.columns) + 1
	}

	updateColumnWidths(m.columns, m.targetTotalWidth)

	m.recalculateLastHorizontalColumn()
}

// Updates column width in-place.  This could be optimized but should be called
// very rarely so we prioritize simplicity over performance here.
func updateColumnWidths(cols []Column, totalWidth int) {
	totalFlexWidth := totalWidth - len(cols) - 1
	totalFlexFactor := 0
	flexGCD := 0

	for _, col := range cols {
		if !col.isFlex() {
			totalFlexWidth -= col.width
			// Do not set Width on col.style here; the rendering functions apply
			// Width(column.width + borderOverhead) at render time to correctly
			// account for lipgloss v2's total-outer-width semantics.
		} else {
			totalFlexFactor += col.flexFactor
			flexGCD = gcd(flexGCD, col.flexFactor)
		}
	}

	if totalFlexFactor == 0 {
		return
	}

	// We use the GCD here because otherwise very large values won't divide
	// nicely as ints
	totalFlexFactor /= flexGCD

	flexUnit := totalFlexWidth / totalFlexFactor
	leftoverWidth := totalFlexWidth % totalFlexFactor

	for index := range cols {
		if !cols[index].isFlex() {
			continue
		}

		width := flexUnit * (cols[index].flexFactor / flexGCD)

		if leftoverWidth > 0 {
			width++
			leftoverWidth--
		}

		if index == len(cols)-1 {
			width += leftoverWidth
			leftoverWidth = 0
		}

		width = max(width, 1)

		cols[index].width = width

		// Do not set Width on the style; rendering functions apply
		// Width(column.width + borderOverhead) at render time.
	}
}

func (m *Model) recalculateHeight() {
	header := m.renderHeaders()
	headerHeight := 1 // Header always has the top border
	if m.headerVisible {
		headerHeight = lipgloss.Height(header)
	}

	footer := m.renderFooter(lipgloss.Width(header), false)
	var footerHeight int
	if footer != "" {
		footerHeight = lipgloss.Height(footer)
	}

	m.metaHeight = headerHeight + footerHeight
}

// ensurePageMap rebuilds the page-start index cache if it is stale.
// For targetHeight mode it uses a two-pass approach: first calculate without a
// footer to determine whether multiple pages are needed, then if they are,
// recalculate with the footer included in the height budget and rebuild.
func (m *Model) ensurePageMap() {
	if m.pageStartIndices == nil {
		m.recalculateHeight()
		m.buildPageStartIndices()

		if m.targetHeight != 0 && len(m.pageStartIndices) > 1 {
			// Footer is now active (multi-page); redo with footer in budget.
			m.recalculateHeight()
			m.buildPageStartIndices()
		}
	}
}

// buildPageStartIndices computes which row index each page starts on and stores
// the result in m.pageStartIndices.  Must be called after recalculateHeight so
// that m.metaHeight is up to date.
func (m *Model) buildPageStartIndices() {
	rows := m.GetVisibleRows()

	if len(rows) == 0 {
		m.pageStartIndices = []int{}

		return
	}

	// metaHeight covers the header (including top border) and footer.
	// The bottom border is appended to the last data row by assembleRowOutput,
	// so subtract 1 more to account for it.
	availableLines := m.targetHeight - m.metaHeight - 1

	if availableLines < 1 {
		availableLines = 1
	}

	pageStarts := []int{0}
	currentPageLines := 0

	for rowIdx, row := range rows {
		rowLines := m.rowLineCount(row)

		var linesNeeded int
		if currentPageLines == 0 {
			linesNeeded = rowLines
		} else if m.rowSeparator {
			linesNeeded = rowLines + 1 // separator between rows
		} else {
			linesNeeded = rowLines
		}

		if currentPageLines+linesNeeded > availableLines && currentPageLines > 0 {
			// Row doesn't fit on the current page; start a new one.
			pageStarts = append(pageStarts, rowIdx)
			currentPageLines = m.rowLineCount(row)
		} else {
			currentPageLines += linesNeeded
		}
	}

	m.pageStartIndices = pageStarts
}

func (m *Model) calculatePadding(numRows int) int {
	if m.minimumHeight == 0 {
		return 0
	}

	padding := m.minimumHeight - m.metaHeight - numRows - 1 // additional 1 for bottom border

	if padding == 0 && numRows == 0 {
		// This is an edge case where we want to add 1 additional line of height, i.e.
		// add a border without an empty row. However, this is not possible, so we need
		// to add an extra row which will result in the table being 1 row taller than
		// the requested minimum height.
		return 1
	}

	if padding < 0 {
		// Table is already larger than minimum height, do nothing.
		return 0
	}

	return padding
}
