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

func getVisibleRows(m *Model) []Row {
	start, end := m.VisibleIndices()

	return m.rows[start : end+1]
}

func TestPaginationAccessors(t *testing.T) {
	const (
		numRows  = 100
		pageSize = 20
	)

	model := genPaginationTable(numRows, pageSize)

	assert.Equal(t, numRows, model.TotalRows())
	assert.Equal(t, pageSize, model.PageSize())
}

func TestPaginationNoPageSizeReturnsAll(t *testing.T) {
	const (
		numRows  = 100
		pageSize = 0
	)

	model := genPaginationTable(numRows, pageSize)

	paginatedRows := getVisibleRows(&model)

	assert.Len(t, paginatedRows, numRows)
	assert.Equal(t, 1, model.MaxPages())
}

func TestPaginationEmptyTableReturnsNoRows(t *testing.T) {
	const (
		numRows  = 0
		pageSize = 10
	)

	model := genPaginationTable(numRows, pageSize)

	paginatedRows := getVisibleRows(&model)

	assert.Len(t, paginatedRows, numRows)
}

func TestPaginationDefaultsToAllRows(t *testing.T) {
	const numRows = 100

	model := genPaginationTable(numRows, 0)

	paginatedRows := getVisibleRows(&model)

	assert.Len(t, paginatedRows, numRows)
}

func TestPaginationReturnsPartialFirstPage(t *testing.T) {
	const (
		numRows  = 10
		pageSize = 20
	)

	model := genPaginationTable(numRows, pageSize)

	paginatedRows := getVisibleRows(&model)

	assert.Len(t, paginatedRows, numRows)
}

func TestPaginationReturnsFirstFullPage(t *testing.T) {
	const (
		pageSize = 10
		numRows  = 20
	)

	model := genPaginationTable(numRows, pageSize)

	paginatedRows := getVisibleRows(&model)

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

	paginatedRows := getVisibleRows(&model)

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

	paginatedRows := getVisibleRows(&model)

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

	paginatedRows := getVisibleRows(&model)

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

	paginatedRows := getVisibleRows(&model)

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

	paginatedRows := getVisibleRows(&model)

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

	paginatedRows := getVisibleRows(&model)

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

	paginatedRows := getVisibleRows(&model)

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

// This is long because of various test cases, not because of logic
// nolint: funlen
func TestExpectedPageForRowIndex(t *testing.T) {
	tests := []struct {
		name         string
		totalRows    int
		pageSize     int
		rowIndex     int
		expectedPage int
	}{
		{
			name: "Empty",
		},
		{
			name:         "No pages",
			totalRows:    50,
			pageSize:     0,
			rowIndex:     37,
			expectedPage: 0,
		},
		{
			name:         "One page",
			totalRows:    50,
			pageSize:     50,
			rowIndex:     37,
			expectedPage: 0,
		},
		{
			name:         "First page",
			totalRows:    50,
			pageSize:     30,
			rowIndex:     17,
			expectedPage: 0,
		},
		{
			name:         "Second page",
			totalRows:    50,
			pageSize:     30,
			rowIndex:     37,
			expectedPage: 1,
		},
		{
			name:         "First page first row",
			totalRows:    50,
			pageSize:     30,
			rowIndex:     0,
			expectedPage: 0,
		},
		{
			name:         "First page last row",
			totalRows:    50,
			pageSize:     30,
			rowIndex:     29,
			expectedPage: 0,
		},
		{
			name:         "Second page first row",
			totalRows:    50,
			pageSize:     30,
			rowIndex:     30,
			expectedPage: 1,
		},
		{
			name:         "Second page last row",
			totalRows:    50,
			pageSize:     30,
			rowIndex:     49,
			expectedPage: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			model := genPaginationTable(test.totalRows, test.pageSize)

			page := model.expectedPageForRowIndex(test.rowIndex)

			assert.Equal(t, test.expectedPage, page)
		})
	}
}

func TestClearPagination(t *testing.T) {
	const (
		pageSize = 10
		numRows  = 20
	)

	model := genPaginationTable(numRows, pageSize)

	assert.Equal(t, 1, model.expectedPageForRowIndex(11))

	model = model.WithNoPagination()

	assert.Equal(t, 0, model.expectedPageForRowIndex(11))
}
