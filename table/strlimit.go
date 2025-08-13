package table

import (
	"strings"

	"github.com/muesli/reflow/ansi"
	"github.com/muesli/reflow/truncate"
)

func limitStr(str string, maxLen int) string {
	if maxLen == 0 {
		return ""
	}

	newLineIndex := strings.Index(str, "\n")
	if newLineIndex > -1 {
		str = str[:newLineIndex] + "…"
	}

	if ansi.PrintableRuneWidth(str) > maxLen {
		// #nosec: G115
		return truncate.StringWithTail(str, uint(maxLen), "…")
	}

	return str
}
