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
	tableDefault        table.Model
	tableWithRowIndices table.Model

	rowCount int
}

func genRows(columnCount int, rowCount int) []table.Row {
	rows := []table.Row{}

	for row := 1; row <= rowCount; row++ {
		rowData := table.RowData{}

		for column := 0; column < columnCount; column++ {
			columnStr := fmt.Sprintf("%d", column+1)
			rowData[columnStr] = fmt.Sprintf("%d - %d", column+1, row)
		}

		rows = append(rows, table.NewRow(rowData))
	}

	return rows
}

func genTable(columnCount int, rowCount int) table.Model {
	columns := []table.Column{}

	for column := 0; column < columnCount; column++ {
		columnStr := fmt.Sprintf("%d", column+1)
		columns = append(columns, table.NewColumn(columnStr, columnStr, 8))
	}

	rows := genRows(columnCount, rowCount)

	return table.New(columns).WithRows(rows).HeaderStyle(lipgloss.NewStyle().Bold(true))
}

func NewModel() Model {
	const startingRowCount = 105

	m := Model{
		rowCount:            startingRowCount,
		tableDefault:        genTable(3, startingRowCount).WithPageSize(10).Focused(true),
		tableWithRowIndices: genTable(3, startingRowCount).WithPageSize(10).Focused(false),
	}

	m.regenTableRows()

	return m
}

func (m *Model) regenTableRows() {
	m.tableDefault = m.tableDefault.WithRows(genRows(3, m.rowCount))
	m.tableWithRowIndices = m.tableWithRowIndices.WithRows(genRows(3, m.rowCount))
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			cmds = append(cmds, tea.Quit)

		case "a":
			m.tableDefault = m.tableDefault.Focused(true)
			m.tableWithRowIndices = m.tableWithRowIndices.Focused(false)

		case "b":
			m.tableDefault = m.tableDefault.Focused(false)
			m.tableWithRowIndices = m.tableWithRowIndices.Focused(true)

		case "z":
			if m.rowCount < 10 {
				break
			}

			m.rowCount -= 10
			m.regenTableRows()

		case "x":
			m.rowCount += 10
			m.regenTableRows()
		}
	}

	m.tableDefault, cmd = m.tableDefault.Update(msg)
	cmds = append(cmds, cmd)

	m.tableWithRowIndices, cmd = m.tableWithRowIndices.Update(msg)
	cmds = append(cmds, cmd)

	// Write a custom footer
	start, end := m.tableWithRowIndices.VisibleIndices()
	m.tableWithRowIndices = m.tableWithRowIndices.WithStaticFooter(
		fmt.Sprintf("%d-%d of %d", start+1, end+1, m.tableWithRowIndices.TotalRows()),
	)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	body := strings.Builder{}

	body.WriteString("Table demo with pagination! Press left/right to move pages, or use page up/down\nPress 'a' for left table, 'b' for right table\nPress 'z' to reduce rows by 10, 'y' to increase rows by 10\nPress q or ctrl+c to quit\n\n")

	pad := lipgloss.NewStyle().Padding(1)

	tables := []string{
		lipgloss.JoinVertical(lipgloss.Center, "A", pad.Render(m.tableDefault.View())),
		lipgloss.JoinVertical(lipgloss.Center, "B", pad.Render(m.tableWithRowIndices.View())),
	}

	body.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, tables...))

	return body.String()
}

func main() {
	p := tea.NewProgram(NewModel())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
