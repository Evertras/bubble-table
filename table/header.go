package table

import (
	"fmt"
)

type Header struct {
	Title string
	Key   string
	Width int

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
