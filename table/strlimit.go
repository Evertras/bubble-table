package table

func limitStr(str string, maxLen int) string {
	if maxLen == 0 {
		return ""
	}

	if len([]rune(str)) > maxLen {
		return string([]rune(str)[:maxLen-1]) + "…"
	}

	return str
}
