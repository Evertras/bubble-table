package table

import "github.com/charmbracelet/lipgloss"

// HeaderBorder defines the borders around the header
type BorderHeader struct {
	Top      string
	Left     string
	Right    string
	TopRight string
	TopLeft  string
	TopInner string

	BottomFlat      string
	BottomRightFlat string
	BottomLeftFlat  string
	BottomInnerFlat string

	BottomJunction      string
	BottomLeftJunction  string
	BottomRightJunction string
	BottomInnerJunction string

	InnerDivider string

	BaseStyle lipgloss.Style

	styleLeftFlat  lipgloss.Style
	styleInnerFlat lipgloss.Style
	styleRightFlat lipgloss.Style

	styleLeftJunction  lipgloss.Style
	styleInnerJunction lipgloss.Style
	styleRightJunction lipgloss.Style
}

// BodyBorder defines the row-by-row borders/dividers
type BorderBody struct {
	Left    string
	Right   string
	Divider string
}

var (
	// https://www.w3.org/TR/xml-entity-names/025.html

	borderHeaderDefault = BorderHeader{
		Top:      "━",
		Left:     "┃",
		Right:    "┃",
		TopRight: "┓",
		TopLeft:  "┏",

		TopInner:     "┳",
		InnerDivider: "┃",

		BottomFlat:      "━",
		BottomRightFlat: "┛",
		BottomLeftFlat:  "┗",
		BottomInnerFlat: "┻",

		BottomJunction:      "━",
		BottomInnerJunction: "╋",
		BottomLeftJunction:  "┣",
		BottomRightJunction: "┫",
	}
)

func init() {
	borderHeaderDefault.generateStyles()
}

func (b *BorderHeader) generateStyles() {
	b.styleLeftFlat = b.BaseStyle.Copy().BorderStyle(
		lipgloss.Border{
			Top:    b.Top,
			Left:   b.Left,
			Right:  b.InnerDivider,
			Bottom: b.BottomFlat,

			TopLeft:     b.TopLeft,
			TopRight:    b.TopInner,
			BottomLeft:  b.BottomLeftFlat,
			BottomRight: b.BottomInnerFlat,
		},
	)

	b.styleInnerFlat = b.BaseStyle.Copy().BorderStyle(
		lipgloss.Border{
			Top:    b.Top,
			Right:  b.InnerDivider,
			Bottom: b.BottomFlat,

			TopRight:    b.TopInner,
			BottomRight: b.BottomInnerFlat,
		},
	).BorderTop(true).BorderBottom(true).BorderRight(true)

	b.styleRightFlat = b.BaseStyle.Copy().BorderStyle(
		lipgloss.Border{
			Top:    b.Top,
			Right:  b.Right,
			Bottom: b.BottomFlat,

			TopRight:    b.TopRight,
			BottomRight: b.BottomRightFlat,
		},
	).BorderTop(true).BorderBottom(true).BorderRight(true)

	b.styleLeftJunction = b.BaseStyle.Copy().BorderStyle(
		lipgloss.Border{
			Top:    b.Top,
			Left:   b.Left,
			Right:  b.InnerDivider,
			Bottom: b.BottomJunction,

			TopLeft:     b.TopLeft,
			TopRight:    b.TopInner,
			BottomLeft:  b.BottomLeftJunction,
			BottomRight: b.BottomInnerJunction,
		},
	)

	b.styleInnerJunction = b.BaseStyle.Copy().BorderStyle(
		lipgloss.Border{
			Top:    b.Top,
			Right:  b.InnerDivider,
			Bottom: b.BottomJunction,

			TopRight:    b.TopInner,
			BottomRight: b.BottomInnerJunction,
		},
	).BorderTop(true).BorderBottom(true).BorderRight(true)

	b.styleRightJunction = b.BaseStyle.Copy().BorderStyle(
		lipgloss.Border{
			Top:    b.Top,
			Right:  b.Right,
			Bottom: b.BottomFlat,

			TopRight:    b.TopRight,
			BottomRight: b.BottomRightJunction,
		},
	).BorderTop(true).BorderBottom(true).BorderRight(true)
}

func (m Model) BorderHeaderDefault() Model {
	// Already generated styles
	m.borderHeader = borderHeaderDefault
	return m
}

func (m Model) BorderHeader(border BorderHeader) Model {
	m.borderHeader = border
	m.borderHeader.generateStyles()
	return m
}
