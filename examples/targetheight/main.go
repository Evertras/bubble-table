package main

import (
	"fmt"
	"log"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/evertras/bubble-table/table"
)

const (
	columnKeyID          = "id"
	columnKeyName        = "name"
	columnKeyDescription = "description"

	// The height we want the table to occupy, including borders/header/footer.
	targetHeight = 12
	columnWidth  = 30
)

// The first few rows have short Pokédex entries (single-line rows); the remainder
// have longer entries that wrap to multiple lines, demonstrating that WithTargetHeight
// keeps the table height consistent regardless of row content.
var pokeRows = []table.Row{
	// Short descriptions — single-line rows
	table.NewRow(table.RowData{columnKeyID: 1, columnKeyName: "Bulbasaur", columnKeyDescription: "A strange seed was planted on its back at birth."}),
	table.NewRow(table.RowData{columnKeyID: 4, columnKeyName: "Charmander", columnKeyDescription: "The flame on its tail shows its life force."}),
	table.NewRow(table.RowData{columnKeyID: 7, columnKeyName: "Squirtle", columnKeyDescription: "Shoots water at prey while in the water."}),
	table.NewRow(table.RowData{columnKeyID: 25, columnKeyName: "Pikachu", columnKeyDescription: "Has electric sacs on each cheek."}),
	table.NewRow(table.RowData{columnKeyID: 39, columnKeyName: "Jigglypuff", columnKeyDescription: "Uses its round eyes to entrance foes."}),
	// Longer Pokédex entries — wrap to multiple lines
	table.NewRow(table.RowData{columnKeyID: 143, columnKeyName: "Snorlax", columnKeyDescription: "Very lazy. It just eats and sleeps. As its rotund bulk builds, it becomes steadily more slothful."}),
	table.NewRow(table.RowData{columnKeyID: 131, columnKeyName: "Lapras", columnKeyDescription: "A gentle soul that can read the hearts of people. It can ferry people across the sea on its back."}),
	table.NewRow(table.RowData{columnKeyID: 147, columnKeyName: "Dratini", columnKeyDescription: "Long considered a mythical Pokémon until a fisherman landed a live specimen after hooking it."}),
	table.NewRow(table.RowData{columnKeyID: 137, columnKeyName: "Porygon", columnKeyDescription: "A Pokémon that consists entirely of programming code. Capable of moving freely in cyberspace."}),
	table.NewRow(table.RowData{columnKeyID: 133, columnKeyName: "Eevee", columnKeyDescription: "Its genetic code is irregular. It may mutate if it is exposed to radiation from element stones."}),
}

type Model struct {
	tableModel table.Model
}

func newTable(rows []table.Row) table.Model {
	columns := []table.Column{
		table.NewColumn(columnKeyID, "ID", 4),
		table.NewColumn(columnKeyName, "Name", 12),
		table.NewColumn(columnKeyDescription, "Pokédex Entry", columnWidth),
	}

	highlight := lipgloss.NewStyle().
		Foreground(lipgloss.Color("212")).
		Bold(true)

	return table.New(columns).
		WithRows(rows).
		WithMultiline(true).
		WithTargetHeight(targetHeight).
		WithMinimumHeight(targetHeight).
		HighlightStyle(highlight).
		Focused(true)
}

func NewModel() Model {
	return Model{
		tableModel: newTable(pokeRows),
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			cmds = append(cmds, tea.Quit)
		}
	}

	m.tableModel, cmd = m.tableModel.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// ruler draws a vertical ruler showing the target height boundary.
func ruler(height int) string {
	lines := make([]string, height)
	for i := range lines {
		if i == height-1 {
			lines[i] = fmt.Sprintf("%2d ◄── target bottom", i+1)
		} else {
			lines[i] = fmt.Sprintf("%2d │", i+1)
		}
	}
	return strings.Join(lines, "\n")
}

func (m Model) View() tea.View {
	body := strings.Builder{}

	body.WriteString("Pokédex — WithTargetHeight keeps the table at a fixed height across pages.\n")
	fmt.Fprintf(&body, "Short entries render as single lines; long entries wrap — target height: %d lines.\n\n", targetHeight)

	rulerStr := ruler(targetHeight)

	body.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, m.tableModel.View(), "  ", rulerStr))
	body.WriteString("\n\nPress q to quit")

	return tea.NewView(body.String())
}

func main() {
	p := tea.NewProgram(NewModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
