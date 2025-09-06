package table

import "github.com/charmbracelet/lipgloss"

// StyledCell represents a cell in the table that has a particular style applied.
// The cell style takes highest precedence and will overwrite more general styles
// from the row, column, or table as a whole.  This style should be generally
// limited to colors, font style, and alignments - spacing style such as margin
// will break the table format.
type StyledCell struct {
	// Data is the content of the cell.
	Data any

	// Style is the specific style to apply. This is ignored if StyleFunc is not nil.
	Style lipgloss.Style

	// StyleFunc is a function that takes the row/column of the cell and
	// returns a lipgloss.Style allowing for dynamic styling based on the cell's
	// content or position. Overrides Style if set.
	StyleFunc StyledCellFunc
}

// StyledCellFuncInput is the input to the StyledCellFunc. Sent as a struct
// to allow for future additions without breaking changes.
type StyledCellFuncInput struct {
	// Data is the data in the cell.
	Data any

	// Column is the column that the cell belongs to.
	Column Column

	// Row is the row that the cell belongs to.
	Row Row

	// GlobalMetadata is the global table metadata that's been set by WithGlobalMetadata
	GlobalMetadata map[string]any
}

// StyledCellFunc is a function that takes various information about the cell and
// returns a lipgloss.Style allowing for easier dynamic styling based on the cell's
// content or position.
type StyledCellFunc = func(input StyledCellFuncInput) lipgloss.Style

// NewStyledCell creates an entry that can be set in the row data and show as
// styled with the given style.
func NewStyledCell(data any, style lipgloss.Style) StyledCell {
	return StyledCell{
		Data:  data,
		Style: style,
	}
}

// NewStyledCellWithStyleFunc creates an entry that can be set in the row data and show as
// styled with the given style function.
func NewStyledCellWithStyleFunc(data any, styleFunc StyledCellFunc) StyledCell {
	return StyledCell{
		Data:      data,
		StyleFunc: styleFunc,
	}
}
