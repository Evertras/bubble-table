package main

import (
	"fmt"
	"log"

	tea "charm.land/bubbletea/v2"
	"github.com/evertras/bubble-table/table"
)

const (
	columnKeyName    = "name"
	columnKeyElement = "element"
	columnKeyHP      = "hp"
)

type Model struct {
	tables       []table.Model
	tableLabels  []string
	currentTable int
}

func NewModel() Model {
	columns := []table.Column{
		table.NewColumn(columnKeyName, "Name", 12),
		table.NewColumn(columnKeyElement, "Element", 10),
		table.NewColumn(columnKeyHP, "HP", 6),
	}

	rows := []table.Row{
		table.NewRow(table.RowData{
			columnKeyName:    "Bulbasaur",
			columnKeyElement: "Grass",
			columnKeyHP:      45,
		}),
		table.NewRow(table.RowData{
			columnKeyName:    "Charmander",
			columnKeyElement: "Fire",
			columnKeyHP:      39,
		}),
		table.NewRow(table.RowData{
			columnKeyName:    "Squirtle",
			columnKeyElement: "Water",
			columnKeyHP:      44,
		}),
		table.NewRow(table.RowData{
			columnKeyName:    "Pikachu",
			columnKeyElement: "Electric",
			columnKeyHP:      35,
		}),
	}

	return Model{
		tables: []table.Model{
			table.New(columns).WithRows(rows).WithRowBorder(true),
			table.New(columns).WithRows(rows).WithRowBorder(true).BorderRounded(),
			table.New(columns).WithRows(rows).WithRowBorder(true).WithOuterBorder(false),
		},
		tableLabels: []string{
			"Default border",
			"Rounded border",
			"No outer border",
		},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit
		case "left":
			m.currentTable = (m.currentTable - 1 + len(m.tables)) % len(m.tables)
		case "right":
			m.currentTable = (m.currentTable + 1) % len(m.tables)
		}
	}

	return m, nil
}

func (m Model) View() tea.View {
	label := m.tableLabels[m.currentTable]
	position := fmt.Sprintf("%d/%d", m.currentTable+1, len(m.tables))

	header := fmt.Sprintf("Row border example  [← →] cycle  [q] quit\n%s (%s)\n\n", label, position)

	return tea.NewView(header + m.tables[m.currentTable].View())
}

func main() {
	p := tea.NewProgram(NewModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
