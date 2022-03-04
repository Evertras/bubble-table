package table

import (
	"fmt"
	"strings"
)

func (m Model) hasFooter() bool {
	return m.staticFooter != "" || m.pageSize != 0 || m.filtered
}

func (m Model) renderFooter() string {
	if !m.hasFooter() {
		return ""
	}

	if m.staticFooter != "" {
		return m.border.styleFooter.Width(m.totalWidth).Render(m.staticFooter)
	}

	sections := []string{}

	if m.filtered && (m.filterTextInput.Focused() || m.filterTextInput.Value() != "") {
		sections = append(sections, fmt.Sprintf("/%s", m.filterTextInput.View()))
	}

	// paged feature enabled
	if m.pageSize != 0 {
		sections = append(sections, fmt.Sprintf("%d/%d", m.CurrentPage(), m.MaxPages()))
	}

	footerText := strings.Join(sections, " ")

	styleFooter := m.baseStyle.Copy().Inherit(m.border.styleFooter).Width(m.totalWidth)

	return styleFooter.Render(footerText)
}
