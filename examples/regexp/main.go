package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/evertras/bubble-table/table"
)

const (
	columnKeyTime    = "time"
	columnKeyLevel   = "level"
	columnKeyMessage = "message"
)

// regexMatcher caches the last compiled pattern so typing a query doesn't
// recompile the regex on every row on every keystroke. It's safe to reuse
// here because bubbletea's Update runs on a single goroutine.
type regexMatcher struct {
	pattern string
	re      *regexp.Regexp
}

func (m *regexMatcher) matches(input table.FilterFuncInput) bool {
	if input.Filter == "" {
		return true
	}

	if input.Filter != m.pattern {
		m.pattern = input.Filter
		m.re, _ = regexp.Compile(input.Filter) // nil on invalid pattern, handled below
	}

	if m.re == nil {
		return false
	}

	for _, column := range input.Columns {
		data, ok := input.Row.Data[column.Key()]
		if !ok {
			continue
		}

		if m.re.MatchString(fmt.Sprint(data)) {
			return true
		}
	}

	return false
}

type Model struct {
	table table.Model
}

func NewModel() Model {
	columns := []table.Column{
		table.NewColumn(columnKeyTime, "Time", 8),
		table.NewColumn(columnKeyLevel, "Level", 7),
		table.NewColumn(columnKeyMessage, "Message", 32),
	}

	data := [][3]string{
		{"10:01:02", "INFO", "GET /home 200"},
		{"10:01:05", "WARN", "slow response 812ms"},
		{"10:01:10", "ERROR", "GET /api/users 500"},
		{"10:01:12", "INFO", "POST /login 200"},
		{"10:01:15", "ERROR", "DB connection timeout"},
		{"10:01:20", "WARN", "GET /api/orders 404"},
		{"10:01:25", "INFO", "GET /health 200"},
		{"10:01:30", "ERROR", "GET /api/users 503"},
		{"10:01:35", "DEBUG", "cache miss key=42"},
		{"10:01:40", "INFO", "GET /about 200"},
	}

	rows := make([]table.Row, 0, len(data))
	for _, d := range data {
		rows = append(rows, table.NewRow(table.RowData{
			columnKeyTime:    d[0],
			columnKeyLevel:   d[1],
			columnKeyMessage: d[2],
		}))
	}

	matchStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffd700"))
	highlightStyle := lipgloss.NewStyle().Background(lipgloss.Color("#7d56f4"))

	filterMatcher := &regexMatcher{}
	searchMatcher := &regexMatcher{}

	return Model{
		table: table.New(columns).
			WithRows(rows).
			Focused(true).
			Filtered(true).
			WithFilterFunc(filterMatcher.matches).
			Searchable(true).
			WithSearchFunc(searchMatcher.matches).
			WithPageSize(5).
			// A row style func replaces the default highlight, so handle it here too
			WithRowStyleFunc(func(input table.RowStyleFuncInput) lipgloss.Style {
				if input.IsHighlighted {
					return highlightStyle
				}

				if input.IsSearchMatch {
					return matchStyle
				}

				return lipgloss.NewStyle()
			}),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.table, cmd = m.table.Update(msg)

	if msg, ok := msg.(tea.KeyPressMsg); ok {
		switch msg.String() {
		case "ctrl+c", "q":
			if !m.table.GetIsFilterInputFocused() && !m.table.GetIsSearchInputFocused() {
				return m, tea.Quit
			}
		}
	}

	return m, cmd
}

func (m Model) View() tea.View {
	body := strings.Builder{}

	body.WriteString("Press / to filter as a regex, e.g. ^(ERROR|WARN)$ on the Level column\n" +
		"shows only warnings and errors, something plain substring matching can't express.\n" +
		"Press ctrl+f to search, e.g. \\b[45]\\d\\d\\b jumps between 4xx/5xx status codes,\n" +
		"skipping the 200s. esc stops typing, n / N move to next / previous search match.\n" +
		"Press q or ctrl+c to quit\n\n")

	body.WriteString(m.table.View())

	return tea.NewView(body.String())
}

func main() {
	p := tea.NewProgram(NewModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
