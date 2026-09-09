package main

import (
	"log"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/evertras/bubble-table/table"
)

const (
	columnKeyName   = "name"
	columnKeyType   = "type"
	columnKeyWins   = "wins"
	columnKeyLevel  = "level"
	columnKeyRarity = "rarity"
)

// rarityRank gives each rarity label its collector order. Neither plain
// string sort nor SortTypeNatural could express this: alphabetically
// "Common" < "Legendary" < "Rare" < "Uncommon", which isn't the order
// anyone actually wants. A custom comparator via WithSortFunc is the only
// way to sort by an arbitrary, domain-specific ordering like this.
var rarityRank = map[string]int{
	"Common":    0,
	"Uncommon":  1,
	"Rare":      2,
	"Legendary": 3,
}

type Model struct {
	simpleTable table.Model

	columnSortKey string
	sortDirection string
}

func NewModel() Model {
	return Model{
		simpleTable: table.New([]table.Column{
			table.NewColumn(columnKeyName, "Name", 13),
			table.NewColumn(columnKeyType, "Type", 13),
			table.NewColumn(columnKeyWins, "Win %", 8).
				WithFormatString("%.1f%%"),
			// A plain string sort would order this "Level 1", "Level 10", "Level 100",
			// "Level 2", ... since it compares character by character. SortTypeNatural
			// compares the embedded numbers by value instead, so it sorts the way a
			// human would expect.
			table.NewColumn(columnKeyLevel, "Level", 10).
				WithSortType(table.SortTypeNatural),
			table.NewColumn(columnKeyRarity, "Rarity", 10).
				WithSortFunc(func(a, b any) bool {
					aStr, _ := a.(string)
					bStr, _ := b.(string)

					return rarityRank[aStr] < rarityRank[bStr]
				}),
		}).WithRows([]table.Row{
			table.NewRow(table.RowData{
				columnKeyName:   "ピカピカ",
				columnKeyType:   "Pikachu",
				columnKeyWins:   78.3,
				columnKeyLevel:  "Level 1",
				columnKeyRarity: "Uncommon",
			}),
			table.NewRow(table.RowData{
				columnKeyName:   "Zapmouse",
				columnKeyType:   "Pikachu",
				columnKeyWins:   3.3,
				columnKeyLevel:  "Level 10",
				columnKeyRarity: "Common",
			}),
			table.NewRow(table.RowData{
				columnKeyName:   "Burninator",
				columnKeyType:   "Charmander",
				columnKeyWins:   32.1,
				columnKeyLevel:  "Level 2",
				columnKeyRarity: "Legendary",
			}),
			table.NewRow(table.RowData{
				columnKeyName:   "Alphonse",
				columnKeyType:   "Pikachu",
				columnKeyWins:   13.8,
				columnKeyLevel:  "Level 20",
				columnKeyRarity: "Rare",
			}),
			table.NewRow(table.RowData{
				columnKeyName:   "Trogdor",
				columnKeyType:   "Charmander",
				columnKeyWins:   99.9,
				columnKeyLevel:  "Level 3",
				columnKeyRarity: "Legendary",
			}),
			table.NewRow(table.RowData{
				columnKeyName:   "Dihydrogen Monoxide",
				columnKeyType:   "Squirtle",
				columnKeyWins:   31.348,
				columnKeyLevel:  "Level 100",
				columnKeyRarity: "Common",
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

	m.simpleTable, cmd = m.simpleTable.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			cmds = append(cmds, tea.Quit)

		case "n":
			m.columnSortKey = columnKeyName
			m.simpleTable = m.simpleTable.SortByAsc(m.columnSortKey)

		case "t":
			m.columnSortKey = columnKeyType
			// Within the same type, order each by wins
			m.simpleTable = m.simpleTable.SortByAsc(m.columnSortKey).ThenSortByDesc(columnKeyWins)

		case "w":
			m.columnSortKey = columnKeyWins
			m.simpleTable = m.simpleTable.SortByDesc(m.columnSortKey)

		case "l":
			m.columnSortKey = columnKeyLevel
			m.simpleTable = m.simpleTable.SortByAsc(m.columnSortKey)

		case "r":
			m.columnSortKey = columnKeyRarity
			m.simpleTable = m.simpleTable.SortByAsc(m.columnSortKey)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() tea.View {
	body := strings.Builder{}

	body.WriteString("A sorted simple default table\n")
	body.WriteString("Sort by (n)ame, (t)ype->wins combo, (w)ins, (l)evel (natural sort), or (r)arity (custom sort func)\n")
	body.WriteString("Currently sorting by: " + m.columnSortKey + "\n")
	body.WriteString("Press q or ctrl+c to quit\n\n")

	body.WriteString(m.simpleTable.View())

	return tea.NewView(body.String())
}

func main() {
	p := tea.NewProgram(NewModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
