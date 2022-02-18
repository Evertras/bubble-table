// This file contains a full demo of most available features, for both testing
// and for reference
package main

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

const (
	columnKeyID          = "id"
	columnKeyName        = "name"
	columnKeyDescription = "description"
	columnKeyCount       = "count"
)

var (
	customBorder = table.Border{
		Top:    "─",
		Left:   "│",
		Right:  "│",
		Bottom: "─",

		TopRight:    "╮",
		TopLeft:     "╭",
		BottomRight: "╯",
		BottomLeft:  "╰",

		TopJunction:    "╥",
		LeftJunction:   "├",
		RightJunction:  "┤",
		BottomJunction: "╨",
		InnerJunction:  "╫",

		InnerDivider: "║",
	}
)

type Model struct {
	tableModel table.Model
}

func NewModel() Model {
	columns := []table.Column{
		table.NewColumn(columnKeyID, "ID", 5),
		table.NewColumn(columnKeyName, "Name", 10),
		table.NewColumn(columnKeyDescription, "Description", 30),
		table.NewColumn(columnKeyCount, "#", 5),
	}

	rows := []table.Row{
		table.NewRow(table.RowData{
			columnKeyID:          "abc",
			columnKeyName:        "Hello",
			columnKeyDescription: "The first table entry, ever",
			columnKeyCount:       4,
		}),
		table.NewRow(table.RowData{
			columnKeyID:          "123",
			columnKeyName:        "Oh no",
			columnKeyDescription: "Super bold!",
			columnKeyCount:       17,
		}).WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)),
		table.NewRow(table.RowData{
			columnKeyID:          "def",
			columnKeyName:        "Yay",
			columnKeyDescription: "This is a really, really, really long description that will get cut off",
			columnKeyCount:       "N/A",
		}),
	}

	// Start with the default key map and change it slightly, just for demoing
	keys := table.DefaultKeyMap()
	keys.RowDown.SetKeys("j", "down", "s")
	keys.RowUp.SetKeys("j", "down", "w")

	return Model{
		tableModel: table.New(columns).
			WithRows(rows).
			HeaderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)).
			SelectableRows(true).
			Focused(true).
			Border(customBorder).
			WithKeyMap(keys),
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

	highlightedRow := m.tableModel.HighlightedRow()
	body.WriteString("Table demo with all features enabled!\nPress space/enter to select a row, q or ctrl+c to quit\n")

	body.WriteString(fmt.Sprintf("Currently looking at ID: %s\n", highlightedRow.Data[columnKeyID]))

	selectedIDs := []string{}

	for _, row := range m.tableModel.SelectedRows() {
		// Slightly dangerous type assumption but fine for demo
		selectedIDs = append(selectedIDs, row.Data[columnKeyID].(string))
	}

	body.WriteString(fmt.Sprintf("SelectedIDs: %s\n", strings.Join(selectedIDs, ", ")))

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
