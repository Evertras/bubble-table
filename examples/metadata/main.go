// This is a more data-driven example of the Pokemon table
package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

type Element string

const (
	columnKeyName              = "name"
	columnKeyElement           = "element"
	columnKeyConversations     = "convos"
	columnKeyPositiveSentiment = "positive"
	columnKeyNegativeSentiment = "negative"

	// This is not a visible column, but is used to attach useful reference data
	// to the row itself for easier retrieval
	columnKeyPokemonData = "pokedata"

	elementNormal   Element = "Normal"
	elementFire     Element = "Fire"
	elementElectric Element = "Electric"
	elementWater    Element = "Water"
	elementPlant    Element = "Plant"
)

var (
	styleSubtle = lipgloss.NewStyle().Foreground(lipgloss.Color("#888"))

	styleBase = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#a7a")).
			BorderForeground(lipgloss.Color("#a38")).
			Align(lipgloss.Right)

	elementColors = map[Element]string{
		elementNormal:   "#fa0",
		elementFire:     "#f64",
		elementElectric: "#ff0",
		elementWater:    "#44f",
		elementPlant:    "#8b8",
	}
)

type Pokemon struct {
	Name                     string
	Element                  Element
	ConversationCount        int
	PositiveSentimentPercent float32
	NegativeSentimentPercent float32
}

func NewPokemon(name string, element Element, conversationCount int, positiveSentimentPercent float32, negativeSentimentPercent float32) Pokemon {
	return Pokemon{
		Name:                     name,
		Element:                  element,
		ConversationCount:        conversationCount,
		PositiveSentimentPercent: positiveSentimentPercent,
		NegativeSentimentPercent: negativeSentimentPercent,
	}
}

func (p Pokemon) ToRow() table.Row {
	color, exists := elementColors[p.Element]

	if !exists {
		color = elementColors[elementNormal]
	}

	return table.NewRow(table.RowData{
		columnKeyName:              p.Name,
		columnKeyElement:           table.NewStyledCell(p.Element, lipgloss.NewStyle().Foreground(lipgloss.Color(color))),
		columnKeyConversations:     p.ConversationCount,
		columnKeyPositiveSentiment: p.PositiveSentimentPercent,
		columnKeyNegativeSentiment: p.NegativeSentimentPercent,

		// This isn't a visible column, but we can add the data here anyway for later retrieval
		columnKeyPokemonData: p,
	})
}

type Model struct {
	pokeTable table.Model
}

func NewModel() Model {
	pokemon := []Pokemon{
		NewPokemon("Pikachu", elementElectric, 2300648, 21.9, 8.54),
		NewPokemon("Eevee", elementNormal, 636373, 26.4, 7.37),
		NewPokemon("Bulbasaur", elementPlant, 352190, 25.7, 9.02),
		NewPokemon("Squirtle", elementWater, 241259, 25.6, 5.96),
		NewPokemon("Blastoise", elementWater, 162794, 19.5, 6.04),
		NewPokemon("Charmander", elementFire, 265760, 31.2, 5.25),
		NewPokemon("Charizard", elementFire, 567763, 25.6, 7.56),
	}

	rows := []table.Row{}

	for _, p := range pokemon {
		rows = append(rows, p.ToRow())
	}

	return Model{
		pokeTable: table.New([]table.Column{
			table.NewColumn(columnKeyName, "Name", 13),
			table.NewColumn(columnKeyElement, "Element", 10),
			table.NewColumn(columnKeyConversations, "# Conversations", 15),
			table.NewColumn(columnKeyPositiveSentiment, ":D %", 5).WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#8c8"))),
			table.NewColumn(columnKeyNegativeSentiment, ":( %", 5).WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#c88"))),
		}).WithRows(rows).
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
	// Get the metadata back out of the row
	selected := m.pokeTable.HighlightedRow().Data[columnKeyPokemonData].(Pokemon)

	view := lipgloss.JoinVertical(
		lipgloss.Left,
		styleSubtle.Render("Press q or ctrl+c to quit - Sorted by # Conversations"),
		styleSubtle.Render("Highlighted: "+fmt.Sprintf("%s (%s)", selected.Name, selected.Element)),
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
