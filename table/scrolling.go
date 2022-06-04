package table

func (m *Model) scrollRight() {
	if m.horizontalScrollOffsetCol < len(m.columns)-1 {
		m.horizontalScrollOffsetCol++
	}
}

func (m *Model) scrollLeft() {
	if m.horizontalScrollOffsetCol > 0 {
		m.horizontalScrollOffsetCol--
	}
}
