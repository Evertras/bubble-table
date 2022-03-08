package table

import "github.com/mattn/go-runewidth"

func limitStr(str string, maxLen int) string {
	if maxLen == 0 {
		return ""
	}
	if runewidth.StringWidth(str) > maxLen {
		return runewidth.Truncate(str, maxLen-1, "") + "…"
	}

	return str
}
