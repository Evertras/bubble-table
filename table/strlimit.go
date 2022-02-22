package table

func limitStr(s string, maxLen int) string {
	if maxLen == 0 {
		return ""
	}

	if len(s) > maxLen {
		return s[:maxLen-1] + "…"
	}

	return s
}
