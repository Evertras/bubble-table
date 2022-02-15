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
	table3x3 table.Model
	table1x3 table.Model
	table3x1 table.Model
	table1x1 table.Model
	table5x5 table.Model
}

func genTable(columnCount int, rowCount int) table.Model {
	columns := []table.Column{}

	for column := 0; column < columnCount; column++ {
		columnStr := fmt.Sprintf("%d", column+1)
		columns = append(columns, table.NewColumn(columnStr, columnStr, 4))
	}

	rows := []table.Row{}

	for row := 1; row < rowCount; row++ {
		rowData := table.RowData{}

		for column := 0; column < columnCount; column++ {
			columnStr := fmt.Sprintf("%d", column+1)
			rowData[columnStr] = fmt.Sprintf("%d,%d", column+1, row+1)
		}

		rows = append(rows, table.NewRow(rowData))
	}

	return table.New(columns).WithRows(rows).HeaderStyle(lipgloss.NewStyle().Bold(true))
}

func NewModel() Model {
	return Model{
		table1x1: genTable(1, 1),
		table3x1: genTable(3, 1),
		table1x3: genTable(1, 3),
		table3x3: genTable(3, 3),
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

	m.table3x1, cmd = m.table3x1.Update(msg)
	cmds = append(cmds, cmd)

	m.table1x3, cmd = m.table1x3.Update(msg)
	cmds = append(cmds, cmd)

	m.table3x3, cmd = m.table3x3.Update(msg)
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

	body.WriteString("Table demo with various sized tables!\nPress q or ctrl+c to quit\n")

	pad := lipgloss.NewStyle().Padding(1)

	tablesSmall := lipgloss.JoinHorizontal(
		lipgloss.Top,
		pad.Render(m.table1x1.View()),
		pad.Render(m.table1x3.View()),
		pad.Render(m.table3x1.View()),
		pad.Render(m.table3x3.View()),
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
