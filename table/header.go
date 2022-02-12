package table

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type Header struct {
	Title string
	Key   string
	Width int
	Style lipgloss.Style

	fmtString string
}

func NewHeader(key, title string, width int) Header {
	return Header{
		Key:       key,
		Title:     title,
		Width:     width,
		fmtString: fmt.Sprintf("%%%ds", width),
	}
}

func (h Header) WithStyle(style lipgloss.Style) Header {
	h.Style = style.Copy()
	return h
}

// https://www.w3.org/TR/xml-entity-names/025.html

var borderHeaderFirst = lipgloss.Border{
	Top:         "━",
	Bottom:      "━",
	Left:        "┃",
	Right:       "┃",
	TopRight:    "┳",
	TopLeft:     "┏",
	BottomRight: "╋",
	BottomLeft:  "┣",
}

var borderHeaderTriangleFirst = lipgloss.Border{
	Top:         "━",
	Bottom:      "━",
	Left:        "┃",
	Right:       "┃",
	TopRight:    "┳",
	TopLeft:     "◤",
	BottomRight: "╋",
	BottomLeft:  "◣",
}

var borderHeaderMiddle = lipgloss.Border{
	Top:         "━",
	Bottom:      "━",
	Left:        "",
	Right:       "┃",
	TopRight:    "┳",
	TopLeft:     "",
	BottomRight: "╋",
	BottomLeft:  "",
}

var borderHeaderLast = lipgloss.Border{
	Top:         "━",
	Bottom:      "━",
	Left:        "",
	Right:       "┃",
	TopRight:    "┓",
	TopLeft:     "",
	BottomRight: "┫",
	BottomLeft:  "",
}

var borderHeaderTriangleLast = lipgloss.Border{
	Top:         "━",
	Bottom:      "━",
	Left:        "",
	Right:       "┃",
	TopRight:    "◥",
	TopLeft:     "",
	BottomRight: "◢",
	BottomLeft:  "",
}
