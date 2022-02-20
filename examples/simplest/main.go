package main

import (
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
)

const (
	columnKeyName    = "name"
	columnKeyElement = "element"
)

type Model struct {
	simpleTable table.Model
}

func NewModel() Model {
	return Model{
		simpleTable: table.New([]table.Column{
			table.NewColumn(columnKeyName, "Name", 13),
			table.NewColumn(columnKeyElement, "Element", 10),
		}).WithRows([]table.Row{
			table.NewRow(table.RowData{
				columnKeyName:    "Pikachu",
				columnKeyElement: "Electric",
			}),
			table.NewRow(table.RowData{
				columnKeyName:    "Charmander",
				columnKeyElement: "Fire",
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
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	body := strings.Builder{}

	body.WriteString("A very simple default table (non-interactive)\nPress q or ctrl+c to quit\n\n")

	body.WriteString(m.simpleTable.View())

	return body.String()
}

func main() {
	p := tea.NewProgram(NewModel())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
