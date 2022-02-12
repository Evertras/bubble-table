package table

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type RowData map[string]interface{}

type Row struct {
	Style lipgloss.Style
	Data  RowData

	selected bool
}

func NewRow(data RowData) Row {
	d := Row{
		Data: make(map[string]interface{}),
	}

	for key, val := range data {
		// Doesn't deep copy val, but close enough for now...
		d.Data[key] = val
	}

	return d
}

func (r Row) WithStyle(style lipgloss.Style) Row {
	r.Style = style

	return r
}

var borderRowLeft = lipgloss.Border{
	Left:        "┃",
	Right:       "┃",
	Bottom:      "━",
	BottomLeft:  "┗",
	BottomRight: "┻",
}

var borderRowMiddle = lipgloss.Border{
	Right:       "┃",
	Bottom:      "━",
	BottomRight: "┻",
}

var borderRowLast = lipgloss.Border{
	Right:       "┃",
	Bottom:      "━",
	BottomRight: "┛",
}

func (m Model) renderRow(i int) string {
	row := m.rows[i]
	last := i == len(m.rows)-1
	highlighted := i == m.rowCursorIndex

	columnStrings := []string{}

	baseStyle := lipgloss.NewStyle()

	if m.focused && highlighted {
		baseStyle = m.highlightStyle
	}

	for i, header := range m.headers {
		var str string

		if header.Key == columnKeySelect {
			if row.selected {
				str = "[x]"
			} else {
				str = "[ ]"
			}
		} else if entry, exists := row.Data[header.Key]; exists {
			str = fmt.Sprintf("%v", entry)
		}

		cellStyle := baseStyle.Copy()

		if i == 0 {
			cellStyle = cellStyle.BorderStyle(borderRowLeft).BorderRight(true).BorderLeft(true)
		} else if i < len(m.headers)-1 {
			cellStyle = cellStyle.BorderStyle(borderRowMiddle).BorderRight(true)
		} else {
			cellStyle = cellStyle.BorderStyle(borderRowLast).BorderRight(true)
		}

		if last {
			cellStyle = cellStyle.BorderBottom(true)
		}

		dataStr := row.Style.Render(fmt.Sprintf(header.fmtString, limitStr(str, header.Width)))

		columnStrings = append(columnStrings, cellStyle.Render(dataStr))
	}

	return lipgloss.JoinHorizontal(lipgloss.Bottom, columnStrings...)
}
