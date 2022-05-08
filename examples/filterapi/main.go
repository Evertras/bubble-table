package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
)

const (
	columnKeyTitle       = "title"
	columnKeyAuthor      = "author"
	columnKeyDescription = "description"
)

type Model struct {
	table           table.Model
	filterTextInput textinput.Model
}

func NewModel() Model {
	columns := []table.Column{
		table.NewColumn(columnKeyTitle, "Title", 13).WithFiltered(true),
		table.NewColumn(columnKeyAuthor, "Author", 13).WithFiltered(true),
		table.NewColumn(columnKeyDescription, "Description", 50),
	}
	return Model{
		table: table.
			New(columns).
			Filtered(true).
			Focused(true).
			WithFooterVisibility(false).
			WithPageSize(10).
			WithRows([]table.Row{
				table.NewRow(table.RowData{
					columnKeyTitle:       "Computer Systems : A Programmer's Perspective",
					columnKeyAuthor:      "Randal E. Bryant、David R. O'Hallaron / Prentice Hall ",
					columnKeyDescription: "This book explains the important and enduring concepts underlying all computer...",
				}),
				table.NewRow(table.RowData{
					columnKeyTitle:       "Effective Java : 3rd Edition",
					columnKeyAuthor:      "Joshua Bloch",
					columnKeyDescription: "The Definitive Guide to Java Platform Best Practices—Updated for Java 9 Java ...",
				}),
				table.NewRow(table.RowData{
					columnKeyTitle:       "Structure and Interpretation of Computer Programs - 2nd Edition (MIT)",
					columnKeyAuthor:      "Harold Abelson、Gerald Jay Sussman",
					columnKeyDescription: "Structure and Interpretation of Computer Programs has had a dramatic impact on...",
				}),
				table.NewRow(table.RowData{
					columnKeyTitle:       "Game Programming Patterns",
					columnKeyAuthor:      "Robert Nystrom / Genever Benning",
					columnKeyDescription: "The biggest challenge facing many game programmers is completing their game. M...",
				}),
			}),
		filterTextInput: textinput.New(),
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

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// global
		if msg.String() == "ctrl+c" {
			cmds = append(cmds, tea.Quit)
			return m, tea.Batch(cmds...)
		}
		// event to filter
		if m.filterTextInput.Focused() {
			if msg.String() == "enter" {
				m.filterTextInput.Blur()
			} else {
				m.filterTextInput, cmd = m.filterTextInput.Update(msg)
			}
			m.table = m.table.WithFilterInput(m.filterTextInput)
			return m, tea.Batch(cmds...)
		}

		// others component
		switch msg.String() {
		case "/":
			m.filterTextInput.Focus()
		case "q":
			cmds = append(cmds, tea.Quit)
			return m, tea.Batch(cmds...)
		default:
			m.table, cmd = m.table.Update(msg)
			cmds = append(cmds, cmd)
		}

	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	body := strings.Builder{}

	body.WriteString("A filtered simple default table\n" +
		"Currently filter by Title and Author, press / + letters to start filtering, and escape to clear filter.\nPress q or ctrl+c to quit\n\n")

	body.WriteString(m.filterTextInput.View() + "\n")
	body.WriteString(m.table.View())

	return body.String()
}

func main() {
	p := tea.NewProgram(NewModel())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
