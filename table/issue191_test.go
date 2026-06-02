package table

import (
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/stretchr/testify/assert"
)

// TestIssue191_HeaderStyleBorderDoesNotBreakLayout verifies that border
// settings passed to HeaderStyle do not affect the structural layout of the
// table.  Top/bottom border LINES are rendered as separate plain strings;
// allowing them to also appear as cell borders adds extra blank rows and
// breaks column-width accounting.
func TestIssue191_HeaderStyleBorderDoesNotBreakLayout(t *testing.T) {
	cols := []Column{
		NewColumn("a", "ColA", 6),
		NewColumn("b", "ColB", 6),
	}
	rows := []Row{
		NewRow(RowData{"a": "val1", "b": "val2"}),
	}

	baseline := New(cols).WithRows(rows).View()

	cases := []struct {
		name  string
		style lipgloss.Style
	}{
		{
			name:  "HiddenBorder BorderTop true",
			style: lipgloss.NewStyle().BorderStyle(lipgloss.HiddenBorder()).BorderTop(true),
		},
		{
			name:  "HiddenBorder BorderBottom true",
			style: lipgloss.NewStyle().BorderStyle(lipgloss.HiddenBorder()).BorderBottom(true),
		},
		{
			name: "HiddenBorder all sides true",
			style: lipgloss.NewStyle().BorderStyle(lipgloss.HiddenBorder()).
				BorderTop(true).BorderBottom(true).BorderLeft(true).BorderRight(true),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := New(cols).WithRows(rows).HeaderStyle(tc.style).View()

			assert.Equal(t, lipgloss.Height(baseline), lipgloss.Height(got),
				"HeaderStyle border settings changed table height\nbaseline:\n%s\ngot:\n%s", baseline, got)
			assert.Equal(t, lipgloss.Width(baseline), lipgloss.Width(got),
				"HeaderStyle border settings changed table width\nbaseline:\n%s\ngot:\n%s", baseline, got)
		})
	}
}
