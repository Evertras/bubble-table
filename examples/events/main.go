package main

import (
	"fmt"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

type Element string

const (
	columnKeyName    = "name"
	columnKeyElement = "element"

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
		columnKeyName:    p.Name,
		columnKeyElement: table.NewStyledCell(p.Element, lipgloss.NewStyle().Foreground(lipgloss.Color(color))),

		// This isn't a visible column, but we can add the data here anyway for later retrieval
		columnKeyPokemonData: p,
	})
}

type Model struct {
	pokeTable table.Model

	currentPokemonData Pokemon

	lastSelectedEvent table.UserEventRowSelectToggled
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
		}).WithRows(rows).
			BorderRounded().
			WithBaseStyle(styleBase).
			WithPageSize(4).
			Focused(true).
			SelectableRows(true),
		currentPokemonData: pokemon[0],
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

	for _, e := range m.pokeTable.GetLastUpdateUserEvents() {
		switch e := e.(type) {
		case table.UserEventHighlightedIndexChanged:
			// We can pretend this is an async data retrieval, but really we already
			// have the data, so just return it after some fake delay.  Also note
			// that the event has some data attached to it, but we're ignoring
			// that for this example as we just want the current highlighted row.
			selectedPokemon := m.pokeTable.HighlightedRow().Data[columnKeyPokemonData].(Pokemon)

			cmds = append(cmds, func() tea.Msg {
				time.Sleep(time.Millisecond * 200)

				return selectedPokemon
			})

		case table.UserEventRowSelectToggled:
			m.lastSelectedEvent = e
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			cmds = append(cmds, tea.Quit)
		}

	case Pokemon:
		m.currentPokemonData = msg
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	view := lipgloss.JoinVertical(
		lipgloss.Left,
		styleSubtle.Render("Press q or ctrl+c to quit"),
		fmt.Sprintf("Highlighted (200 ms delay): %s (%s)", m.currentPokemonData.Name, m.currentPokemonData.Element),
		fmt.Sprintf("Last selected event: %d (%v)", m.lastSelectedEvent.RowIndex, m.lastSelectedEvent.IsSelected),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#8c8")).
			Render(":D %"+fmt.Sprintf("%.1f", m.currentPokemonData.PositiveSentimentPercent)),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#c88")).
			Render(":( %"+fmt.Sprintf("%.1f", m.currentPokemonData.NegativeSentimentPercent)),
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
