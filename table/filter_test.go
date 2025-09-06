package table

import (
	"fmt"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
)

func TestIsRowMatched(t *testing.T) {
	columns := []Column{
		NewColumn("title", "title", 10).WithFiltered(true),
		NewColumn("description", "description", 10)}

	assert.True(t, filterFuncContains(FilterFuncInput{
		Columns: columns,
		Row: NewRow(RowData{
			"title":       "AAA",
			"description": "",
		}),
		Filter: "",
	}))

	type testCase struct {
		name        string
		filter      string
		title       any
		description any
		shouldMatch bool
	}

	timeFrom2020 := time.Date(2020, time.July, 1, 1, 1, 1, 1, time.UTC)

	cases := []testCase{
		{"empty filter matches all", "", "AAA", "", true},
		{"exact match", "AAA", "AAA", "", true},
		{"partial match start", "A", "AAA", "", true},
		{"partial match middle", "AA", "AAA", "", true},
		{"too long", "AAAA", "AAA", "", false},
		{"lowercase", "aaa", "AAA", "", true},
		{"mixed case", "AaA", "AAA", "", true},
		{"wrong input", "B", "AAA", "", false},
		{"ignore description", "BBB", "AAA", "BBB", false},
		{"time filterable success", "2020", timeFrom2020, "", true},
		{"time filterable wrong input", "2021", timeFrom2020, "", false},
		{"styled cell", "AAA", NewStyledCell("AAA", lipgloss.NewStyle()), "", true},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.shouldMatch, filterFuncContains(FilterFuncInput{
				Columns: columns,
				Row: NewRow(RowData{
					"title":       testCase.title,
					"description": testCase.description,
				}),
				Filter: testCase.filter,
			}))
		})
	}

	// Styled check
}

func TestIsRowMatchedForNonStringer(t *testing.T) {
	columns := []Column{
		NewColumn("val", "val", 10).WithFiltered(true),
	}

	type testCase struct {
		name        string
		filter      string
		val         any
		shouldMatch bool
	}

	cases := []testCase{
		{"exact match", "12", 12, true},
		{"partial match", "1", 12, true},
		{"partial match end", "2", 12, true},
		{"wrong input", "3", 12, false},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.shouldMatch, filterFuncContains(FilterFuncInput{
				Columns: columns,
				Row: NewRow(RowData{
					"val": testCase.val,
				}),
				Filter: testCase.filter,
			}))
		})
	}
}

func TestGetFilteredRowsNoColumnFiltered(t *testing.T) {
	columns := []Column{NewColumn("title", "title", 10)}
	rows := []Row{
		NewRow(RowData{
			"title":       "AAA",
			"description": "",
		}),
		NewRow(RowData{
			"title":       "BBB",
			"description": "",
		}),
		NewRow(RowData{
			"title":       "CCC",
			"description": "",
		}),
	}

	model := New(columns).WithRows(rows).Filtered(true)
	model.filterTextInput.SetValue("AAA")

	filteredRows := model.getFilteredRows(rows)

	assert.Len(t, filteredRows, len(rows))
}

func TestGetFilteredRowsUnfiltered(t *testing.T) {
	columns := []Column{NewColumn("title", "title", 10)}
	rows := []Row{
		NewRow(RowData{
			"title": "AAA",
		}),
		NewRow(RowData{
			"title": "BBB",
		}),
	}

	model := New(columns).WithRows(rows)

	filteredRows := model.getFilteredRows(rows)

	assert.Len(t, filteredRows, len(rows))
}

func TestGetFilteredRowsFiltered(t *testing.T) {
	columns := []Column{NewColumn("title", "title", 10).WithFiltered(true)}
	rows := []Row{
		NewRow(RowData{
			"title":       "AAA",
			"description": "",
		}),
		NewRow(RowData{
			"title":       "BBB",
			"description": "",
		}),
		// Empty
		NewRow(RowData{}),
	}
	model := New(columns).WithRows(rows).Filtered(true)
	model.filterTextInput.SetValue("AaA")

	filteredRows := model.getFilteredRows(rows)

	assert.Len(t, filteredRows, 1)
}

func TestGetFilteredRowsRefocusAfterFilter(t *testing.T) {
	columns := []Column{NewColumn("title", "title", 10).WithFiltered(true)}
	rows := []Row{
		NewRow(RowData{
			"title":       "a",
			"description": "",
		}),
		NewRow(RowData{
			"title":       "b",
			"description": "",
		}),
		NewRow(RowData{
			"title":       "c",
			"description": "",
		}),
		NewRow(RowData{
			"title":       "d1",
			"description": "",
		}),
		NewRow(RowData{
			"title":       "d2",
			"description": "",
		}),
	}
	model := New(columns).WithRows(rows).Filtered(true).WithPageSize(1)
	model = model.PageDown()
	assert.Len(t, model.GetVisibleRows(), 5)
	assert.Equal(t, 1, model.PageSize())
	assert.Equal(t, 2, model.CurrentPage())
	assert.Equal(t, 5, model.MaxPages())
	assert.Equal(t, 5, model.TotalRows())

	model.filterTextInput.SetValue("c")
	model, _ = model.updateFilterTextInput(tea.KeyMsg{})
	assert.Len(t, model.GetVisibleRows(), 1)
	assert.Equal(t, 1, model.PageSize())
	assert.Equal(t, 1, model.CurrentPage())
	assert.Equal(t, 1, model.MaxPages())
	assert.Equal(t, 1, model.TotalRows())

	model.filterTextInput.SetValue("not-exist")
	model, _ = model.updateFilterTextInput(tea.KeyMsg{})
	assert.Len(t, model.GetVisibleRows(), 0)
	assert.Equal(t, 1, model.PageSize())
	assert.Equal(t, 1, model.CurrentPage())
	assert.Equal(t, 1, model.MaxPages())
	assert.Equal(t, 0, model.TotalRows())
}

func TestFilterWithExternalTextInput(t *testing.T) {
	columns := []Column{NewColumn("title", "title", 10).WithFiltered(true)}
	rows := []Row{
		NewRow(RowData{
			"title":       "AAA",
			"description": "",
		}),
		NewRow(RowData{
			"title":       "BBB",
			"description": "",
		}),
		// Empty
		NewRow(RowData{}),
	}

	// Page size 1 to test scrolling back if input changes
	model := New(columns).WithRows(rows).Filtered(true).WithPageSize(1)
	model.pageDown()
	assert.Equal(t, 2, model.CurrentPage(), "Should start on second page for test")
	input := textinput.New()
	input.SetValue("AaA")
	model = model.WithFilterInput(input)
	assert.Equal(t, 1, model.CurrentPage(), "Did not go back to first page")

	filteredRows := model.getFilteredRows(rows)

	assert.Len(t, filteredRows, 1)
}

func TestFilterWithSetValue(t *testing.T) {
	columns := []Column{NewColumn("title", "title", 10).WithFiltered(true)}
	rows := []Row{
		NewRow(RowData{
			"title":       "AAA",
			"description": "",
		}),
		NewRow(RowData{
			"title":       "BBB",
			"description": "",
		}),
		// Empty
		NewRow(RowData{}),
	}

	// Page size 1 to make sure we scroll back correctly
	model := New(columns).WithRows(rows).Filtered(true).WithPageSize(1)
	model.pageDown()
	assert.Equal(t, 2, model.CurrentPage(), "Should start on second page for test")
	model = model.WithFilterInputValue("AaA")

	assert.Equal(t, 1, model.CurrentPage(), "Did not go back to first page")

	filteredRows := model.getFilteredRows(rows)
	assert.Len(t, filteredRows, 1)

	// Make sure it holds true after an update
	model, _ = model.Update(tea.KeyRight)
	filteredRows = model.getFilteredRows(rows)
	assert.Len(t, filteredRows, 1)

	// Remove filter
	model = model.WithFilterInputValue("")
	filteredRows = model.getFilteredRows(rows)
	assert.Len(t, filteredRows, 3)
}

func TestFilterFunc(t *testing.T) {
	const (
		colTitle = "title"
		colDesc  = "description"
	)

	columns := []Column{NewColumn("title", "title", 10).WithFiltered(true)}
	rows := []Row{
		NewRow(RowData{
			colTitle: "AAA",
			colDesc:  "",
		}),
		NewRow(RowData{
			colTitle: "BBB",
			colDesc:  "",
		}),
		// Empty
		NewRow(RowData{}),
	}

	filterFunc := func(input FilterFuncInput) bool {
		// Completely arbitrary check for testing purposes
		title := fmt.Sprintf("%v", input.Row.Data["title"])

		return title == "AAA" && input.Filter == "x" && input.GlobalMetadata["testValue"] == 3
	}

	// First check that the table won't match with different case
	model := New(columns).WithRows(rows).Filtered(true).WithGlobalMetadata(map[string]any{
		"testValue": 3,
	})
	model = model.WithFilterInputValue("x")

	filteredRows := model.getFilteredRows(rows)
	assert.Len(t, filteredRows, 0)

	// The filter func should then match the one row
	model = model.WithFilterFunc(filterFunc)
	filteredRows = model.getFilteredRows(rows)
	assert.Len(t, filteredRows, 1)

	// Remove filter
	model = model.WithFilterInputValue("")
	filteredRows = model.getFilteredRows(rows)
	assert.Len(t, filteredRows, 3)
}

func BenchmarkFilteredScrolling(b *testing.B) {
	// Scrolling through a filtered table with many rows should be quick
	// https://github.com/Evertras/bubble-table/issues/135
	const rowCount = 40000
	columns := []Column{NewColumn("title", "title", 10).WithFiltered(true)}
	rows := make([]Row, rowCount)

	for i := 0; i < rowCount; i++ {
		rows[i] = NewRow(RowData{
			"title": fmt.Sprintf("%d", i),
		})
	}

	model := New(columns).WithRows(rows).Filtered(true)
	model = model.WithFilterInputValue("1")

	hitKey := func(key rune) {
		model, _ = model.Update(
			tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune{key},
			})
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		hitKey('j')
	}
}

func BenchmarkFilteredScrollingPaged(b *testing.B) {
	// Scrolling through a filtered table with many rows should be quick
	// https://github.com/Evertras/bubble-table/issues/135
	const rowCount = 40000
	columns := []Column{NewColumn("title", "title", 10).WithFiltered(true)}
	rows := make([]Row, rowCount)

	for i := 0; i < rowCount; i++ {
		rows[i] = NewRow(RowData{
			"title": fmt.Sprintf("%d", i),
		})
	}

	model := New(columns).WithRows(rows).Filtered(true).WithPageSize(50)
	model = model.WithFilterInputValue("1")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		model, _ = model.Update(
			tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune{'j'},
			})
	}
}

func BenchmarkFilteredRenders(b *testing.B) {
	// Rendering a filtered table should be fast
	// https://github.com/Evertras/bubble-table/issues/135
	const rowCount = 40000
	columns := []Column{NewColumn("title", "title", 10).WithFiltered(true)}
	rows := make([]Row, rowCount)

	for i := 0; i < rowCount; i++ {
		rows[i] = NewRow(RowData{
			"title": fmt.Sprintf("%d", i),
		})
	}

	model := New(columns).WithRows(rows).Filtered(true).WithPageSize(50)
	model = model.WithFilterInputValue("1")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Don't care about result, just rendering
		_ = model.View()
	}
}

func TestFuzzyFilter_EmptyFilterMatchesAll(t *testing.T) {
	cols := []Column{
		NewColumn("name", "Name", 10).WithFiltered(true),
	}
	rows := []Row{
		NewRow(RowData{"name": "Acme Steel"}),
		NewRow(RowData{"name": "Globex"}),
	}

	for index, row := range rows {
		if !filterFuncFuzzy(FilterFuncInput{
			Columns: cols,
			Row:     row,
			Filter:  "",
		}) {
			t.Fatalf("row %d should match empty filter", index)
		}
	}
}

func TestFuzzyFilter_SubsequenceAcrossColumns(t *testing.T) {
	cols := []Column{
		NewColumn("name", "Name", 10).WithFiltered(true),
		NewColumn("city", "City", 10).WithFiltered(true),
	}
	row := NewRow(RowData{
		"name": "Acme",
		"city": "Stuttgart",
	})

	type testCase struct {
		name        string
		filter      string
		shouldMatch bool
	}

	testCases := []testCase{
		{"subsequence match", "agt", true},
		{"case-insensitive match", "ACM", true},
		{"not a subsequence", "zzt", false},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.shouldMatch, filterFuncFuzzy(FilterFuncInput{
			Columns: cols,
			Row:     row,
			Filter:  tc.filter,
		}))
	}
}

func TestFuzzyFilter_ColumnNotInRow(t *testing.T) {
	cols := []Column{
		NewColumn("column_name_doesnt_match", "Name", 10).WithFiltered(true),
	}
	row := NewRow(RowData{
		"name": "Acme Steel",
	})

	assert.False(t, filterFuncFuzzy(FilterFuncInput{
		Columns: cols,
		Row:     row,
		Filter:  "steel",
	}), "Shouldn't match")
}

func TestFuzzyFilter_RowHasEmptyHaystack(t *testing.T) {
	cols := []Column{
		NewColumn("name", "Name", 10).WithFiltered(true),
	}
	row := NewRow(RowData{"name": ""})

	// literally any value other than an empty string
	// should not match
	assert.False(t, filterFuncFuzzy(FilterFuncInput{
		Columns: cols,
		Row:     row,
		Filter:  "a",
	}), "Shouldn't match")
}

func TestFuzzyFilter_MultiToken_AND(t *testing.T) {
	cols := []Column{
		NewColumn("name", "Name", 10).WithFiltered(true),
		NewColumn("dept", "Dept", 10).WithFiltered(true),
	}
	row := NewRow(RowData{
		"name": "Wayne Enterprises",
		"dept": "R&D",
	})

	// Both tokens must match as subsequences somewhere in the concatenated haystack
	assert.True(t, filterFuncFuzzy(FilterFuncInput{
		Columns: cols,
		Row:     row,
		Filter:  "wy ent",
	}), "Should match wy ent") // "wy" in Wayne, "ent" in Enterprises
	assert.False(t, filterFuncFuzzy(FilterFuncInput{
		Columns: cols,
		Row:     row,
		Filter:  "wy zzz",
	}), "Shouldn't match wy zzz")
}

func TestFuzzyFilter_IgnoresNonFilterableColumns(t *testing.T) {
	cols := []Column{
		NewColumn("name", "Name", 10).WithFiltered(true),
		NewColumn("secret", "Secret", 10).WithFiltered(false), // should be ignored
	}
	row := NewRow(RowData{
		"name":   "Acme",
		"secret": "topsecretpattern",
	})

	assert.False(t, filterFuncFuzzy(FilterFuncInput{
		Columns: cols,
		Row:     row,
		Filter:  "topsecret",
	}), "Shouldn't match on non-filterable")
}

func TestFuzzyFilter_UnwrapsStyledCell(t *testing.T) {
	cols := []Column{
		NewColumn("name", "Name", 10).WithFiltered(true),
	}
	row := NewRow(RowData{
		"name": NewStyledCell("Nakatomi Plaza", lipgloss.NewStyle()),
	})

	assert.True(t, filterFuncFuzzy(FilterFuncInput{
		Columns: cols,
		Row:     row,
		Filter:  "nak plz",
	}), "Expected fuzzy subsequence to match within StyledCell data")
}

func TestFuzzyFilter_NonStringValuesFormatted(t *testing.T) {
	cols := []Column{
		NewColumn("id", "ID", 6).WithFiltered(true),
	}
	row := NewRow(RowData{
		"id": 12345, // should be formatted via fmt.Sprintf("%v", v)
	})

	assert.True(t, filterFuncFuzzy(FilterFuncInput{
		Columns: cols,
		Row:     row,
		Filter:  "245", // subsequence of "12345"
	}), "expected matcher to format non-strings and match subsequence")
}

func TestFuzzySubSequenceMatch_EmptyString(t *testing.T) {
	assert.True(t, fuzzySubsequenceMatch("anything", ""), "empty needle should match anything")
	assert.False(t, fuzzySubsequenceMatch("", "a"), "non-empty needle should not match empty haystack")
	assert.True(t, fuzzySubsequenceMatch("", ""), "empty needle should match empty haystack")
}
