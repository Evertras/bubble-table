package main

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
)

const (
	columnKeyID = "id"

	numCols = 100
	numRows = 10
)

type Model struct {
	scrollableTable table.Model
}

func colKey(colNum int) string {
	return fmt.Sprintf("%d", colNum)
}

func genRow(id int) table.Row {
	data := table.RowData{
		columnKeyID: fmt.Sprintf("ID %d", id),
	}

	for i := 0; i < numCols; i++ {
		data[colKey(i)] = i + 1
	}

	return table.NewRow(data)
}

func NewModel() Model {
	rows := []table.Row{}

	for i := 0; i < numRows; i++ {
		rows = append(rows, genRow(i))
	}

	cols := []table.Column{
		table.NewColumn(columnKeyID, "ID", 5),
	}

	for i := 0; i < numCols; i++ {
		cols = append(cols, table.NewColumn(colKey(i), colKey(i+1), 5))
	}

	return Model{
		scrollableTable: table.New(cols).WithRows(rows).WithMaxTotalWidth(30).Focused(true),
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

	m.scrollableTable, cmd = m.scrollableTable.Update(msg)
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

	body.WriteString("A scrollable table\nPress shift+left or shift+right to scroll\nPress q or ctrl+c to quit\n\n")

	body.WriteString(m.scrollableTable.View())

	return body.String()
}

func main() {
	p := tea.NewProgram(NewModel())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
