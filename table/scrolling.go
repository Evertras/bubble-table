package table

func (m *Model) scrollRight() {
	if m.horizontalScrollOffsetCol < m.maxHorizontalColumnIndex {
		m.horizontalScrollOffsetCol++
	}
}

func (m *Model) scrollLeft() {
	if m.horizontalScrollOffsetCol > 0 {
		m.horizontalScrollOffsetCol--
	}
}

func (m *Model) recalculateLastHorizontalColumn() {
	columns := make([]Column, 0, len(m.columns))
	for _, column := range m.columns {
		if column.hidden {
			continue
		}

		columns = append(columns, column)
	}

	if m.horizontalScrollFreezeColumnsCount >= len(columns) {
		m.maxHorizontalColumnIndex = 0

		return
	}

	if m.totalWidth <= m.maxTotalWidth {
		m.maxHorizontalColumnIndex = 0

		return
	}

	const (
		leftOverflowWidth = 2
		borderAdjustment  = 1
	)

	// Always have left border
	visibleWidth := borderAdjustment + leftOverflowWidth

	for i := 0; i < m.horizontalScrollFreezeColumnsCount; i++ {
		visibleWidth += columns[i].width + borderAdjustment
	}

	m.maxHorizontalColumnIndex = len(columns) - 1

	// Work backwards from the right
	for i := len(columns) - 1; i >= m.horizontalScrollFreezeColumnsCount && visibleWidth <= m.maxTotalWidth; i-- {
		visibleWidth += columns[i].width + borderAdjustment

		if visibleWidth <= m.maxTotalWidth {
			m.maxHorizontalColumnIndex = i - m.horizontalScrollFreezeColumnsCount
		}
	}
}
