package table

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func genPaginationTable(count, pageSize int) Model {
	model := New([]Column{
		NewColumn("id", "ID", 3),
	})

	rows := []Row{}

	for i := 1; i <= count; i++ {
		rows = append(rows, NewRow(RowData{
			"id": i,
		}))
	}

	return model.WithRows(rows).WithPageSize(pageSize)
}

func paginationRowID(row Row) int {
	rowID, ok := row.Data["id"].(int)

	if !ok {
		panic("id not int, bad test")
	}

	return rowID
}

func TestPaginationNoPageSizeReturnsAll(t *testing.T) {
	const (
		numRows  = 100
		pageSize = 0
	)

	model := genPaginationTable(numRows, pageSize)

	paginatedRows := model.getVisibleRows()

	assert.Len(t, paginatedRows, numRows)
}

func TestPaginationEmptyTableReturnsNoRows(t *testing.T) {
	const (
		numRows  = 0
		pageSize = 10
	)

	model := genPaginationTable(numRows, pageSize)

	paginatedRows := model.getVisibleRows()

	assert.Len(t, paginatedRows, numRows)
}

func TestPaginationDefaultsToAllRows(t *testing.T) {
	const numRows = 100

	model := genPaginationTable(numRows, 0)

	paginatedRows := model.getVisibleRows()

	assert.Len(t, paginatedRows, numRows)
}

func TestPaginationReturnsPartialFirstPage(t *testing.T) {
	const (
		numRows  = 10
		pageSize = 20
	)

	model := genPaginationTable(numRows, pageSize)

	paginatedRows := model.getVisibleRows()

	assert.Len(t, paginatedRows, numRows)
}

func TestPaginationReturnsFirstFullPage(t *testing.T) {
	const (
		pageSize = 10
		numRows  = 20
	)

	model := genPaginationTable(numRows, pageSize)

	paginatedRows := model.getVisibleRows()

	assert.Len(t, paginatedRows, pageSize)

	for i, row := range paginatedRows {
		assert.Equal(t, i+1, paginationRowID(row))
	}
}

func TestPaginationReturnsSecondFullPageAfterMoving(t *testing.T) {
	const (
		pageSize = 10
		numRows  = 30
	)

	model := genPaginationTable(numRows, pageSize)

	model.pageDown()

	paginatedRows := model.getVisibleRows()

	assert.Len(t, paginatedRows, pageSize)

	for i, row := range paginatedRows {
		assert.Equal(t, i+11, paginationRowID(row))
	}
}

func TestPaginationReturnsPartialFinalPage(t *testing.T) {
	const (
		pageSize = 10
		numRows  = 15
	)

	model := genPaginationTable(numRows, pageSize)

	model.pageDown()

	paginatedRows := model.getVisibleRows()

	assert.Len(t, paginatedRows, numRows-pageSize)

	for i, row := range paginatedRows {
		assert.Equal(t, i+11, paginationRowID(row))
	}
}

func TestPaginationWrapsUpPartial(t *testing.T) {
	const (
		pageSize = 10
		numRows  = 15
	)

	model := genPaginationTable(numRows, pageSize)

	model.pageUp()

	paginatedRows := model.getVisibleRows()

	assert.Len(t, paginatedRows, numRows-pageSize)

	for i, row := range paginatedRows {
		assert.Equal(t, i+11, paginationRowID(row))
	}
}

func TestPaginationWrapsUpFull(t *testing.T) {
	const (
		pageSize = 10
		numRows  = 20
	)

	model := genPaginationTable(numRows, pageSize)

	model.pageUp()

	paginatedRows := model.getVisibleRows()

	assert.Len(t, paginatedRows, numRows-pageSize)

	for i, row := range paginatedRows {
		assert.Equal(t, i+11, paginationRowID(row))
	}
}

func TestPaginationWrapsUpSelf(t *testing.T) {
	const (
		pageSize = 10
		numRows  = 10
	)

	model := genPaginationTable(numRows, pageSize)

	model.pageUp()

	paginatedRows := model.getVisibleRows()

	assert.Len(t, paginatedRows, numRows)

	for i, row := range paginatedRows {
		assert.Equal(t, i+1, paginationRowID(row))
	}
}

func TestPaginationWrapsDown(t *testing.T) {
	const (
		pageSize = 10
		numRows  = 15
	)

	model := genPaginationTable(numRows, pageSize)

	model.pageDown()
	model.pageDown()

	paginatedRows := model.getVisibleRows()

	assert.Len(t, paginatedRows, pageSize)

	for i, row := range paginatedRows {
		assert.Equal(t, i+1, paginationRowID(row))
	}
}

func TestPaginationWrapsDownSelf(t *testing.T) {
	const (
		pageSize = 10
		numRows  = 10
	)

	model := genPaginationTable(numRows, pageSize)

	model.pageDown()
	model.pageDown()

	paginatedRows := model.getVisibleRows()

	assert.Len(t, paginatedRows, pageSize)

	for i, row := range paginatedRows {
		assert.Equal(t, i+1, paginationRowID(row))
	}
}

func TestPaginationHighlightFirstOnPageDown(t *testing.T) {
	const (
		pageSize = 10
		numRows  = 20
	)

	model := genPaginationTable(numRows, pageSize)

	assert.Equal(t, 1, paginationRowID(model.HighlightedRow()), "Initial test setup wrong, test code may be bad")

	model.pageDown()

	assert.Equal(t, 11, paginationRowID(model.HighlightedRow()), "Did not highlight expected row")
}
