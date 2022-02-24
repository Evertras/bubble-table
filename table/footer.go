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
	var pagination string
	var filter string

	if m.staticFooter != "" {
		footerText = m.staticFooter
	} else {
		// paged feature enabled
		if m.pageSize != 0 {
			pagination = fmt.Sprintf("%d/%d", m.CurrentPage(), m.MaxPages())
		}
		if m.filtered {
			// filter pressing
			if m.filterTextInput.Focused() {
				filter = fmt.Sprintf("/%s", m.filterTextInput.View())
			} else {
				if m.filterTextInput.Value() != "" {
					filter = fmt.Sprintf("/%s", m.filterTextInput.Value())
				}
			}
		}
		footerText = fmt.Sprintf("%s %s", filter, pagination)
	}

	return m.border.styleFooter.Width(m.totalWidth).Render(footerText)
}
