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

func rowStyleFunc(input table.RowStyleFuncInput) lipgloss.Style {
	calculatedStyle := lipgloss.NewStyle()

	switch input.Row.Data[columnKeyStatus] {
	case "Critical":
		calculatedStyle = styleCritical.Copy()
	case "Stable":
		calculatedStyle = styleStable.Copy()
	case "Good":
		calculatedStyle = styleGood.Copy()
	}

	if input.Index%2 == 0 {
		calculatedStyle = calculatedStyle.Background(lipgloss.Color("#222"))
	} else {
		calculatedStyle = calculatedStyle.Background(lipgloss.Color("#444"))
	}

	return calculatedStyle
}

func NewModel() Model {
	return Model{
		table:       table.New(generateColumns(0)).WithRowStyleFunc(rowStyleFunc),
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

// Generate columns based on how many are critical to show some summary
func generateColumns(numCritical int) []table.Column {
	// Show how many critical there are
	statusStr := fmt.Sprintf("Score (%d)", numCritical)
	statusCol := table.NewColumn(columnKeyStatus, statusStr, 10)

	if numCritical > 3 {
		// This normally applies the critical style to everything in the column,
		// but in this case we apply a row style which overrides it anyway.
		statusCol = statusCol.WithStyle(styleCritical)
	}

	return []table.Column{
		table.NewColumn(columnKeyID, "ID", 10),
		table.NewColumn(columnKeyScore, "Score", 8),
		statusCol,
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

		numCritical := 0

		for _, d := range msg {
			if d.Status == "Critical" {
				numCritical++
			}
		}

		// Reapply the new data and the new columns based on critical count
		m.table = m.table.WithRows(generateRowsFromData(m.data)).WithColumns(generateColumns(numCritical))

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
