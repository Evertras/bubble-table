package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

const (
	columnKeyName              = "name"
	columnKeyElement           = "element"
	columnKeyConversations     = "convos"
	columnKeyPositiveSentiment = "positive"
	columnKeyNegativeSentiment = "negative"

	colorNormal   = "#fa0"
	colorFire     = "#f64"
	colorElectric = "#ff0"
	colorWater    = "#44f"
	colorPlant    = "#8b8"
)

var (
	styleSubtle = lipgloss.NewStyle().Foreground(lipgloss.Color("#888"))

	styleBase = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#a7a")).
			BorderForeground(lipgloss.Color("#a38")).
			Align(lipgloss.Right)
)

type Model struct {
	pokeTable table.Model
}

func makeRow(name, element, colorStr string, numConversations int, positiveSentiment, negativeSentiment float32) table.Row {
	return table.NewRow(table.RowData{
		columnKeyName:              name,
		columnKeyElement:           table.NewStyledCell(element, lipgloss.NewStyle().Foreground(lipgloss.Color(colorStr))),
		columnKeyConversations:     numConversations,
		columnKeyPositiveSentiment: positiveSentiment,
		columnKeyNegativeSentiment: negativeSentiment,
	})
}

func NewModel() Model {
	return Model{
		pokeTable: table.New([]table.Column{
			table.NewColumn(columnKeyName, "Name", 13),
			table.NewColumn(columnKeyElement, "Element", 10),
			table.NewColumn(columnKeyConversations, "# Conversations", 15),
			table.NewColumn(columnKeyPositiveSentiment, ":D %", 5).WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#8c8"))),
			table.NewColumn(columnKeyNegativeSentiment, ":( %", 5).WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#c88"))),
		}).WithRows([]table.Row{
			makeRow("Pikachu", "Electric", colorElectric, 2300648, 21.9, 8.54),
			makeRow("Eevee", "Normal", colorNormal, 636373, 26.4, 7.37),
			makeRow("Bulbasaur", "Plant", colorPlant, 352190, 25.7, 9.02),
			makeRow("Squirtle", "Water", colorWater, 241259, 25.6, 5.96),
			makeRow("Blastoise", "Water", colorWater, 162794, 19.5, 6.04),
			makeRow("Charmander", "Fire", colorFire, 265760, 31.2, 5.25),
			makeRow("Charizard", "Fire", colorFire, 567763, 25.6, 7.56),
		}).
			BorderRounded().
			WithBaseStyle(styleBase).
			WithPageSize(6).
			SortByDesc(columnKeyConversations).
			Focused(true),
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

	m.pokeTable, cmd = m.pokeTable.Update(msg)
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
	selected := m.pokeTable.HighlightedRow().Data[columnKeyName].(string)
	view := lipgloss.JoinVertical(
		lipgloss.Left,
		styleSubtle.Render("Press q or ctrl+c to quit - Sorted by # Conversations"),
		styleSubtle.Render("Highlighted: "+selected),
		styleSubtle.Render("https://www.nintendolife.com/news/2021/11/these-are-the-most-loved-and-most-hated-pokemon-according-to-a-new-study"),
		m.pokeTable.View(),
	) + "\n"

	return lipgloss.NewStyle().MarginLeft(1).Render(view)
}

func main() {
	p := tea.NewProgram(NewModel())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
