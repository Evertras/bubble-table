package table

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
)

func TestBasicTableShowsAllHeaders(t *testing.T) {
	const (
		firstKey   = "first-key"
		firstTitle = "First Title"
		firstWidth = 13

		secondKey   = "second-key"
		secondTitle = "Second Title"
		secondWidth = 20
	)

	columns := []Column{
		NewColumn(firstKey, firstTitle, firstWidth),
		NewColumn(secondKey, secondTitle, secondWidth),
	}

	model := New(columns)

	rendered := model.View()

	assert.Contains(t, rendered, firstTitle)
	assert.Contains(t, rendered, secondTitle)

	assert.False(t, strings.HasSuffix(rendered, "\n"), "Should not end in newline")
}

func TestBasicTableTruncatesLongHeaders(t *testing.T) {
	const (
		firstKey   = "first-key"
		firstTitle = "First Title"
		firstWidth = 3

		secondKey   = "second-key"
		secondTitle = "Second Title"
		secondWidth = 3
	)

	columns := []Column{
		NewColumn(firstKey, firstTitle, firstWidth),
		NewColumn(secondKey, secondTitle, secondWidth),
	}

	model := New(columns)

	rendered := model.View()

	assert.Contains(t, rendered, "Fi…")
	assert.Contains(t, rendered, "Se…")

	assert.False(t, strings.HasSuffix(rendered, "\n"), "Should not end in newline")
}

func TestNilColumnsSafelyReturnsEmptyView(t *testing.T) {
	model := New(nil)

	assert.Equal(t, "", model.View())
}

func TestSingleCellView(t *testing.T) {
	model := New([]Column{
		NewColumn("id", "ID", 4),
	})

	const expectedTable = `┏━━━━┓
┃  ID┃
┗━━━━┛`

	rendered := model.View()

	assert.Equal(t, expectedTable, rendered)
}

func TestSingleColumnView(t *testing.T) {
	model := New([]Column{
		NewColumn("id", "ID", 4),
	}).WithRows([]Row{
		NewRow(RowData{"id": "1"}),
		NewRow(RowData{"id": "2"}),
	})

	const expectedTable = `┏━━━━┓
┃  ID┃
┣━━━━┫
┃   1┃
┃   2┃
┗━━━━┛`

	rendered := model.View()

	assert.Equal(t, expectedTable, rendered)
}

func TestSingleColumnViewSorted(t *testing.T) {
	model := New([]Column{
		NewColumn("id", "ID", 4),
	}).WithRows([]Row{
		NewRow(RowData{"id": "1"}),
		NewRow(RowData{"id": "2"}),
	}).SortByDesc("id")

	const expectedTable = `┏━━━━┓
┃  ID┃
┣━━━━┫
┃   2┃
┃   1┃
┗━━━━┛`

	rendered := model.View()

	assert.Equal(t, expectedTable, rendered)
}

func TestSingleRowView(t *testing.T) {
	model := New([]Column{
		NewColumn("1", "1", 4),
		NewColumn("2", "2", 4),
		NewColumn("3", "3", 4),
	})

	const expectedTable = `┏━━━━┳━━━━┳━━━━┓
┃   1┃   2┃   3┃
┗━━━━┻━━━━┻━━━━┛`

	rendered := model.View()

	assert.Equal(t, expectedTable, rendered)
}

func TestSimple3x3(t *testing.T) {
	model := New([]Column{
		NewColumn("1", "1", 4),
		NewColumn("2", "2", 4),
		NewColumn("3", "3", 4),
	})

	rows := []Row{}

	for rowIndex := 1; rowIndex <= 3; rowIndex++ {
		rowData := RowData{}

		for columnIndex := 1; columnIndex <= 3; columnIndex++ {
			id := fmt.Sprintf("%d", columnIndex)

			rowData[id] = fmt.Sprintf("%d,%d", columnIndex, rowIndex)
		}

		rows = append(rows, NewRow(rowData))
	}

	model = model.WithRows(rows)

	const expectedTable = `┏━━━━┳━━━━┳━━━━┓
┃   1┃   2┃   3┃
┣━━━━╋━━━━╋━━━━┫
┃ 1,1┃ 2,1┃ 3,1┃
┃ 1,2┃ 2,2┃ 3,2┃
┃ 1,3┃ 2,3┃ 3,3┃
┗━━━━┻━━━━┻━━━━┛`

	rendered := model.View()

	assert.Equal(t, expectedTable, rendered)
}

func TestSingleHeaderWithFooter(t *testing.T) {
	model := New([]Column{
		NewColumn("id", "ID", 4),
	}).WithStaticFooter("Foot")

	const expectedTable = `┏━━━━┓
┃  ID┃
┣━━━━┫
┃Foot┃
┗━━━━┛`
	rendered := model.View()

	assert.Equal(t, expectedTable, rendered)
}

func TestSingleRowWithFooterView(t *testing.T) {
	model := New([]Column{
		NewColumn("1", "1", 4),
		NewColumn("2", "2", 4),
		NewColumn("3", "3", 4),
	}).WithStaticFooter("Footer")

	const expectedTable = `┏━━━━┳━━━━┳━━━━┓
┃   1┃   2┃   3┃
┣━━━━┻━━━━┻━━━━┫
┃        Footer┃
┗━━━━━━━━━━━━━━┛`

	rendered := model.View()

	assert.Equal(t, expectedTable, rendered)
}

func TestSingleRowWithFooterViewAndBaseStyle(t *testing.T) {
	model := New([]Column{
		NewColumn("1", "1", 4),
		NewColumn("2", "2", 4),
		NewColumn("3", "3", 4),
	}).WithStaticFooter("Footer").WithBaseStyle(lipgloss.NewStyle().Align(lipgloss.Left))

	const expectedTable = `┏━━━━┳━━━━┳━━━━┓
┃1   ┃2   ┃3   ┃
┣━━━━┻━━━━┻━━━━┫
┃Footer        ┃
┗━━━━━━━━━━━━━━┛`

	rendered := model.View()

	assert.Equal(t, expectedTable, rendered)
}

func TestSingleColumnWithFooterView(t *testing.T) {
	model := New([]Column{
		NewColumn("id", "ID", 4),
	}).WithRows([]Row{
		NewRow(RowData{"id": "1"}),
		NewRow(RowData{"id": "2"}),
	}).WithStaticFooter("Foot")

	const expectedTable = `┏━━━━┓
┃  ID┃
┣━━━━┫
┃   1┃
┃   2┃
┣━━━━┫
┃Foot┃
┗━━━━┛`

	rendered := model.View()

	assert.Equal(t, expectedTable, rendered)
}

func TestSimple3x3WithFooterView(t *testing.T) {
	model := New([]Column{
		NewColumn("1", "1", 4),
		NewColumn("2", "2", 4),
		NewColumn("3", "3", 4),
	})

	rows := []Row{}

	for rowIndex := 1; rowIndex <= 3; rowIndex++ {
		rowData := RowData{}

		for columnIndex := 1; columnIndex <= 3; columnIndex++ {
			id := fmt.Sprintf("%d", columnIndex)

			rowData[id] = fmt.Sprintf("%d,%d", columnIndex, rowIndex)
		}

		rows = append(rows, NewRow(rowData))
	}

	model = model.WithRows(rows).WithStaticFooter("Footer")

	const expectedTable = `┏━━━━┳━━━━┳━━━━┓
┃   1┃   2┃   3┃
┣━━━━╋━━━━╋━━━━┫
┃ 1,1┃ 2,1┃ 3,1┃
┃ 1,2┃ 2,2┃ 3,2┃
┃ 1,3┃ 2,3┃ 3,3┃
┣━━━━┻━━━━┻━━━━┫
┃        Footer┃
┗━━━━━━━━━━━━━━┛`

	rendered := model.View()

	assert.Equal(t, expectedTable, rendered)
}

func TestPaged3x3WithNoSpecifiedFooter(t *testing.T) {
	model := New([]Column{
		NewColumn("1", "1", 4),
		NewColumn("2", "2", 4),
		NewColumn("3", "3", 4),
	})

	rows := []Row{}

	for rowIndex := 1; rowIndex <= 3; rowIndex++ {
		rowData := RowData{}

		for columnIndex := 1; columnIndex <= 3; columnIndex++ {
			id := fmt.Sprintf("%d", columnIndex)

			rowData[id] = fmt.Sprintf("%d,%d", columnIndex, rowIndex)
		}

		rows = append(rows, NewRow(rowData))
	}

	model = model.WithRows(rows).WithPageSize(2)

	const expectedTable = `┏━━━━┳━━━━┳━━━━┓
┃   1┃   2┃   3┃
┣━━━━╋━━━━╋━━━━┫
┃ 1,1┃ 2,1┃ 3,1┃
┃ 1,2┃ 2,2┃ 3,2┃
┣━━━━┻━━━━┻━━━━┫
┃           1/2┃
┗━━━━━━━━━━━━━━┛`

	rendered := model.View()

	assert.Equal(t, expectedTable, rendered)
}

func TestPaged3x3WithStaticFooter(t *testing.T) {
	model := New([]Column{
		NewColumn("1", "1", 4),
		NewColumn("2", "2", 4),
		NewColumn("3", "3", 4),
	})

	rows := []Row{}

	for rowIndex := 1; rowIndex <= 3; rowIndex++ {
		rowData := RowData{}

		for columnIndex := 1; columnIndex <= 3; columnIndex++ {
			id := fmt.Sprintf("%d", columnIndex)

			rowData[id] = fmt.Sprintf("%d,%d", columnIndex, rowIndex)
		}

		rows = append(rows, NewRow(rowData))
	}

	model = model.WithRows(rows).WithPageSize(2).WithStaticFooter("Override")

	const expectedTable = `┏━━━━┳━━━━┳━━━━┓
┃   1┃   2┃   3┃
┣━━━━╋━━━━╋━━━━┫
┃ 1,1┃ 2,1┃ 3,1┃
┃ 1,2┃ 2,2┃ 3,2┃
┣━━━━┻━━━━┻━━━━┫
┃      Override┃
┗━━━━━━━━━━━━━━┛`

	rendered := model.View()

	assert.Equal(t, expectedTable, rendered)
}

func TestSimple3x3StyleOverridesAsBaseColumnRowCell(t *testing.T) {
	model := New([]Column{
		NewColumn("1", "1", 6),
		NewColumn("2", "2", 6).WithStyle(lipgloss.NewStyle().Align(lipgloss.Left)),
		NewColumn("3", "3", 6),
	}).WithBaseStyle(lipgloss.NewStyle().Align(lipgloss.Center))

	rows := []Row{}

	for rowIndex := 1; rowIndex <= 3; rowIndex++ {
		rowData := RowData{}

		for columnIndex := 1; columnIndex <= 3; columnIndex++ {
			id := fmt.Sprintf("%d", columnIndex)

			rowData[id] = fmt.Sprintf("%d,%d", columnIndex, rowIndex)
		}

		rows = append(rows, NewRow(rowData))
	}

	// Test overrides with alignment because it's easy to check output string
	rows[0] = rows[0].WithStyle(lipgloss.NewStyle().Align(lipgloss.Left))
	rows[0].Data["2"] = NewStyledCell("R", lipgloss.NewStyle().Align(lipgloss.Right))

	rows[2] = rows[2].WithStyle(lipgloss.NewStyle().Align(lipgloss.Right))

	model = model.WithRows(rows)

	const expectedTable = `┏━━━━━━┳━━━━━━┳━━━━━━┓
┃  1   ┃2     ┃  3   ┃
┣━━━━━━╋━━━━━━╋━━━━━━┫
┃1,1   ┃     R┃3,1   ┃
┃ 1,2  ┃2,2   ┃ 3,2  ┃
┃   1,3┃   2,3┃   3,3┃
┗━━━━━━┻━━━━━━┻━━━━━━┛`

	rendered := model.View()

	assert.Equal(t, expectedTable, rendered)
}

// This is a long test due to typing and multiple big table strings, that's okay
// nolint: funlen
func Test3x3WithFilterFooter(t *testing.T) {
	model := New([]Column{
		NewColumn("1", "1", 4).WithFiltered(true),
		NewColumn("2", "2", 4),
		NewColumn("3", "3", 4),
	})

	rows := []Row{}

	for rowIndex := 1; rowIndex <= 3; rowIndex++ {
		rowData := RowData{}

		for columnIndex := 1; columnIndex <= 3; columnIndex++ {
			id := fmt.Sprintf("%d", columnIndex)

			rowData[id] = fmt.Sprintf("%d,%d", columnIndex, rowIndex)
		}

		rows = append(rows, NewRow(rowData))
	}

	model = model.WithRows(rows).Filtered(true).Focused(true)

	const expectedTable = `┏━━━━┳━━━━┳━━━━┓
┃   1┃   2┃   3┃
┣━━━━╋━━━━╋━━━━┫
┃ 1,1┃ 2,1┃ 3,1┃
┃ 1,2┃ 2,2┃ 3,2┃
┃ 1,3┃ 2,3┃ 3,3┃
┣━━━━┻━━━━┻━━━━┫
┃              ┃
┗━━━━━━━━━━━━━━┛`

	assert.Equal(t, expectedTable, model.View())

	hitKey := func(key rune) {
		model, _ = model.Update(
			tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune{key},
			})
	}

	hitKey('/')
	hitKey('3')

	// The byte code near the bottom is a block cursor from the text box
	const expectedFilteredTypingTable = `┏━━━━┳━━━━┳━━━━┓
┃   1┃   2┃   3┃
┣━━━━╋━━━━╋━━━━┫
┃ 1,3┃ 2,3┃ 3,3┃
┣━━━━┻━━━━┻━━━━┫
┃           /3` + "\x1b[7m \x1b[0m" + `┃
┗━━━━━━━━━━━━━━┛`

	assert.Equal(t, expectedFilteredTypingTable, model.View())

	const expectedFilteredDoneTable = `┏━━━━┳━━━━┳━━━━┓
┃   1┃   2┃   3┃
┣━━━━╋━━━━╋━━━━┫
┃ 1,3┃ 2,3┃ 3,3┃
┣━━━━┻━━━━┻━━━━┫
┃            /3┃
┗━━━━━━━━━━━━━━┛`

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, expectedFilteredDoneTable, model.View())
}

func TestSingleCellFlexView(t *testing.T) {
	model := New([]Column{
		NewFlexColumn("id", "ID", 1),
	}).WithTargetWidth(6)

	const expectedTable = `┏━━━━┓
┃  ID┃
┗━━━━┛`

	rendered := model.View()

	assert.Equal(t, expectedTable, rendered)
}

func TestSimpleFlex3x3(t *testing.T) {
	model := New([]Column{
		NewFlexColumn("1", "1", 1),
		NewFlexColumn("2", "2", 1),
		NewFlexColumn("3", "3", 2),
	}).WithTargetWidth(20)

	rows := []Row{}

	for rowIndex := 1; rowIndex <= 3; rowIndex++ {
		rowData := RowData{}

		for columnIndex := 1; columnIndex <= 3; columnIndex++ {
			id := fmt.Sprintf("%d", columnIndex)

			rowData[id] = fmt.Sprintf("%d,%d", columnIndex, rowIndex)
		}

		rows = append(rows, NewRow(rowData))
	}

	model = model.WithRows(rows)

	const expectedTable = `┏━━━━┳━━━━┳━━━━━━━━┓
┃   1┃   2┃       3┃
┣━━━━╋━━━━╋━━━━━━━━┫
┃ 1,1┃ 2,1┃     3,1┃
┃ 1,2┃ 2,2┃     3,2┃
┃ 1,3┃ 2,3┃     3,3┃
┗━━━━┻━━━━┻━━━━━━━━┛`

	rendered := model.View()

	assert.Equal(t, expectedTable, rendered)
}
