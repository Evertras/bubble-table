package table

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This function is only long because of repetitive test definitions, this is fine
//
//nolint:funlen
func TestColumnUpdateWidths(t *testing.T) {
	tests := []struct {
		name           string
		columns        []Column
		totalWidth     int
		expectedWidths []int
	}{
		{
			name: "Static",
			columns: []Column{
				NewColumn("abc", "a", 4),
				NewColumn("sdf", "b", 7),
				NewColumn("xyz", "c", 2),
			},
			totalWidth: 13,
			expectedWidths: []int{
				4, 7, 2,
			},
		},
		{
			name: "Even half",
			columns: []Column{
				NewFlexColumn("abc", "a", 1),
				NewFlexColumn("sdf", "b", 1),
			},
			totalWidth: 11,
			expectedWidths: []int{
				4, 4,
			},
		},
		{
			name: "Odd half increases first",
			columns: []Column{
				NewFlexColumn("abc", "a", 1),
				NewFlexColumn("sdf", "b", 1),
			},
			totalWidth: 12,
			expectedWidths: []int{
				5, 4,
			},
		},
		{
			name: "Even fourths",
			columns: []Column{
				NewFlexColumn("abc", "a", 1),
				NewFlexColumn("sdf", "b", 1),
				NewFlexColumn("xyz", "c", 1),
				NewFlexColumn("123", "d", 1),
			},
			totalWidth: 17,
			expectedWidths: []int{
				3, 3, 3, 3,
			},
		},
		{
			name: "Odd fourths",
			columns: []Column{
				NewFlexColumn("abc", "a", 1),
				NewFlexColumn("sdf", "b", 1),
				NewFlexColumn("xyz", "c", 1),
				NewFlexColumn("123", "d", 1),
			},
			totalWidth: 20,
			expectedWidths: []int{
				4, 4, 4, 3,
			},
		},
		{
			name: "Simple mix static and flex",
			columns: []Column{
				NewColumn("abc", "a", 5),
				NewFlexColumn("flex", "flex", 1),
			},
			totalWidth: 18,
			expectedWidths: []int{
				5, 10,
			},
		},
		{
			name: "Static and flex with high flex factor",
			columns: []Column{
				NewColumn("abc", "a", 5),
				NewFlexColumn("flex", "flex", 1000),
			},
			totalWidth: 18,
			expectedWidths: []int{
				5, 10,
			},
		},
		{
			name: "Static and multiple flexes with high flex factor",
			columns: []Column{
				NewColumn("abc", "a", 5),
				NewFlexColumn("flex", "flex", 1000),
				NewFlexColumn("flex", "flex", 1000),
				NewFlexColumn("flex", "flex", 1000),
			},
			totalWidth: 22,
			expectedWidths: []int{
				5, 4, 4, 4,
			},
		},
		{
			name: "Static and multiple flexes of different sizes",
			columns: []Column{
				NewFlexColumn("flex", "flex", 1),
				NewColumn("abc", "a", 5),
				NewFlexColumn("flex", "flex", 2),
				NewFlexColumn("flex", "flex", 1),
			},
			totalWidth: 22,
			expectedWidths: []int{
				3, 5, 6, 3,
			},
		},
		{
			name: "Width is too small",
			columns: []Column{
				NewColumn("abc", "a", 5),
				NewFlexColumn("flex", "flex", 2),
				NewFlexColumn("flex", "flex", 1),
			},
			totalWidth: 3,
			expectedWidths: []int{
				5, 1, 1,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			updateColumnWidths(test.columns, test.totalWidth)

			for i, col := range test.columns {
				assert.Equal(t, test.expectedWidths[i], col.width, fmt.Sprintf("index %d", i))
			}
		})
	}
}

// This function is long because it has many test cases
//
//nolint:funlen
func TestRecalculateHeight(t *testing.T) {
	columns := []Column{
		NewColumn("ka", "a", 3),
		NewColumn("kb", "b", 4),
		NewColumn("kc", "c", 5),
	}

	rows := []Row{
		NewRow(RowData{"ka": 1, "kb": 23, "kc": "zyx"}),
		NewRow(RowData{"ka": 3, "kb": 34, "kc": "wvu"}),
		NewRow(RowData{"ka": 5, "kb": 45, "kc": "zyx"}),
		NewRow(RowData{"ka": 7, "kb": 56, "kc": "wvu"}),
	}

	tests := []struct {
		name           string
		model          Model
		expectedHeight int
	}{
		{
			name:           "Default header",
			model:          New(columns).WithRows(rows),
			expectedHeight: 3,
		},
		{
			name:           "Empty page with default header",
			model:          New(columns),
			expectedHeight: 3,
		},
		{
			name:           "Filtered with default header",
			model:          New(columns).WithRows(rows).Filtered(true),
			expectedHeight: 5,
		},
		{
			name:           "Static footer one line",
			model:          New(columns).WithRows(rows).WithStaticFooter("single line"),
			expectedHeight: 5,
		},
		{
			name: "Static footer overflow",
			model: New(columns).WithRows(rows).
				WithStaticFooter("single line but it's long"),
			expectedHeight: 6,
		},
		{
			name: "Static footer multi-line",
			model: New(columns).WithRows(rows).
				WithStaticFooter("footer with\nmultiple lines"),
			expectedHeight: 6,
		},
		{
			name:           "Paginated",
			model:          New(columns).WithRows(rows).WithPageSize(2),
			expectedHeight: 5,
		},
		{
			name:           "No pagination",
			model:          New(columns).WithRows(rows).WithPageSize(2).WithNoPagination(),
			expectedHeight: 3,
		},
		{
			name:           "Footer not visible",
			model:          New(columns).WithRows(rows).Filtered(true).WithFooterVisibility(false),
			expectedHeight: 3,
		},
		{
			name:           "Header not visible",
			model:          New(columns).WithRows(rows).WithHeaderVisibility(false),
			expectedHeight: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.model.recalculateHeight()
			assert.Equal(t, test.expectedHeight, test.model.metaHeight)
		})
	}
}
