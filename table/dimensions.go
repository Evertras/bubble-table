package table

func (m *Model) recalculateWidth() {
	total := 0

	for _, column := range m.columns {
		total += column.width
	}

	m.totalWidth = total + len(m.columns) - 1
}
