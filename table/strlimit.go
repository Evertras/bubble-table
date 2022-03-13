package table

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

func limitStr(str string, maxLen int) string {
	if maxLen == 0 {
		return ""
	}

	newLineIndex := strings.Index(str, "\n")
	if newLineIndex > -1 {
		str = str[:newLineIndex] + "…"
	}

	if runewidth.StringWidth(str) > maxLen {
		return runewidth.Truncate(str, maxLen-1, "") + "…"
	}

	return str
}
