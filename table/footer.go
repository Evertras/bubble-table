package table

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
)

func (m Model) hasFooter() bool {
	multiPageTargetHeight := m.targetHeight != 0 && m.pageStartIndices != nil && len(m.pageStartIndices) > 1

	return m.footerVisible &&
		(m.staticFooter != "" || m.pageSize != 0 || multiPageTargetHeight || m.filtered || m.searchable)
}

func footerInputSection(enabled bool, input textinput.Model) string {
	if !enabled {
		return ""
	}

	if input.Focused() {
		return input.View()
	}

	if input.Value() != "" {
		// When not focused, show only the prompt + value without the
		// cursor character that bubbles v2 renders in its View().
		return input.Prompt + input.Value()
	}

	return ""
}

func (m Model) footerPageSection() string {
	isPaged := m.pageSize != 0 || (m.targetHeight != 0 && m.pageStartIndices != nil && len(m.pageStartIndices) > 1)
	if !isPaged {
		return ""
	}

	str := fmt.Sprintf("%d/%d", m.CurrentPage(), m.MaxPages())

	if (m.filtered && m.filterTextInput.Focused()) || (m.searchable && m.searchTextInput.Focused()) {
		// Need to apply inline style here in case of filter input cursor, because
		// the input cursor resets the style after rendering.  Note that Inline(true)
		// creates a copy, so it's safe to use here without mutating the underlying
		// base style.
		str = m.baseStyle.Inline(true).Render(str)
	}

	return str
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

	if section := footerInputSection(m.filtered, m.filterTextInput); section != "" {
		sections = append(sections, section)
	}

	if section := footerInputSection(m.searchable, m.searchTextInput); section != "" {
		sections = append(sections, section)
	}

	if section := m.footerPageSection(); section != "" {
		sections = append(sections, section)
	}

	return styleFooter.Render(strings.Join(sections, " "))
}
