package table

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// RowData is a map of string column keys to interface{} data.
type RowData map[string]interface{}

// Row represents a row in the table with some data keyed to the table columns>
// Can have a style applied to it such as color/bold.  Create using NewRow().
type Row struct {
	Style lipgloss.Style
	Data  RowData

	selected bool
}

// NewRow creates a new row and copies the given row data.
func NewRow(data RowData) Row {
	row := Row{
		Data: make(map[string]interface{}),
	}

	for key, val := range data {
		// Doesn't deep copy val, but close enough for now...
		row.Data[key] = val
	}

	return row
}

// WithStyle uses the given style for the text in the row.
func (r Row) WithStyle(style lipgloss.Style) Row {
	r.Style = style.Copy()

	return r
}

// This is somewhat complicated but at the moment splitting this out feels like
// it would just make things harder to read.  May revisit later.
// nolint: cyclop
func (m Model) renderRow(rowIndex int, last bool) string {
	numColumns := len(m.columns)
	row := m.GetRows()[rowIndex]
	highlighted := rowIndex == m.rowCursorIndex

	columnStrings := []string{}

	baseStyle := row.Style.Copy()

	if m.focused && highlighted {
		baseStyle = baseStyle.Inherit(m.highlightStyle)
	}

	stylesInner, stylesLast := m.styleRows()

	for columnIndex, column := range m.columns {
		cellStyle := baseStyle.Copy().Inherit(column.style)

		var str string

		if column.Key == columnKeySelect {
			if row.selected {
				str = m.selectedText
			} else {
				str = m.unselectedText
			}
		} else if entry, exists := row.Data[column.Key]; exists {
			switch entry := entry.(type) {
			case StyledCell:
				str = fmt.Sprintf("%v", entry.Data)
				cellStyle = entry.Style.Copy().Inherit(cellStyle)
			default:
				str = fmt.Sprintf("%v", entry)
			}
		}

		var rowStyles borderStyleRow

		if !last {
			rowStyles = stylesInner
		} else {
			rowStyles = stylesLast
		}

		if columnIndex == 0 {
			cellStyle = cellStyle.Inherit(rowStyles.left)
		} else if columnIndex < numColumns-1 {
			cellStyle = cellStyle.Inherit(rowStyles.inner)
		} else {
			cellStyle = cellStyle.Inherit(rowStyles.right)
		}

		cellStr := cellStyle.Render(limitStr(str, column.Width))

		columnStrings = append(columnStrings, cellStr)
	}

	return lipgloss.JoinHorizontal(lipgloss.Bottom, columnStrings...)
}
