package main

import (
	"log"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/evertras/bubble-table/table"
)

const (
	columnKeyName  = "name"
	columnKeyType  = "type"
	columnKeyColor = "color"
)

type Model struct {
	table table.Model
}

func NewModel() Model {
	columns := []table.Column{
		table.NewColumn(columnKeyName, "Name", 14),
		table.NewColumn(columnKeyType, "Type", 10).WithFiltered(true),
		table.NewColumn(columnKeyColor, "Color", 10),
	}

	data := [][3]string{
		{"Apple", "Fruit", "Red"},
		{"Carrot", "Vegetable", "Orange"},
		{"Banana", "Fruit", "Yellow"},
		{"Leek", "Vegetable", "Green"},
		{"Cherry", "Fruit", "Red"},
		{"Beetroot", "Vegetable", "Purple"},
		{"Lime", "Fruit", "Green"},
		{"Pepper", "Vegetable", "Red"},
		{"Grape", "Fruit", "Purple"},
		{"Pumpkin", "Vegetable", "Orange"},
	}

	rows := make([]table.Row, 0, len(data))
	for _, d := range data {
		rows = append(rows, table.NewRow(table.RowData{
			columnKeyName:  d[0],
			columnKeyType:  d[1],
			columnKeyColor: d[2],
		}))
	}

	matchStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffd700"))
	highlightStyle := lipgloss.NewStyle().Background(lipgloss.Color("#7d56f4"))

	return Model{
		table: table.New(columns).
			WithRows(rows).
			Focused(true).
			Filtered(true).
			Searchable(true).
			WithPageSize(4).
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

	body.WriteString("Press ctrl+f to search all columns (jumps as you type), enter/esc to stop typing.\n" +
		"Then n / N for next / previous match. Matches are highlighted without hiding rows.\n" +
		"Press / to filter by Type; search only looks at the filtered rows.\n" +
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
