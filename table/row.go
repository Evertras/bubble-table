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

func (m Model) renderRow(i int) string {
	numHeaders := len(m.headers)
	row := m.rows[i]
	last := i == len(m.rows)-1
	highlighted := i == m.rowCursorIndex

	columnStrings := []string{}

	baseStyle := lipgloss.NewStyle()

	if m.focused && highlighted {
		baseStyle = m.highlightStyle
	}

	var (
		rowStyleLeft  lipgloss.Style
		rowStyleInner lipgloss.Style
		rowStyleRight lipgloss.Style

		rowLastStyleLeft  lipgloss.Style
		rowLastStyleInner lipgloss.Style
		rowLastStyleRight lipgloss.Style
	)

	if numHeaders == 1 {
		rowStyleLeft = m.border.styleSingleColumnInner

		rowLastStyleLeft = m.border.styleSingleColumnBottom
	} else {
		rowStyleLeft = m.border.styleMultiLeft
		rowStyleInner = m.border.styleMultiInner
		rowStyleRight = m.border.styleMultiRight

		rowLastStyleLeft = m.border.styleMultiBottomLeft
		rowLastStyleInner = m.border.styleMultiBottom
		rowLastStyleRight = m.border.styleMultiBottomRight
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

		if !last {
			if i == 0 {
				cellStyle = rowStyleLeft
			} else if i < numHeaders-1 {
				cellStyle = rowStyleInner
			} else {
				cellStyle = rowStyleRight
			}
		} else {
			if i == 0 {
				cellStyle = rowLastStyleLeft
			} else if i < numHeaders-1 {
				cellStyle = rowLastStyleInner
			} else {
				cellStyle = rowLastStyleRight
			}
		}

		dataStr := row.Style.Render(fmt.Sprintf(header.fmtString, limitStr(str, header.Width)))

		columnStrings = append(columnStrings, cellStyle.Render(dataStr))
	}

	return lipgloss.JoinHorizontal(lipgloss.Bottom, columnStrings...)
}
