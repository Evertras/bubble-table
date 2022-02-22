package table

func limitStr(str string, maxLen int) string {
	if maxLen == 0 {
		return ""
	}

	if len(str) > maxLen {
		return str[:maxLen-1] + "…"
	}

	return str
}
