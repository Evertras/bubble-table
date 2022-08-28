package main

import (
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
)

const (
	columnKeyName = "name"
	columnKeyType = "type"
	columnKeyWins = "wins"
)

type Model struct {
	simpleTable table.Model

	columnSortKey string
	sortDirection string
}

func NewModel() Model {
	return Model{
		simpleTable: table.New([]table.Column{
			table.NewColumn(columnKeyName, "Name", 13),
			table.NewColumn(columnKeyType, "Type", 13),
			table.NewColumn(columnKeyWins, "Win %", 8).
				WithFormatString("%.1f%%"),
		}).WithRows([]table.Row{
			table.NewRow(table.RowData{
				columnKeyName: "ピカピカ",
				columnKeyType: "Pikachu",
				columnKeyWins: 78.3,
			}),
			table.NewRow(table.RowData{
				columnKeyName: "Zapmouse",
				columnKeyType: "Pikachu",
				columnKeyWins: 3.3,
			}),
			table.NewRow(table.RowData{
				columnKeyName: "Burninator",
				columnKeyType: "Charmander",
				columnKeyWins: 32.1,
			}),
			table.NewRow(table.RowData{
				columnKeyName: "Alphonse",
				columnKeyType: "Pikachu",
				columnKeyWins: 13.8,
			}),
			table.NewRow(table.RowData{
				columnKeyName: "Trogdor",
				columnKeyType: "Charmander",
				columnKeyWins: 99.9,
			}),
			table.NewRow(table.RowData{
				columnKeyName: "Dihydrogen Monoxide",
				columnKeyType: "Squirtle",
				columnKeyWins: 31.348,
			}),
		}),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.simpleTable, cmd = m.simpleTable.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			cmds = append(cmds, tea.Quit)

		case "n":
			m.columnSortKey = columnKeyName
			m.simpleTable = m.simpleTable.SortByAsc(m.columnSortKey)

		case "t":
			m.columnSortKey = columnKeyType
			// Within the same type, order each by wins
			m.simpleTable = m.simpleTable.SortByAsc(m.columnSortKey).ThenSortByDesc(columnKeyWins)

		case "w":
			m.columnSortKey = columnKeyWins
			m.simpleTable = m.simpleTable.SortByDesc(m.columnSortKey)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	body := strings.Builder{}

	body.WriteString("A sorted simple default table\nSort by (n)ame, (t)ype->wins combo, or (w)ins\nCurrently sorting by: " + m.columnSortKey + "\nPress q or ctrl+c to quit\n\n")

	body.WriteString(m.simpleTable.View())

	return body.String()
}

func main() {
	p := tea.NewProgram(NewModel())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
