package table

import (
	"fmt"
	"testing"
)

var benchView string

func benchTable(numColumns, numDataRows int) Model {
	columns := []Column{}

	for i := 0; i < numColumns; i++ {
		iStr := fmt.Sprintf("%d", i)
		columns = append(columns, NewColumn(iStr, iStr, 6))
	}

	rows := []Row{}

	for i := 0; i < numDataRows; i++ {
		rowData := RowData{}

		for columnIndex, column := range columns {
			rowData[column.key] = fmt.Sprintf("%d", columnIndex)
		}

		rows = append(rows, NewRow(rowData))
	}

	return New(columns).WithRows(rows)
}

func BenchmarkPlain3x3TableView(b *testing.B) {
	makeRow := func(id, name string, score int) Row {
		return NewRow(RowData{
			"id":    id,
			"name":  name,
			"score": score,
		})
	}

	model := New([]Column{
		NewColumn("id", "ID", 3),
		NewColumn("name", "Name", 8),
		NewColumn("score", "Score", 6),
	}).WithRows([]Row{
		makeRow("abc", "First", 17),
		makeRow("def", "Second", 1034),
		makeRow("123", "Third", 841),
	})

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		benchView = model.View()
	}
}

func BenchmarkPlainTableViews(b *testing.B) {
	sizes := []struct {
		numColumns int
		numRows    int
	}{
		{
			numColumns: 1,
			numRows:    0,
		},
		{
			numColumns: 10,
			numRows:    0,
		},
		{
			numColumns: 1,
			numRows:    4,
		},
		{
			numColumns: 1,
			numRows:    19,
		},
		{
			numColumns: 9,
			numRows:    19,
		},
		{
			numColumns: 9,
			numRows:    49,
		},
	}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("%dx%d", size.numColumns, size.numRows+1), func(b *testing.B) {
			model := benchTable(size.numColumns, size.numRows)
			b.ResetTimer()

			for n := 0; n < b.N; n++ {
				benchView = model.View()
			}
		})
	}
}
