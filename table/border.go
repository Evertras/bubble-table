package table

import (
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
)

// Border defines the borders in and around the table.
type Border struct {
	Top         string
	Left        string
	Right       string
	Bottom      string
	TopRight    string
	TopLeft     string
	BottomRight string
	BottomLeft  string

	TopJunction    string
	LeftJunction   string
	RightJunction  string
	BottomJunction string

	InnerJunction string

	InnerDivider string

	// Styles for 2x2 tables and larger
	styleMultiTopLeft     lipgloss.Style
	styleMultiTop         lipgloss.Style
	styleMultiTopRight    lipgloss.Style
	styleMultiRight       lipgloss.Style
	styleMultiBottomRight lipgloss.Style
	styleMultiBottom      lipgloss.Style
	styleMultiBottomLeft  lipgloss.Style
	styleMultiLeft        lipgloss.Style
	styleMultiInner       lipgloss.Style

	// Styles for a single column table
	styleSingleColumnTop    lipgloss.Style
	styleSingleColumnInner  lipgloss.Style
	styleSingleColumnBottom lipgloss.Style

	// Styles for a single row table
	styleSingleRowLeft  lipgloss.Style
	styleSingleRowInner lipgloss.Style
	styleSingleRowRight lipgloss.Style

	// Style for a table with only one cell
	styleSingleCell lipgloss.Style

	// Style for the footer
	styleFooter lipgloss.Style

	// foreground is the colour applied to plain-string border lines (top,
	// separator, bottom).  nil means use the terminal default.
	foreground color.Color
}

var (
	// https://www.w3.org/TR/xml-entity-names/025.html

	borderDefault = Border{
		Top:    "━",
		Left:   "┃",
		Right:  "┃",
		Bottom: "━",

		TopRight:    "┓",
		TopLeft:     "┏",
		BottomRight: "┛",
		BottomLeft:  "┗",

		TopJunction:    "┳",
		LeftJunction:   "┣",
		RightJunction:  "┫",
		BottomJunction: "┻",
		InnerJunction:  "╋",

		InnerDivider: "┃",
	}

	borderRounded = Border{
		Top:    "─",
		Left:   "│",
		Right:  "│",
		Bottom: "─",

		TopRight:    "╮",
		TopLeft:     "╭",
		BottomRight: "╯",
		BottomLeft:  "╰",

		TopJunction:    "┬",
		LeftJunction:   "├",
		RightJunction:  "┤",
		BottomJunction: "┴",
		InnerJunction:  "┼",

		InnerDivider: "│",
	}
)

func init() {
	borderDefault.generateStyles()
	borderRounded.generateStyles()
}

func (b *Border) generateStyles() {
	b.generateMultiStyles()
	b.generateSingleColumnStyles()
	b.generateSingleRowStyles()
	b.generateSingleCellStyle()

	// The footer is a full-width single cell with a bottom border.  It needs
	// the actual border characters for all edges because renderFooter may
	// optionally add a top border (when there are no data rows above).
	b.styleFooter = lipgloss.NewStyle().
		BorderStyle(lipgloss.Border{
			Left:  b.Left,
			Right: b.Right,

			Top:      b.Top,
			TopLeft:  b.TopLeft,
			TopRight: b.TopRight,

			Bottom:      b.Bottom,
			BottomLeft:  b.BottomLeft,
			BottomRight: b.BottomRight,
		}).
		Align(lipgloss.Right).
		BorderBottom(true).
		BorderRight(true).
		BorderLeft(true)
}

// buildBorderLine constructs a horizontal border line for the given column
// widths. leftChar is the leftmost character, midChar fills each column,
// junctionChar separates columns, and rightChar closes the line.
// If b.foreground is set the line is wrapped in a lipgloss Foreground style
// so that it renders in the same colour as the cell border characters.
func (b *Border) buildBorderLine(columnWidths []int, leftChar, midChar, junctionChar, rightChar string) string {
	var buf strings.Builder

	buf.WriteString(leftChar)

	for i, w := range columnWidths {
		buf.WriteString(strings.Repeat(midChar, w))

		if i < len(columnWidths)-1 {
			buf.WriteString(junctionChar)
		}
	}

	buf.WriteString(rightChar)

	line := buf.String()

	if b.foreground != nil {
		line = lipgloss.NewStyle().Foreground(b.foreground).Render(line)
	}

	return line
}

// buildTopBorderLine builds the top border line for a multi-column table.
func (b *Border) buildTopBorderLine(columnWidths []int) string {
	return b.buildBorderLine(columnWidths, b.TopLeft, b.Top, b.TopJunction, b.TopRight)
}

// buildSeparatorLine builds the separator line between the header and first
// data row for a multi-column table.
func (b *Border) buildSeparatorLine(columnWidths []int) string {
	return b.buildBorderLine(columnWidths, b.LeftJunction, b.Bottom, b.InnerJunction, b.RightJunction)
}

// buildInnerSeparatorLine builds a separator line with no outer junction
// characters, used when the outer border is hidden.
func (b *Border) buildInnerSeparatorLine(columnWidths []int) string {
	return b.buildBorderLine(columnWidths, "", b.Bottom, b.InnerJunction, "")
}

// buildBottomBorderLine builds the bottom border line for a multi-column
// table. When hasFooter is true the corners are replaced with junctions so
// the footer can attach below.
func (b *Border) buildBottomBorderLine(columnWidths []int, hasFooter bool) string {
	left := b.BottomLeft
	right := b.BottomRight

	if hasFooter {
		left = b.LeftJunction
		right = b.RightJunction
	}

	return b.buildBorderLine(columnWidths, left, b.Bottom, b.BottomJunction, right)
}

func (b *Border) styleLeftWithFooter(original lipgloss.Style) lipgloss.Style {
	border := original.GetBorderStyle()

	border.BottomLeft = b.LeftJunction

	return original.BorderStyle(border)
}

func (b *Border) styleRightWithFooter(original lipgloss.Style) lipgloss.Style {
	border := original.GetBorderStyle()

	border.BottomRight = b.RightJunction

	return original.BorderStyle(border)
}

func (b *Border) styleBothWithFooter(original lipgloss.Style) lipgloss.Style {
	border := original.GetBorderStyle()

	border.BottomLeft = b.LeftJunction
	border.BottomRight = b.RightJunction

	return original.BorderStyle(border)
}

func (b *Border) generateMultiStyles() {
	// All cell styles have ONLY left and/or right borders.
	// Top, bottom, and separator lines are rendered as plain strings by
	// the header/row rendering functions, not by these styles.

	b.styleMultiTopLeft = lipgloss.NewStyle().BorderStyle(
		lipgloss.Border{
			Left:  b.Left,
			Right: b.InnerDivider,
		},
	).BorderLeft(true).BorderRight(true)

	b.styleMultiTop = lipgloss.NewStyle().BorderStyle(
		lipgloss.Border{
			Right: b.InnerDivider,
		},
	).BorderRight(true)

	b.styleMultiTopRight = lipgloss.NewStyle().BorderStyle(
		lipgloss.Border{
			Right: b.Right,
		},
	).BorderRight(true)

	b.styleMultiLeft = lipgloss.NewStyle().BorderStyle(
		lipgloss.Border{
			Left:  b.Left,
			Right: b.InnerDivider,
		},
	).BorderRight(true).BorderLeft(true)

	b.styleMultiRight = lipgloss.NewStyle().BorderStyle(
		lipgloss.Border{
			Right: b.Right,
		},
	).BorderRight(true)

	b.styleMultiInner = lipgloss.NewStyle().BorderStyle(
		lipgloss.Border{
			Right: b.InnerDivider,
		},
	).BorderRight(true)

	b.styleMultiBottomLeft = lipgloss.NewStyle().BorderStyle(
		lipgloss.Border{
			Left:  b.Left,
			Right: b.InnerDivider,
		},
	).BorderLeft(true).BorderRight(true)

	b.styleMultiBottom = lipgloss.NewStyle().BorderStyle(
		lipgloss.Border{
			Right: b.InnerDivider,
		},
	).BorderRight(true)

	b.styleMultiBottomRight = lipgloss.NewStyle().BorderStyle(
		lipgloss.Border{
			Right: b.Right,
		},
	).BorderRight(true)
}

func (b *Border) generateSingleColumnStyles() {
	// Cell styles have only left/right borders; top/bottom border lines are
	// built as plain strings by the rendering functions.

	b.styleSingleColumnTop = lipgloss.NewStyle().BorderStyle(
		lipgloss.Border{
			Left:  b.Left,
			Right: b.Right,
		},
	).BorderLeft(true).BorderRight(true)

	b.styleSingleColumnInner = lipgloss.NewStyle().BorderStyle(
		lipgloss.Border{
			Left:  b.Left,
			Right: b.Right,
		},
	).BorderRight(true).BorderLeft(true)

	b.styleSingleColumnBottom = lipgloss.NewStyle().BorderStyle(
		lipgloss.Border{
			Left:  b.Left,
			Right: b.Right,
		},
	).BorderRight(true).BorderLeft(true)
}

func (b *Border) generateSingleRowStyles() {
	// Cell styles have only left/right borders; top/bottom border lines are
	// built as plain strings by the rendering functions.

	b.styleSingleRowLeft = lipgloss.NewStyle().BorderStyle(
		lipgloss.Border{
			Left:  b.Left,
			Right: b.InnerDivider,
		},
	).BorderLeft(true).BorderRight(true)

	b.styleSingleRowInner = lipgloss.NewStyle().BorderStyle(
		lipgloss.Border{
			Right: b.InnerDivider,
		},
	).BorderRight(true)

	b.styleSingleRowRight = lipgloss.NewStyle().BorderStyle(
		lipgloss.Border{
			Right: b.Right,
		},
	).BorderRight(true)
}

func (b *Border) generateSingleCellStyle() {
	// Cell style has only left/right borders; top/bottom border lines are
	// built as plain strings by the rendering functions.
	b.styleSingleCell = lipgloss.NewStyle().BorderStyle(
		lipgloss.Border{
			Left:  b.Left,
			Right: b.Right,
		},
	).BorderLeft(true).BorderRight(true)
}

// BorderDefault uses the basic square border, useful to reset the border if
// it was changed somehow.
func (m Model) BorderDefault() Model {
	fg := m.border.foreground
	m.border = borderDefault
	m.border.foreground = fg

	return m
}

// BorderRounded uses a thin, rounded border.
func (m Model) BorderRounded() Model {
	fg := m.border.foreground
	m.border = borderRounded
	m.border.foreground = fg

	return m
}

// Border uses the given border components to render the table.
func (m Model) Border(border Border) Model {
	border.generateStyles()

	m.border = border

	return m
}

type borderStyleRow struct {
	left  lipgloss.Style
	inner lipgloss.Style
	right lipgloss.Style
}

func (b *borderStyleRow) inherit(s lipgloss.Style) {
	b.left = b.left.Inherit(s)
	b.inner = b.inner.Inherit(s)
	b.right = b.right.Inherit(s)
}

// There's a lot of branches here, but splitting it up further would make it
// harder to follow. So just be careful with comments and make sure it's tested!
//
//nolint:nestif
func (m Model) styleHeaders() borderStyleRow {
	hasRows := len(m.GetVisibleRows()) > 0 || m.calculatePadding(0) > 0
	singleColumn := len(m.columns) == 1
	styles := borderStyleRow{}

	// Possible configurations:
	// - Single cell
	// - Single row
	// - Single column
	// - Multi

	if singleColumn {
		if hasRows {
			// Single column
			styles.left = m.border.styleSingleColumnTop
			styles.inner = styles.left
			styles.right = styles.left
		} else {
			// Single cell
			styles.left = m.border.styleSingleCell
			styles.inner = styles.left
			styles.right = styles.left

			if m.hasFooter() {
				styles.left = m.border.styleBothWithFooter(styles.left)
			}
		}
	} else if !hasRows {
		// Single row
		styles.left = m.border.styleSingleRowLeft
		styles.inner = m.border.styleSingleRowInner
		styles.right = m.border.styleSingleRowRight

		if m.hasFooter() {
			styles.left = m.border.styleLeftWithFooter(styles.left)
			styles.right = m.border.styleRightWithFooter(styles.right)
		}
	} else {
		// Multi
		styles.left = m.border.styleMultiTopLeft
		styles.inner = m.border.styleMultiTop
		styles.right = m.border.styleMultiTopRight
	}

	styles.inherit(m.headerStyle)

	// HeaderStyle is for text styling only. Strip any border-side flags that
	// may have been inherited: top/bottom border lines are rendered as separate
	// plain strings, and inner/right cells never carry a left border (each
	// column's left divider is provided by the previous column's right border).
	styles.left = styles.left.BorderTop(false).BorderBottom(false)
	styles.inner = styles.inner.BorderTop(false).BorderBottom(false).BorderLeft(false)
	styles.right = styles.right.BorderTop(false).BorderBottom(false).BorderLeft(false)

	if !m.outerBorder {
		styles.left = styles.left.BorderLeft(false)
		styles.right = styles.right.BorderRight(false)
	}

	return styles
}

func (m Model) styleRows() (inner borderStyleRow, last borderStyleRow) {
	if len(m.columns) == 1 {
		inner.left = m.border.styleSingleColumnInner
		inner.inner = inner.left
		inner.right = inner.left

		last.left = m.border.styleSingleColumnBottom

		if m.hasFooter() {
			last.left = m.border.styleBothWithFooter(last.left)
		}

		last.inner = last.left
		last.right = last.left
	} else {
		inner.left = m.border.styleMultiLeft
		inner.inner = m.border.styleMultiInner
		inner.right = m.border.styleMultiRight

		last.left = m.border.styleMultiBottomLeft
		last.inner = m.border.styleMultiBottom
		last.right = m.border.styleMultiBottomRight

		if m.hasFooter() {
			last.left = m.border.styleLeftWithFooter(last.left)
			last.right = m.border.styleRightWithFooter(last.right)
		}
	}

	if !m.outerBorder {
		inner.left = inner.left.BorderLeft(false)
		inner.right = inner.right.BorderRight(false)
		last.left = last.left.BorderLeft(false)
		last.right = last.right.BorderRight(false)
	}

	return inner, last
}
