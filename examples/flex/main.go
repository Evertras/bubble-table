package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

const (
	columnKeyName        = "name"
	columnKeyElement     = "element"
	columnKeyDescription = "description"

	minWidth = 30
)

type Model struct {
	flexTable   table.Model
	totalMargin int
	totalWidth  int
}

func NewModel() Model {
	return Model{
		flexTable: table.New([]table.Column{
			table.NewColumn(columnKeyName, "Name", 10),
			table.NewFlexColumn(columnKeyElement, "Element", 1),
			table.NewFlexColumn(columnKeyDescription, "Description", 3),
		}).WithRows([]table.Row{
			table.NewRow(table.RowData{
				columnKeyName:        "Pikachu",
				columnKeyElement:     "Electric",
				columnKeyDescription: "Super zappy mouse, handle with care",
			}),
			table.NewRow(table.RowData{
				columnKeyName:    "Charmander",
				columnKeyElement: "Fire",
				// TODO: Fix double width string length limiting!
				//columnKeyDescription: "直立した恐竜のような身体と、尻尾の先端に常に燃えている炎が特徴。",
				columnKeyDescription: "Lots of fire",
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

	m.flexTable, cmd = m.flexTable.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			cmds = append(cmds, tea.Quit)

		case "left":
			if m.totalWidth-m.totalMargin > minWidth {
				m.totalMargin++
				m.recalculateTable()
			}

		case "right":
			if m.totalMargin > 0 {
				m.totalMargin--
				m.recalculateTable()
			}
		}

	case tea.WindowSizeMsg:
		m.totalWidth = msg.Width

		m.recalculateTable()
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) recalculateTable() {
	m.flexTable = m.flexTable.WithTargetWidth(m.totalWidth - m.totalMargin)
}

func (m Model) View() string {
	strs := []string{
		"A flexible table that fills available space (Name is fixed-width)",
		fmt.Sprintf("Total margin: %d (left/right to adjust)", m.totalMargin),
		"Press q or ctrl+c to quit",
		m.flexTable.View(),
	}

	return lipgloss.JoinVertical(lipgloss.Left, strs...) + "\n"
}

func main() {
	p := tea.NewProgram(NewModel())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
