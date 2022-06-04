package table

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

const columnKeyOverflow = "___overflow___"

// RowData is a map of string column keys to interface{} data.  Data with a key
// that matches a column key will be displayed.  Data with a key that does not
// match a column key will not be displayed, but will remain attached to the Row.
// This can be useful for attaching hidden metadata for future reference when
// retrieving rows.
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

func (m Model) renderRowColumnData(row Row, column Column, rowStyle lipgloss.Style, borderStyle lipgloss.Style) string {
	cellStyle := rowStyle.Copy().Inherit(column.style).Inherit(m.baseStyle)

	var str string

	if column.key == columnKeySelect {
		if row.selected {
			str = m.selectedText
		} else {
			str = m.unselectedText
		}
	} else if column.key == columnKeyOverflow {
		str = ">"
	} else if entry, exists := row.Data[column.key]; exists {
		switch entry := entry.(type) {
		case StyledCell:
			str = fmt.Sprintf("%v", entry.Data)
			cellStyle = entry.Style.Copy().Inherit(cellStyle)
		default:
			str = fmt.Sprintf("%v", entry)
		}
	}

	cellStyle = cellStyle.Inherit(borderStyle)
	cellStr := cellStyle.Render(limitStr(str, column.width))

	return cellStr
}

func (m Model) renderRow(rowIndex int, last bool) string {
	numColumns := len(m.columns)
	row := m.GetVisibleRows()[rowIndex]
	highlighted := rowIndex == m.rowCursorIndex
	totalRenderedWidth := 0

	columnStrings := []string{}

	rowStyle := row.Style.Copy()

	if m.focused && highlighted {
		rowStyle = rowStyle.Inherit(m.highlightStyle)
	}

	stylesInner, stylesLast := m.styleRows()

	for columnIndex, column := range m.columns {
		var borderStyle lipgloss.Style

		var rowStyles borderStyleRow

		if !last {
			rowStyles = stylesInner
		} else {
			rowStyles = stylesLast
		}

		if columnIndex == 0 {
			borderStyle = rowStyles.left
		} else if columnIndex < numColumns-1 {
			borderStyle = rowStyles.inner
		} else {
			borderStyle = rowStyles.right
		}

		cellStr := m.renderRowColumnData(row, column, rowStyle, borderStyle)

		if m.maxTotalWidth != 0 {
			renderedWidth := lipgloss.Width(cellStr)

			const borderAdjustment = 1

			if totalRenderedWidth+renderedWidth > m.maxTotalWidth-borderAdjustment*2 {
				overflowWidth := m.maxTotalWidth - totalRenderedWidth - borderAdjustment
				overflowStyle := genOverflowStyle(rowStyles.right, overflowWidth)
				overflowColumn := genOverflowColumn(overflowWidth)
				overflowStr := m.renderRowColumnData(row, overflowColumn, rowStyle, overflowStyle)

				columnStrings = append(columnStrings, overflowStr)

				break
			}

			totalRenderedWidth += renderedWidth
		}

		columnStrings = append(columnStrings, cellStr)
	}

	return lipgloss.JoinHorizontal(lipgloss.Bottom, columnStrings...)
}
