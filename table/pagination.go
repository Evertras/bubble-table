package table

// PageSize returns the current page size for the table, or 0 if there is no
// pagination enabled.
func (m *Model) PageSize() int {
	return m.pageSize
}

// CurrentPage returns the current page that the table is on, starting from an
// index of 1.
func (m *Model) CurrentPage() int {
	return m.currentPage + 1
}

// MaxPages returns the maximum number of pages that are visible.
func (m *Model) MaxPages() int {
	if m.targetHeight != 0 {
		m.ensurePageMap()

		if len(m.pageStartIndices) == 0 {
			return 1
		}

		return len(m.pageStartIndices)
	}

	totalRows := len(m.GetVisibleRows())

	if m.pageSize == 0 || totalRows == 0 {
		return 1
	}

	return (totalRows-1)/m.pageSize + 1
}

// TotalRows returns the current total row count of the table.  If the table is
// paginated, this is the total number of rows across all pages.
func (m *Model) TotalRows() int {
	return len(m.GetVisibleRows())
}

// VisibleIndices returns the current visible rows by their 0 based index.
// Useful for custom pagination footers.
func (m *Model) VisibleIndices() (start, end int) {
	totalRows := len(m.GetVisibleRows())

	if m.targetHeight != 0 {
		m.ensurePageMap()

		if totalRows == 0 || len(m.pageStartIndices) == 0 {
			return 0, -1
		}

		start = m.pageStartIndices[m.currentPage]

		if m.currentPage+1 < len(m.pageStartIndices) {
			end = m.pageStartIndices[m.currentPage+1] - 1
		} else {
			end = totalRows - 1
		}

		return start, end
	}

	if m.pageSize == 0 {
		start = 0
		end = totalRows - 1

		return start, end
	}

	start = m.pageSize * m.currentPage
	end = start + m.pageSize - 1

	if end >= totalRows {
		end = totalRows - 1
	}

	return start, end
}

func (m *Model) wrappedOrClamped(wrappedValue, clampedValue int) int {
	if m.paginationWrapping {
		return wrappedValue
	}

	return clampedValue
}

func (m *Model) clampCurrentPage(maxPageIndex int) {
	if m.currentPage > maxPageIndex {
		m.currentPage = m.wrappedOrClamped(0, maxPageIndex)
	} else if m.currentPage < 0 {
		m.currentPage = m.wrappedOrClamped(maxPageIndex, 0)
	}
}

func (m *Model) pageDownTargetHeight() {
	m.ensurePageMap()

	if len(m.pageStartIndices) <= 1 {
		return
	}

	m.currentPage++
	m.clampCurrentPage(len(m.pageStartIndices) - 1)
	m.rowCursorIndex = m.pageStartIndices[m.currentPage]
}

func (m *Model) pageUpTargetHeight() {
	m.ensurePageMap()

	if len(m.pageStartIndices) <= 1 {
		return
	}

	m.currentPage--
	m.clampCurrentPage(len(m.pageStartIndices) - 1)
	m.rowCursorIndex = m.pageStartIndices[m.currentPage]
}

func (m *Model) pageDown() {
	if m.targetHeight != 0 {
		m.pageDownTargetHeight()

		return
	}

	if m.pageSize == 0 || len(m.GetVisibleRows()) <= m.pageSize {
		return
	}

	m.currentPage++
	m.clampCurrentPage(m.MaxPages() - 1)
	m.rowCursorIndex = m.currentPage * m.pageSize
}

func (m *Model) pageUp() {
	if m.targetHeight != 0 {
		m.pageUpTargetHeight()

		return
	}

	if m.pageSize == 0 || len(m.GetVisibleRows()) <= m.pageSize {
		return
	}

	m.currentPage--
	m.clampCurrentPage(m.MaxPages() - 1)
	m.rowCursorIndex = m.currentPage * m.pageSize
}

func (m *Model) pageFirst() {
	m.currentPage = 0
	m.rowCursorIndex = 0
}

func (m *Model) pageLast() {
	if m.targetHeight != 0 {
		m.ensurePageMap()

		if len(m.pageStartIndices) == 0 {
			return
		}

		m.currentPage = len(m.pageStartIndices) - 1
		m.rowCursorIndex = m.pageStartIndices[m.currentPage]

		return
	}

	m.currentPage = m.MaxPages() - 1
	m.rowCursorIndex = m.currentPage * m.pageSize
}

func (m *Model) expectedPageForRowIndex(rowIndex int) int {
	if m.targetHeight != 0 {
		m.ensurePageMap()

		// Binary search: find the last page whose start index <= rowIndex.
		low, high := 0, len(m.pageStartIndices)-1

		for low < high {
			mid := (low + high + 1) / 2 //nolint:mnd // standard ceiling-midpoint formula
			if m.pageStartIndices[mid] <= rowIndex {
				low = mid
			} else {
				high = mid - 1
			}
		}

		return low
	}

	if m.pageSize == 0 {
		return 0
	}

	return rowIndex / m.pageSize
}
