package table

import (
	"fmt"
)

// Column is a column in the table
type Column struct {
	Title string
	Key   string
	Width int

	fmtString string
}

// NewColumn creates a new column with the given information
func NewColumn(key, title string, width int) Column {
	return Column{
		Key:       key,
		Title:     title,
		Width:     width,
		fmtString: fmt.Sprintf("%%%ds", width),
	}
}
