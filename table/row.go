package table

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

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

// nolint: nestif // This has many ifs, but they're short
func (m Model) renderRowColumnData(row Row, column Column, rowStyle lipgloss.Style, borderStyle lipgloss.Style) string {
	cellStyle := rowStyle.Copy().Inherit(column.style).Inherit(m.baseStyle)

	var str string

	if column.key == columnKeySelect {
		if row.selected {
			str = m.selectedText
		} else {
			str = m.unselectedText
		}
	} else if column.key == columnKeyOverflowRight {
		str = ">"
	} else if column.key == columnKeyOverflowLeft {
		str = "<"
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

// This is long and could use some refactoring in the future, but not quite sure
// how to pick it apart yet.
// nolint: funlen, cyclop
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

	if m.horizontalScrollOffsetCol > 0 {
		var borderStyle lipgloss.Style

		if !last {
			borderStyle = stylesInner.left.Copy()
		} else {
			borderStyle = stylesLast.left.Copy()
		}

		rendered := m.renderRowColumnData(row, genOverflowColumnLeft(1), rowStyle, borderStyle)

		totalRenderedWidth += lipgloss.Width(rendered)

		columnStrings = append(columnStrings, rendered)
	}

	for columnIndex, column := range m.columns {
		if columnIndex < m.horizontalScrollOffsetCol {
			continue
		}

		var borderStyle lipgloss.Style

		var rowStyles borderStyleRow

		if !last {
			rowStyles = stylesInner
		} else {
			rowStyles = stylesLast
		}

		if len(columnStrings) == 0 {
			borderStyle = rowStyles.left
		} else if columnIndex < numColumns-1 {
			borderStyle = rowStyles.inner
		} else {
			borderStyle = rowStyles.right
		}

		cellStr := m.renderRowColumnData(row, column, rowStyle, borderStyle)

		if m.maxTotalWidth != 0 {
			renderedWidth := lipgloss.Width(cellStr)

			const (
				borderAdjustment = 1
				overflowColWidth = 2
			)

			targetWidth := m.maxTotalWidth - overflowColWidth

			if columnIndex == len(m.columns)-1 {
				// If this is the last header, we don't need to account for the
				// overflow arrow column
				targetWidth = m.maxTotalWidth
			}

			if totalRenderedWidth+renderedWidth > targetWidth {
				overflowWidth := m.maxTotalWidth - totalRenderedWidth - borderAdjustment
				overflowStyle := genOverflowStyle(rowStyles.right, overflowWidth)
				overflowColumn := genOverflowColumnRight(overflowWidth)
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
