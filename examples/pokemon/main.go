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
	colorElectric = "#ff0"
	colorFire     = "#f64"
	colorPlant    = "#8b8"
	colorWater    = "#44f"
)

var (
	styleSubtle = lipgloss.NewStyle().Foreground(lipgloss.Color("#888"))

	styleBase = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#a7a")).
			BorderForeground(lipgloss.Color("#a38")).
			Align(lipgloss.Right)
)

type Model struct {
	pokeTable            table.Model
	favoriteElementIndex int
}

var elementList = []string{
	"Normal",
	"Electric",
	"Fire",
	"Plant",
	"Water",
}

var colorMap = map[any]string{
	"Electric": colorElectric,
	"Fire":     colorFire,
	"Plant":    colorPlant,
	"Water":    colorWater,
}

func makeRow(name, element string, numConversations int, positiveSentiment, negativeSentiment float32) table.Row {
	elementStyleFunc := func(input table.StyledCellFuncInput) lipgloss.Style {
		color := colorNormal

		if val, ok := colorMap[input.Data]; ok {
			color = val
		}

		style := lipgloss.NewStyle().Foreground(lipgloss.Color(color))

		if input.GlobalMetadata["favoriteElement"] == input.Data {
			style = style.Italic(true)
		}

		return style
	}

	return table.NewRow(table.RowData{
		columnKeyName:              name,
		columnKeyElement:           table.NewStyledCellWithStyleFunc(element, elementStyleFunc),
		columnKeyConversations:     numConversations,
		columnKeyPositiveSentiment: positiveSentiment,
		columnKeyNegativeSentiment: negativeSentiment,
	})
}

func genMetadata(favoriteElementIndex int) map[string]any {
	return map[string]any{
		"favoriteElement": elementList[favoriteElementIndex],
	}
}

func NewModel() Model {
	initialFavoriteElementIndex := 0
	return Model{
		favoriteElementIndex: initialFavoriteElementIndex,
		pokeTable: table.New([]table.Column{
			table.NewColumn(columnKeyName, "Name", 13),
			table.NewColumn(columnKeyElement, "Element", 10),
			table.NewColumn(columnKeyConversations, "# Conversations", 15),
			table.NewColumn(columnKeyPositiveSentiment, ":D %", 6).
				WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#8c8"))).
				WithFormatString("%.1f%%"),
			table.NewColumn(columnKeyNegativeSentiment, ":( %", 6).
				WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#c88"))).
				WithFormatString("%.1f%%"),
		}).WithRows([]table.Row{
			makeRow("Pikachu", "Electric", 2300648, 21.9, 8.54),
			makeRow("Eevee", "Normal", 636373, 26.4, 7.37),
			makeRow("Bulbasaur", "Plant", 352190, 25.7, 9.02),
			makeRow("Squirtle", "Water", 241259, 25.6, 5.96),
			makeRow("Blastoise", "Water", 162794, 19.5, 6.04),
			makeRow("Charmander", "Fire", 265760, 31.2, 5.25),
			makeRow("Charizard", "Fire", 567763, 25.6, 7.56),
		}).
			BorderRounded().
			WithBaseStyle(styleBase).
			WithPageSize(6).
			SortByDesc(columnKeyConversations).
			Focused(true).
			WithGlobalMetadata(genMetadata(initialFavoriteElementIndex)),
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

		case "e":
			m.favoriteElementIndex++
			if m.favoriteElementIndex >= len(elementList) {
				m.favoriteElementIndex = 0
			}

			m.pokeTable = m.pokeTable.WithGlobalMetadata(genMetadata(m.favoriteElementIndex))
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
		styleSubtle.Render("Favorite element: "+elementList[m.favoriteElementIndex]),
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
