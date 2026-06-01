package main

import (
	"fmt"
	"log"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/evertras/bubble-table/table"
)

func hyperlink(text string, url string) string {
	return "\x1b]8;;" + url + "\x07" + text + "\x1b]8;;\x07"
}

const (
	columnKeyName    = "name"
	columnKeyElement = "element"
	columnKeyUrl     = "url"
)

type Model struct {
	simpleTable table.Model
}

func NewModel() Model {
	return Model{
		simpleTable: table.New([]table.Column{
			table.NewColumn(columnKeyName, "Name", 13),
			table.NewColumn(columnKeyElement, "Element", 10),
			table.NewColumn(columnKeyUrl, "Url", 30),
		}).WithRows([]table.Row{
			table.NewRow(table.RowData{
				columnKeyName:    "Pikachu",
				columnKeyElement: "Electric",
				columnKeyUrl:     hyperlink("See Pokedex", "https://www.serebii.net/pokemon/pikachu/"),
			}),
			table.NewRow(table.RowData{
				columnKeyName:    "Charmander",
				columnKeyElement: "Fire",
				columnKeyUrl:     hyperlink("See Pokedex", "https://www.serebii.net/pokemon/charmander/"),
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
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			cmds = append(cmds, tea.Quit)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() tea.View {
	body := strings.Builder{}

	body.WriteString("A very simple default table (non-interactive)\nPress q or ctrl+c to quit\n\n")

	body.WriteString(m.simpleTable.View())

	return tea.NewView(body.String())
}

func main() {
	fmt.Println("This is a URL in the terminal: " + hyperlink("Go to Serebii", "https://www.serebii.net/pokemon/charmander/"))

	p := tea.NewProgram(NewModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
