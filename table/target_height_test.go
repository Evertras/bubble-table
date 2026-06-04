package table

import (
	"fmt"
	"strings"
	"testing"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// genTargetHeightTable builds a table with a fixed-width content column and
// WithTargetHeight set.  Rows with long content will wrap when multiline is on.
func genTargetHeightTable(targetHeight int, multiline bool, rows []Row) Model {
	return New([]Column{
		NewColumn("id", "ID", 3),
		NewColumn("content", "Content", 20),
	}).
		WithRows(rows).
		WithMultiline(multiline).
		WithTargetHeight(targetHeight)
}

func shortRow(id int) Row {
	return NewRow(RowData{"id": id, "content": "short"})
}

func longRow(id int) Row {
	// This content exceeds the 20-char column width and will wrap to 2 lines.
	return NewRow(RowData{"id": id, "content": fmt.Sprintf("row%d long content that wraps", id)})
}

// TestTargetHeightUniformShortRows checks that uniform single-line rows produce
// the expected number of pages and visible row counts.
// metaHeight=5 (header=3 + footer=2), so availableLines = 10-5-1 = 4 rows per page.
func TestTargetHeightUniformShortRows(t *testing.T) {
	rows := make([]Row, 12)
	for i := range rows {
		rows[i] = shortRow(i + 1)
	}

	model := genTargetHeightTable(10, false, rows)

	assert.Equal(t, 3, model.MaxPages(), "12 rows at 4 per page = 3 pages")

	start, end := model.VisibleIndices()
	assert.Equal(t, 0, start)
	assert.Equal(t, 3, end, "page 1 should show rows 0-3")
}

// TestTargetHeightPageNavigationShortRows checks cursor and page after pageDown.
func TestTargetHeightPageNavigationShortRows(t *testing.T) {
	rows := make([]Row, 12)
	for i := range rows {
		rows[i] = shortRow(i + 1)
	}

	model := genTargetHeightTable(10, false, rows)
	model.pageDown()

	assert.Equal(t, 2, model.CurrentPage())

	start, end := model.VisibleIndices()
	assert.Equal(t, 4, start)
	assert.Equal(t, 7, end, "page 2 should show rows 4-7")
	assert.Equal(t, 4, model.rowCursorIndex, "cursor should move to first row of page 2")
}

// TestTargetHeightMultilineVariableRows checks that the page map respects actual
// rendered row heights so that wrapping rows push subsequent rows onto the next page.
// targetHeight=8: pass 1 (no footer) availableLines=4 → 2 pages, triggers pass 2.
// Pass 2 (footer) availableLines=2:
//   - short(1)+short(2) = 2 lines → page 1 full
//   - long(3) = 2 lines → page 2 (alone, no room for short(4))
//   - short(4)+short(5) = 2 lines → page 3
func TestTargetHeightMultilineVariableRows(t *testing.T) {
	rows := []Row{
		shortRow(1),
		shortRow(2),
		longRow(3),
		shortRow(4),
		shortRow(5),
	}

	model := genTargetHeightTable(8, true, rows)

	require.Equal(t, 3, model.MaxPages(), "two-pass with footer budget: 3 pages")

	start, end := model.VisibleIndices()
	assert.Equal(t, 0, start)
	assert.Equal(t, 1, end, "page 1 should contain rows 0-1")

	model.pageDown()

	start, end = model.VisibleIndices()
	assert.Equal(t, 2, start)
	assert.Equal(t, 2, end, "page 2 should contain only the long row")

	model.pageDown()

	start, end = model.VisibleIndices()
	assert.Equal(t, 3, start)
	assert.Equal(t, 4, end, "page 3 should contain rows 3-4")
}

// TestTargetHeightCursorStaysInBoundsAcrossPages checks that cursor navigation
// stays within the visible range after crossing a page boundary.
func TestTargetHeightCursorStaysInBoundsAcrossPages(t *testing.T) {
	rows := make([]Row, 10)
	for i := range rows {
		rows[i] = shortRow(i + 1)
	}

	model := genTargetHeightTable(8, false, rows)

	// Move down past the end of page 1 (4 rows) and into page 2.
	for i := 0; i < 5; i++ {
		model.moveHighlightDown()
	}

	start, end := model.VisibleIndices()
	assert.GreaterOrEqual(t, model.rowCursorIndex, start)
	assert.LessOrEqual(t, model.rowCursorIndex, end)
}

// TestTargetHeightExpectedPageForRowIndex checks the binary search page lookup.
// Pages start at [0, 4, 8] for 12 single-line rows with availableLines=4.
func TestTargetHeightExpectedPageForRowIndex(t *testing.T) {
	rows := make([]Row, 12)
	for i := range rows {
		rows[i] = shortRow(i + 1)
	}

	model := genTargetHeightTable(10, false, rows)

	assert.Equal(t, 0, model.expectedPageForRowIndex(0))
	assert.Equal(t, 0, model.expectedPageForRowIndex(3))
	assert.Equal(t, 1, model.expectedPageForRowIndex(4))
	assert.Equal(t, 1, model.expectedPageForRowIndex(7))
	assert.Equal(t, 2, model.expectedPageForRowIndex(8))
	assert.Equal(t, 2, model.expectedPageForRowIndex(11))
}

// TestTargetHeightPageMapRebuildOnWithRows checks that changing rows invalidates
// and rebuilds the page map.
// Uses 7 rows: pass 1 (avail=6) gives 2 pages, triggering pass 2 (avail=4, 2 pages).
func TestTargetHeightPageMapRebuildOnWithRows(t *testing.T) {
	rows := make([]Row, 7)
	for i := range rows {
		rows[i] = shortRow(i + 1)
	}

	model := genTargetHeightTable(10, false, rows)
	assert.Equal(t, 2, model.MaxPages())

	newRows := make([]Row, 3)
	for i := range newRows {
		newRows[i] = shortRow(i + 1)
	}

	model = model.WithRows(newRows)
	assert.Equal(t, 1, model.MaxPages())
}

// TestTargetHeightSingleRowTallerThanBudget checks that a row taller than the
// entire available height still appears alone on its page rather than being dropped.
func TestTargetHeightSingleRowTallerThanBudget(t *testing.T) {
	// targetHeight=4 → availableLines = 4-3-1 = 0, clamped to 1.
	// The long row wraps to 2 lines — exceeds the budget.
	// It must still appear alone on page 1; short(2) goes to page 2.
	rows := []Row{
		longRow(1),
		shortRow(2),
	}

	model := genTargetHeightTable(4, true, rows)

	assert.Equal(t, 2, model.MaxPages())

	start, end := model.VisibleIndices()
	assert.Equal(t, 0, start)
	assert.Equal(t, 0, end, "oversized row should occupy page 1 alone")
}

// TestTargetHeightPageLast checks pageLast navigates to the correct page and cursor.
func TestTargetHeightPageLast(t *testing.T) {
	rows := make([]Row, 12)
	for i := range rows {
		rows[i] = shortRow(i + 1)
	}

	model := genTargetHeightTable(10, false, rows)
	model.pageLast()

	assert.Equal(t, 3, model.CurrentPage())

	start, end := model.VisibleIndices()
	assert.Equal(t, 8, start)
	assert.Equal(t, 11, end)
	assert.Equal(t, 8, model.rowCursorIndex)
}

// TestTargetHeightWithCurrentPage checks the public WithCurrentPage API.
func TestTargetHeightWithCurrentPage(t *testing.T) {
	rows := make([]Row, 12)
	for i := range rows {
		rows[i] = shortRow(i + 1)
	}

	model := genTargetHeightTable(10, false, rows)
	model = model.WithCurrentPage(3)

	assert.Equal(t, 3, model.CurrentPage())
	assert.Equal(t, 8, model.rowCursorIndex)
}

// TestTargetHeightRowSeparatorReducesBudget checks that row separators consume
// budget lines, resulting in fewer rows visible per page.
func TestTargetHeightRowSeparatorReducesBudget(t *testing.T) {
	rows := make([]Row, 10)
	for i := range rows {
		rows[i] = shortRow(i + 1)
	}

	modelPlain := genTargetHeightTable(8, false, rows)
	modelSep := genTargetHeightTable(8, false, rows).WithRowBorder(true)

	_, endPlain := modelPlain.VisibleIndices()
	_, endSep := modelSep.VisibleIndices()

	assert.Greater(t, endPlain, endSep, "separators should reduce rows per page")
}

// TestTargetHeightViewHeightRespected renders the table and confirms the output
// does not exceed targetHeight lines.
func TestTargetHeightViewHeightRespected(t *testing.T) {
	rows := []Row{
		longRow(1),
		longRow(2),
		longRow(3),
		longRow(4),
		longRow(5),
	}

	const target = 10

	model := genTargetHeightTable(target, true, rows)
	rendered := model.View()

	lines := strings.Split(rendered, "\n")
	assert.LessOrEqual(t, len(lines), target,
		"rendered height (%d lines) exceeded target (%d)", len(lines), target)
}

// TestTargetHeightPageUp checks that pageUp navigates from page 2 back to page 1.
func TestTargetHeightPageUp(t *testing.T) {
	rows := make([]Row, 12)
	for i := range rows {
		rows[i] = shortRow(i + 1)
	}

	model := genTargetHeightTable(10, false, rows)
	model.pageDown()
	assert.Equal(t, 2, model.CurrentPage())

	model.pageUp()

	assert.Equal(t, 1, model.CurrentPage())

	start, end := model.VisibleIndices()
	assert.Equal(t, 0, start)
	assert.Equal(t, 3, end, "page 1 should show rows 0-3 after paging up")
	assert.Equal(t, 0, model.rowCursorIndex, "cursor should return to first row of page 1")
}

// TestTargetHeightPageDownSinglePage checks that pageDown is a no-op when all
// rows fit on a single page.
func TestTargetHeightPageDownSinglePage(t *testing.T) {
	rows := []Row{shortRow(1), shortRow(2)}

	model := genTargetHeightTable(8, false, rows)
	assert.Equal(t, 1, model.MaxPages())

	model.pageDown()

	assert.Equal(t, 1, model.CurrentPage(), "pageDown on a single page should not change page")
}

// TestTargetHeightPageUpSinglePage checks that pageUp is a no-op when all
// rows fit on a single page.
func TestTargetHeightPageUpSinglePage(t *testing.T) {
	rows := []Row{shortRow(1), shortRow(2)}

	model := genTargetHeightTable(8, false, rows)
	assert.Equal(t, 1, model.MaxPages())

	model.pageUp()

	assert.Equal(t, 1, model.CurrentPage(), "pageUp on a single page should not change page")
}

// TestTargetHeightWithMultilineInvalidatesPageMap checks that calling WithMultiline
// on a model that already has targetHeight set invalidates and rebuilds the page map.
// With availableLines=4 (targetHeight=10, metaHeight=5, bottom border=1) and 4 long rows:
//   - multiline off: each row is 1 line → all 4 fit on page 1 → 1 page
//   - multiline on:  each row wraps to 2 lines → only 2 fit per page → 2 pages
func TestTargetHeightWithMultilineInvalidatesPageMap(t *testing.T) {
	rows := []Row{longRow(1), longRow(2), longRow(3), longRow(4)}

	model := genTargetHeightTable(10, false, rows)
	assert.Equal(t, 1, model.MaxPages(), "multiline off: all 4 rows fit on one page")

	model = model.WithMultiline(true)
	assert.Equal(t, 2, model.MaxPages(), "multiline on: wrapping rows should split across 2 pages")
}

// TestTargetHeightEmptyRows checks that a table with targetHeight and no rows
// returns sane defaults from MaxPages, VisibleIndices, and pageLast without panicking.
func TestTargetHeightEmptyRows(t *testing.T) {
	model := genTargetHeightTable(8, false, []Row{})

	assert.Equal(t, 1, model.MaxPages(), "empty table should report 1 page")

	start, end := model.VisibleIndices()
	assert.Equal(t, 0, start)
	assert.Equal(t, -1, end, "empty table should return -1 end index")

	// pageLast on an empty table should be a no-op (not panic)
	model.pageLast()
	assert.Equal(t, 1, model.CurrentPage())
}

// TestTargetHeightWithRowsShrinksClampsPage checks that calling WithRows with
// fewer rows than the current page requires correctly clamps to the last page.
func TestTargetHeightWithRowsShrinksClampsPage(t *testing.T) {
	rows := make([]Row, 12)
	for i := range rows {
		rows[i] = shortRow(i + 1)
	}

	model := genTargetHeightTable(10, false, rows)
	model.pageDown()
	assert.Equal(t, 2, model.CurrentPage(), "should be on page 2 before shrink")

	// Reduce to 3 rows (1 page) while on page 2 — WithRows should clamp to page 1.
	newRows := []Row{shortRow(1), shortRow(2), shortRow(3)}
	model = model.WithRows(newRows)

	assert.Equal(t, 1, model.CurrentPage(), "page should be clamped to last after row shrink")
}

// TestTargetHeightWithCurrentPageBelowOne checks that WithCurrentPage clamps
// values less than 1 to page 1 in targetHeight mode.
func TestTargetHeightWithCurrentPageBelowOne(t *testing.T) {
	rows := make([]Row, 12)
	for i := range rows {
		rows[i] = shortRow(i + 1)
	}

	model := genTargetHeightTable(8, false, rows)
	model = model.WithCurrentPage(0)

	assert.Equal(t, 1, model.CurrentPage(), "page below 1 should be clamped to 1")
	assert.Equal(t, 0, model.rowCursorIndex)
}

// TestTargetHeightVisibleDataLinesWithSeparator calls View() on a table with
// both multiline and row separators enabled to exercise the separator branch in
// visibleDataLines.
func TestTargetHeightVisibleDataLinesWithSeparator(t *testing.T) {
	rows := []Row{
		shortRow(1),
		shortRow(2),
		shortRow(3),
	}

	const target = 12

	model := New([]Column{
		NewColumn("id", "ID", 3),
		NewColumn("content", "Content", 20),
	}).
		WithRows(rows).
		WithMultiline(true).
		WithRowBorder(true).
		WithTargetHeight(target).
		WithMinimumHeight(target)

	rendered := model.View()
	lines := strings.Split(rendered, "\n")

	assert.LessOrEqual(t, len(lines), target,
		"rendered height (%d lines) exceeded target (%d)", len(lines), target)
	assert.Greater(t, len(lines), 0, "table should not be empty")
}

// TestTargetHeightNoFooterOnSinglePage checks that no page-count footer is
// rendered when all rows fit on a single page.
func TestTargetHeightNoFooterOnSinglePage(t *testing.T) {
	rows := []Row{shortRow(1), shortRow(2)}

	model := genTargetHeightTable(10, false, rows)

	assert.Equal(t, 1, model.MaxPages(), "all rows should fit on one page")

	rendered := model.View()
	assert.NotContains(t, rendered, "1/1", "page-count footer must not appear on a single-page table")
}

// TestTargetHeightPageMapInvalidatedByOptions checks that options which affect
// layout correctly invalidate the page map so MaxPages stays accurate.
func TestTargetHeightPageMapInvalidatedByOptions(t *testing.T) {
	makeRows := func(n int) []Row {
		rows := make([]Row, n)
		for i := range rows {
			rows[i] = shortRow(i + 1)
		}

		return rows
	}

	base := genTargetHeightTable(10, false, makeRows(12))
	basePages := base.MaxPages()

	require.Greater(t, basePages, 1, "test setup requires a multi-page table")

	// Each of these options should invalidate the page map; MaxPages must still
	// return a sensible (non-zero) value after the rebuild.
	cases := []struct {
		name  string
		apply func(Model) Model
	}{
		{"Filtered", func(m Model) Model { return m.Filtered(true) }},
		{"WithStaticFooter", func(m Model) Model { return m.WithStaticFooter("hint") }},
		{"WithTargetWidth", func(m Model) Model { return m.WithTargetWidth(50) }},
		{"WithColumns", func(m Model) Model {
			return m.WithColumns([]Column{NewColumn("id", "ID", 3), NewColumn("content", "Content", 20)})
		}},
		{"WithFooterVisibility", func(m Model) Model { return m.WithFooterVisibility(false) }},
		{"WithHeaderVisibility", func(m Model) Model { return m.WithHeaderVisibility(false) }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			model := tc.apply(base)
			assert.Greater(t, model.MaxPages(), 0, "MaxPages should be positive after %s", tc.name)
		})
	}
}

// TestTargetHeightFilterReducesRowsBelowOnePage is a regression test for a
// panic caused by pageStartIndices not being invalidated when a filter is
// applied in targetHeight mode.  The page map is built by pointer-receiver
// navigation calls (e.g. pageDown) and cached on the model.  A subsequent
// filter that shrinks visible rows leaves the stale map in place; ensurePageMap
// then skips rebuilding, VisibleIndices returns an endRowIndex beyond the
// filtered slice, and renderRow panics with an index out of range.
func TestTargetHeightFilterReducesRowsBelowOnePage(t *testing.T) {
	// targetHeight=10 gives 4 rows per page (metaHeight=5, 1 border line).
	// 8 rows spans two pages; a navigation call seeds the stale page map;
	// the filter then reduces visible rows to 1 (< one page).
	rows := make([]Row, 8)
	for i := range rows {
		rows[i] = NewRow(RowData{
			"id":      i + 1,
			"content": fmt.Sprintf("item-%d", i+1),
		})
	}

	model := New([]Column{
		NewColumn("id", "ID", 3),
		NewColumn("content", "Content", 20).WithFiltered(true),
	}).
		WithRows(rows).
		WithTargetHeight(10).
		Filtered(true)

	// pageDown is a pointer receiver: it calls ensurePageMap on the model
	// directly, seeding pageStartIndices for 8 unfiltered rows ([0, 4]).
	model.pageDown()

	// Applying the filter sets visibleRowCacheUpdated=false but (before the
	// fix) leaves pageStartIndices=[0,4] intact.
	model = model.WithFilterInputValue("item-3")

	// View must not panic even though only 1 row is visible.
	assert.NotPanics(t, func() {
		_ = model.View()
	})
}

// TestTargetHeightFilterViaWithFilterInput ensures View does not panic when
// WithFilterInput reduces visible rows below one page in targetHeight mode.
func TestTargetHeightFilterViaWithFilterInput(t *testing.T) {
	rows := make([]Row, 8)
	for i := range rows {
		rows[i] = NewRow(RowData{"id": i + 1, "content": fmt.Sprintf("item-%d", i+1)})
	}

	model := New([]Column{
		NewColumn("id", "ID", 3),
		NewColumn("content", "Content", 20).WithFiltered(true),
	}).WithRows(rows).WithTargetHeight(10).Filtered(true)

	model.pageDown()

	input := textinput.New()
	input.SetValue("item-3")
	model = model.WithFilterInput(input)

	assert.NotPanics(t, func() { _ = model.View() })
}

// TestTargetHeightFilterViaUpdateFilterTextInput ensures View does not panic
// when typing in the filter input reduces visible rows below one page in
// targetHeight mode.
func TestTargetHeightFilterViaUpdateFilterTextInput(t *testing.T) {
	rows := make([]Row, 8)
	for i := range rows {
		rows[i] = NewRow(RowData{"id": i + 1, "content": fmt.Sprintf("item-%d", i+1)})
	}

	model := New([]Column{
		NewColumn("id", "ID", 3),
		NewColumn("content", "Content", 20).WithFiltered(true),
	}).WithRows(rows).WithTargetHeight(10).Filtered(true)

	model.pageDown()

	model.filterTextInput.SetValue("item-3")
	model, _ = model.updateFilterTextInput(tea.KeyPressMsg{})

	assert.NotPanics(t, func() { _ = model.View() })
}

// TestTargetHeightFilterClearKey ensures View does not panic after the filter
// is cleared via the FilterClear key in targetHeight mode.
func TestTargetHeightFilterClearKey(t *testing.T) {
	rows := make([]Row, 8)
	for i := range rows {
		rows[i] = NewRow(RowData{"id": i + 1, "content": fmt.Sprintf("item-%d", i+1)})
	}

	model := New([]Column{
		NewColumn("id", "ID", 3),
		NewColumn("content", "Content", 20).WithFiltered(true),
	}).WithRows(rows).WithTargetHeight(10).Filtered(true)

	model.filterTextInput.SetValue("item-3")
	model.visibleRowCacheUpdated = false
	model.pageDown()

	model.handleKeypress(tea.KeyPressMsg{Code: tea.KeyEscape})

	assert.NotPanics(t, func() { _ = model.View() })
}

// TestTargetHeightWithCurrentPageAboveMax checks that WithCurrentPage clamps
// values above MaxPages to the last page in targetHeight mode.
func TestTargetHeightWithCurrentPageAboveMax(t *testing.T) {
	rows := make([]Row, 12)
	for i := range rows {
		rows[i] = shortRow(i + 1)
	}

	model := genTargetHeightTable(10, false, rows)

	require.Greater(t, model.MaxPages(), 1, "test setup requires a multi-page table")

	model = model.WithCurrentPage(999)

	assert.Equal(t, model.MaxPages(), model.CurrentPage(), "page above max should be clamped to last page")
}
