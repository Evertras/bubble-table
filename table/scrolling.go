package table

func (m *Model) scrollRight() {
	maxCol := len(m.columns) - 1 - m.horizontalScrollFreezeColumnsCount
	if m.horizontalScrollOffsetCol < maxCol {
		m.horizontalScrollOffsetCol++
	}
}

func (m *Model) scrollLeft() {
	if m.horizontalScrollOffsetCol > 0 {
		m.horizontalScrollOffsetCol--
	}
}
