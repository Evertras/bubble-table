package table

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This function is only long because of repetitive test definitions, this is fine
// nolint: funlen
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
