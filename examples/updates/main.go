package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

const (
	columnKeyID     = "id"
	columnKeyScore  = "score"
	columnKeyStatus = "status"
)

var (
	styleCritical = lipgloss.NewStyle().Foreground(lipgloss.Color("#f00"))
	styleStable   = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0"))
	styleGood     = lipgloss.NewStyle().Foreground(lipgloss.Color("#0f0"))
)

type Model struct {
	table table.Model

	updateDelay time.Duration

	data []*SomeData
}

func NewModel() Model {
	columns := []table.Column{
		table.NewColumn(columnKeyID, "ID", 10),
		table.NewColumn(columnKeyScore, "Score", 8),
		table.NewColumn(columnKeyStatus, "Status", 10),
	}

	return Model{
		table:       table.New(columns),
		updateDelay: time.Second,
	}
}

// This data is stored somewhere else, maybe on a client or some other thing
func refreshDataCmd() tea.Msg {
	// This could come from some API or something
	return []*SomeData{
		NewSomeData("abc"),
		NewSomeData("def"),
		NewSomeData("123"),
		NewSomeData("ok"),
		NewSomeData("another"),
		NewSomeData("yay"),
		NewSomeData("more"),
	}
}

func (m Model) Init() tea.Cmd {
	return refreshDataCmd
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			cmds = append(cmds, tea.Quit)

		case "up":
			if m.updateDelay < time.Second {
				m.updateDelay *= 10
			}

		case "down":
			if m.updateDelay > time.Millisecond*1 {
				m.updateDelay /= 10
			}
		}

	case []*SomeData:
		m.data = msg

		// Reapply the new data
		m.table = m.table.WithRows(generateRowsFromData(m.data))

		// This can be from any source, but for demo purposes let's party!
		delay := m.updateDelay
		cmds = append(cmds, func() tea.Msg {
			time.Sleep(delay)
			return refreshDataCmd()
		})
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	body := strings.Builder{}

	body.WriteString(
		fmt.Sprintf(
			"Table demo with updating data!  Updating every %v\nPress up/down to update faster/slower\nPress q or ctrl+c to quit\n",
			m.updateDelay,
		))

	pad := lipgloss.NewStyle().Padding(1)

	body.WriteString(pad.Render(m.table.View()))

	return body.String()
}

func generateRowsFromData(data []*SomeData) []table.Row {
	rows := []table.Row{}

	for _, entry := range data {
		row := table.NewRow(table.RowData{
			columnKeyID:     entry.ID,
			columnKeyScore:  entry.Score,
			columnKeyStatus: entry.Status,
		})

		// Highlight different statuses
		switch entry.Status {
		case "Critical":
			row = row.WithStyle(styleCritical)

		case "Stable":
			row = row.WithStyle(styleStable)

		case "Good":
			row = row.WithStyle(styleGood)
		}

		rows = append(rows, row)
	}

	return rows
}

func main() {

	p := tea.NewProgram(NewModel())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
