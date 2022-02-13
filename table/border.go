package table

import "github.com/charmbracelet/lipgloss"

// Border defines the borders in and around the table
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

	BaseStyle lipgloss.Style

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
)

func init() {
	borderDefault.generateStyles()
}

func (b *Border) generateStyles() {
	b.styleMultiTopLeft = b.BaseStyle.Copy().BorderStyle(
		lipgloss.Border{
			TopLeft:     b.TopLeft,
			Top:         b.Top,
			TopRight:    b.TopJunction,
			Right:       b.InnerDivider,
			BottomRight: b.InnerJunction,
			Bottom:      b.Bottom,
			BottomLeft:  b.LeftJunction,
			Left:        b.Left,
		},
	)

	b.styleMultiTop = b.BaseStyle.Copy().BorderStyle(
		lipgloss.Border{
			Top:    b.Top,
			Right:  b.InnerDivider,
			Bottom: b.Bottom,

			TopRight:    b.TopJunction,
			BottomRight: b.InnerJunction,
		},
	).BorderTop(true).BorderBottom(true).BorderRight(true)

	b.styleMultiTopRight = b.BaseStyle.Copy().BorderStyle(
		lipgloss.Border{
			Top:    b.Top,
			Right:  b.Right,
			Bottom: b.Bottom,

			TopRight:    b.TopRight,
			BottomRight: b.RightJunction,
		},
	).BorderTop(true).BorderBottom(true).BorderRight(true)

	b.styleMultiLeft = b.BaseStyle.Copy().BorderStyle(
		lipgloss.Border{
			Left:  b.Left,
			Right: b.InnerDivider,
		},
	).BorderRight(true).BorderLeft(true)

	b.styleMultiRight = b.BaseStyle.Copy().BorderStyle(
		lipgloss.Border{
			Right: b.Right,
		},
	).BorderRight(true)

	b.styleMultiInner = b.BaseStyle.Copy().BorderStyle(
		lipgloss.Border{
			Right: b.InnerDivider,
		},
	).BorderRight(true)

	b.styleMultiBottomLeft = b.BaseStyle.Copy().BorderStyle(
		lipgloss.Border{
			Left:   b.Left,
			Right:  b.Right,
			Bottom: b.Bottom,

			BottomLeft:  b.BottomLeft,
			BottomRight: b.BottomJunction,
		},
	).BorderLeft(true).BorderBottom(true).BorderRight(true)

	b.styleMultiBottom = b.BaseStyle.Copy().BorderStyle(
		lipgloss.Border{
			Right:  b.Right,
			Bottom: b.Bottom,

			BottomRight: b.BottomJunction,
		},
	).BorderBottom(true).BorderRight(true)

	b.styleMultiBottomRight = b.BaseStyle.Copy().BorderStyle(
		lipgloss.Border{
			Right:  b.Right,
			Bottom: b.Bottom,

			BottomRight: b.BottomRight,
		},
	).BorderBottom(true).BorderRight(true)

	b.styleSingleColumnTop = b.BaseStyle.Copy().BorderStyle(
		lipgloss.Border{
			Top:    b.Top,
			Left:   b.Left,
			Right:  b.Right,
			Bottom: b.Bottom,

			TopLeft:     b.TopLeft,
			TopRight:    b.TopRight,
			BottomLeft:  b.LeftJunction,
			BottomRight: b.RightJunction,
		},
	)

	b.styleSingleColumnInner = b.BaseStyle.Copy().BorderStyle(
		lipgloss.Border{
			Left:  b.Left,
			Right: b.Right,
		},
	).BorderRight(true).BorderLeft(true)

	b.styleSingleColumnBottom = b.BaseStyle.Copy().BorderStyle(
		lipgloss.Border{
			Left:   b.Left,
			Right:  b.Right,
			Bottom: b.Bottom,

			BottomLeft:  b.BottomLeft,
			BottomRight: b.BottomRight,
		},
	).BorderRight(true).BorderLeft(true).BorderBottom(true)

	b.styleSingleRowLeft = b.BaseStyle.Copy().BorderStyle(
		lipgloss.Border{
			Top:    b.Top,
			Left:   b.Left,
			Right:  b.Right,
			Bottom: b.Bottom,

			BottomLeft:  b.BottomLeft,
			BottomRight: b.BottomJunction,
			TopRight:    b.TopJunction,
			TopLeft:     b.TopLeft,
		},
	)

	b.styleSingleRowInner = b.BaseStyle.Copy().BorderStyle(
		lipgloss.Border{
			Top:    b.Top,
			Right:  b.Right,
			Bottom: b.Bottom,

			BottomRight: b.BottomJunction,
			TopRight:    b.TopJunction,
		},
	).BorderTop(true).BorderBottom(true).BorderRight(true)

	b.styleSingleRowRight = b.BaseStyle.Copy().BorderStyle(
		lipgloss.Border{
			Top:    b.Top,
			Right:  b.Right,
			Bottom: b.Bottom,

			BottomRight: b.BottomRight,
			TopRight:    b.TopRight,
		},
	).BorderTop(true).BorderBottom(true).BorderRight(true)

	b.styleSingleCell = b.BaseStyle.Copy().BorderStyle(
		lipgloss.Border{
			Top:    b.Top,
			Left:   b.Left,
			Right:  b.Right,
			Bottom: b.Bottom,

			BottomLeft:  b.BottomLeft,
			BottomRight: b.BottomRight,
			TopRight:    b.TopRight,
			TopLeft:     b.TopLeft,
		},
	)
}

func (m Model) BorderDefault() Model {
	// Already generated styles
	m.border = borderDefault

	return m
}

func (m Model) Border(border Border) Model {
	border.generateStyles()

	m.border = border

	return m
}
