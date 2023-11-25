package main

import (
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

const (
	columnKeyName     = "name"
	columnKeyCountry  = "country"
	columnKeyCurrency = "crurrency"
)

type Model struct {
	tableModel table.Model
}

func NewModel() Model {
	columns := []table.Column{
		table.NewColumn(columnKeyName, "Name", 10).WithStyle(
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("#88f")),
		),
		table.NewColumn(columnKeyCountry, "Country", 20),
		table.NewColumn(columnKeyCurrency, "Currency", 10),
	}

	rows := []table.Row{
		table.NewRow(
			table.RowData{
				columnKeyName:     "Talon Stokes",
				columnKeyCountry:  "Mexico",
				columnKeyCurrency: "$23.17",
			}),
		table.NewRow(
			table.RowData{
				columnKeyName:     "Sonia Shepard",
				columnKeyCountry:  "United States",
				columnKeyCurrency: "$76.47",
			}),
		table.NewRow(
			table.RowData{
				columnKeyName:     "Shad Reed",
				columnKeyCountry:  "Turkey",
				columnKeyCurrency: "$62.99",
			}),
		table.NewRow(
			table.RowData{
				columnKeyName:     "Kibo Clay",
				columnKeyCountry:  "Philippines",
				columnKeyCurrency: "$29.82",
			}),
		table.NewRow(
			table.RowData{

				columnKeyName:     "Leslie Kerr",
				columnKeyCountry:  "Singapore",
				columnKeyCurrency: "$70.54",
			}),
		table.NewRow(
			table.RowData{
				columnKeyName:     "Micah Hurst",
				columnKeyCountry:  "Pakistan",
				columnKeyCurrency: "$80.84",
			}),
		table.NewRow(
			table.RowData{
				columnKeyName:     "Dora Miranda",
				columnKeyCountry:  "Colombia",
				columnKeyCurrency: "$34.75",
			}),
		table.NewRow(
			table.RowData{
				columnKeyName:     "Keefe Walters",
				columnKeyCountry:  "China",
				columnKeyCurrency: "$56.82",
			}),
		table.NewRow(
			table.RowData{
				columnKeyName:     "Fujimoto Tarokizaemon no shoutokinori",
				columnKeyCountry:  "Japan",
				columnKeyCurrency: "$89.31",
			}),
		table.NewRow(
			table.RowData{
				columnKeyName:     "Keefe Walters",
				columnKeyCountry:  "China",
				columnKeyCurrency: "$56.82",
			}),
		table.NewRow(
			table.RowData{
				columnKeyName:     "Vincent Sanchez",
				columnKeyCountry:  "Peru",
				columnKeyCurrency: "$71.60",
			}),
		table.NewRow(
			table.RowData{
				columnKeyName:     "Lani Figueroa",
				columnKeyCountry:  "United Kingdom",
				columnKeyCurrency: "$90.67",
			}),
	}

	model := Model{
		tableModel: table.New(columns).
			WithRows(rows).
			HeaderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)).
			Focused(true).
			WithBaseStyle(
				lipgloss.NewStyle().
					BorderForeground(lipgloss.Color("#a38")).
					Foreground(lipgloss.Color("#a7a")).
					Align(lipgloss.Left),
			).
			WithMultiline(true),
	}

	return model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.tableModel, cmd = m.tableModel.Update(msg)
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

	body.WriteString("A table demo with multiline feature enabled!\n")
	body.WriteString("Press up/down or j/k to move around\n")
	body.WriteString(m.tableModel.View())
	body.WriteString("\n")

	return body.String()
}

func main() {
	p := tea.NewProgram(NewModel())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
