package table

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

func limitStr(str string, maxLen int) string {
	if maxLen == 0 {
		return ""
	}

	newLineIndex := strings.Index(str, "\n")
	if newLineIndex > -1 {
		str = str[:newLineIndex] + "…"
	}

	if ansi.StringWidth(str) > maxLen {
		return ansi.Truncate(str, maxLen, "…")
	}

	return str
}
