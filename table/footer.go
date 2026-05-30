package table

import (
	"fmt"
	"strings"
)

func (m Model) hasFooter() bool {
	return m.footerVisible && (m.staticFooter != "" || m.pageSize != 0 || m.filtered)
}

func (m Model) footerFilterSection() string {
	if !m.filtered {
		return ""
	}

	if m.filterTextInput.Focused() {
		return m.filterTextInput.View()
	}

	if m.filterTextInput.Value() != "" {
		// When not focused, show only the prompt + value without the
		// cursor character that bubbles v2 renders in its View().
		return m.filterTextInput.Prompt + m.filterTextInput.Value()
	}

	return ""
}

func (m Model) renderFooter(width int, includeTop bool) string {
	if !m.hasFooter() {
		return ""
	}

	// In lipgloss v2, Width(n) sets the *total* outer width including borders.
	// The footer style has left and right borders (overhead = 2) plus a bottom
	// border.  We pass the full row width directly so that content gets
	// width-2 chars, matching the inner width of the data rows.
	styleFooter := m.baseStyle.Inherit(m.border.styleFooter).Width(width)

	if includeTop {
		styleFooter = styleFooter.BorderTop(true)
	}

	if m.staticFooter != "" {
		return styleFooter.Render(m.staticFooter)
	}

	sections := []string{}

	if section := m.footerFilterSection(); section != "" {
		sections = append(sections, section)
	}

	// paged feature enabled
	if m.pageSize != 0 {
		str := fmt.Sprintf("%d/%d", m.CurrentPage(), m.MaxPages())
		if m.filtered && m.filterTextInput.Focused() {
			// Need to apply inline style here in case of filter input cursor, because
			// the input cursor resets the style after rendering.  Note that Inline(true)
			// creates a copy, so it's safe to use here without mutating the underlying
			// base style.
			str = m.baseStyle.Inline(true).Render(str)
		}
		sections = append(sections, str)
	}

	footerText := strings.Join(sections, " ")

	return styleFooter.Render(footerText)
}
