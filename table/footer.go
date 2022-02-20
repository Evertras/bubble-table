package table

func (m Model) hasFooter() bool {
	return m.staticFooter != ""
}

func (m Model) renderFooter() string {
	if !m.hasFooter() {
		return ""
	}

	return m.border.styleFooter.Width(m.totalWidth).Render(m.staticFooter)
}
