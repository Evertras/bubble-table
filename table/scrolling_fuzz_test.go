//go:build go1.18
// +build go1.18

package table

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

// This is long because of test cases
//
//nolint:funlen,gocognit,cyclop
func FuzzHorizontalScrollingStopEdgeCases(f *testing.F) {
	const (
		minNameWidth = 2
		maxNameWidth = 50

		minColWidth = 4
		maxColWidth = 50

		minNumCols = 1
		maxNumCols = 500

		minMaxWidth = 5
		maxMaxWidth = 200

		borderBuffer = 4
	)

	f.Add(5, 3, 5, 30)
	f.Fuzz(func(t *testing.T, nameWidth, colWidth, numCols, maxWidth int) {
		if nameWidth < minNameWidth ||
			nameWidth > maxNameWidth ||
			nameWidth > maxWidth-colWidth ||
			nameWidth+colWidth+borderBuffer >= maxWidth {
			return
		}

		if colWidth < minColWidth ||
			colWidth > maxColWidth ||
			colWidth >= maxWidth {
			return
		}

		if numCols < minNumCols || numCols > maxNumCols {
			return
		}

		if maxWidth < minMaxWidth || maxWidth > maxMaxWidth {
			return
		}

		cols := []Column{NewColumn("Name", "Name", nameWidth)}
		for i := 0; i < numCols; i++ {
			s := fmt.Sprintf("%d", i+1)
			cols = append(cols, NewColumn(s, s, colWidth))
		}

		rowData := RowData{"Name": "A"}

		for i := 0; i < numCols; i++ {
			s := fmt.Sprintf("%d", i+1)
			rowData[s] = s
		}

		rows := []Row{NewRow(rowData)}

		model := New(cols).
			WithRows(rows).
			WithStaticFooter("Footer").
			WithMaxTotalWidth(maxWidth).
			WithHorizontalFreezeColumnCount(1).
			Focused(true)

		hitScrollRight := func() {
			model, _ = model.Update(tea.KeyMsg{Type: tea.KeyShiftRight})
		}

		// Excessive scrolling attempts to be sure
		for i := 0; i < numCols*2; i++ {
			hitScrollRight()
		}

		rendered := model.View()

		assert.NotContains(t, rendered, ">")

		if !strings.Contains(rendered, "…") {
			assert.Contains(t, rendered, fmt.Sprintf("%d", numCols))
		}
	})
}
