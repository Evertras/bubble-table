package table

import (
	"fmt"
	"strings"
	"testing"

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

func TestSimple3x2(t *testing.T) {
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
