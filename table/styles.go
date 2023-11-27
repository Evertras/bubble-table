package table

import "github.com/charmbracelet/lipgloss"

// forceInheritStyles merges the given styles into the base style, and then
// applies the padding and margin from the last style in the list.
// The Inherit() method skips padding and margin, so this is a workaround.
func forceInheritStyles(base lipgloss.Style, styles ...lipgloss.Style) lipgloss.Style {
	merged := base.Copy()

	for _, style := range styles {
		merged = merged.Inherit(style).
			Padding(style.GetPadding()).
			Margin(style.GetMargin())
	}

	return merged
}
