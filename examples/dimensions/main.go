package main

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

type Model struct {
	table2x2 table.Model
	table1x2 table.Model
	table2x1 table.Model
	table1x1 table.Model
	table5x5 table.Model
}

func genTable(columnCount int, rowCount int) table.Model {
	headers := []table.Header{}

	for column := 0; column < columnCount; column++ {
		columnStr := fmt.Sprintf("%d", column+1)
		headers = append(headers, table.NewHeader(columnStr, columnStr, 4))
	}

	rows := []table.Row{}

	for row := 0; row < rowCount; row++ {
		rowData := table.RowData{}

		for column := 0; column < columnCount; column++ {
			columnStr := fmt.Sprintf("%d", column+1)
			rowData[columnStr] = fmt.Sprintf("%d,%d", column+1, row+1)
		}

		rows = append(rows, table.NewRow(rowData))
	}

	return table.New(headers).WithRows(rows).HeaderStyle(lipgloss.NewStyle().Bold(true))
}

func NewModel() Model {
	return Model{
		table1x1: genTable(1, 1),
		table2x1: genTable(2, 1),
		table1x2: genTable(1, 2),
		table2x2: genTable(2, 2),
		table5x5: genTable(5, 5),
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

	m.table1x1, cmd = m.table1x1.Update(msg)
	cmds = append(cmds, cmd)

	m.table2x1, cmd = m.table2x1.Update(msg)
	cmds = append(cmds, cmd)

	m.table1x2, cmd = m.table1x2.Update(msg)
	cmds = append(cmds, cmd)

	m.table2x2, cmd = m.table2x2.Update(msg)
	cmds = append(cmds, cmd)

	m.table5x5, cmd = m.table5x5.Update(msg)
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

	body.WriteString("Table demo with various sized tables!\nPress space/enter to select a row, q or ctrl+c to quit\n")

	pad := lipgloss.NewStyle().Padding(1)

	tablesSmall := lipgloss.JoinHorizontal(
		lipgloss.Top,
		pad.Render(m.table1x1.View()),
		pad.Render(m.table1x2.View()),
		pad.Render(m.table2x1.View()),
		pad.Render(m.table2x2.View()),
	)

	tableBig := pad.Render(m.table5x5.View())

	body.WriteString(lipgloss.JoinVertical(lipgloss.Center, tablesSmall, tableBig))

	return body.String()
}

func main() {
	p := tea.NewProgram(NewModel())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
