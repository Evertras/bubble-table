package table

import "fmt"

func (m Model) hasFooter() bool {
	return m.staticFooter != "" || m.pageSize != 0
}

func (m Model) renderFooter() string {
	if !m.hasFooter() {
		return ""
	}

	var footerText string

	switch {
	case m.staticFooter != "":
		footerText = m.staticFooter

	case m.pageSize != 0:
		footerText = fmt.Sprintf("%d/%d", m.CurrentPage(), m.MaxPages())
	}

	return m.border.styleFooter.Width(m.totalWidth).Render(footerText)
}
