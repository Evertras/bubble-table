package table

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	columnKeySelect = "___select___"
)

var (
	// defaultHighlightStyle has no visual effect; users opt into highlight colors
	// via WithHighlightStyle.  This matches the historical rendering behavior
	// where the table background color was not automatically applied to focused
	// rows.
	defaultHighlightStyle = lipgloss.NewStyle()
)

// Model is the main table model.  Create using New().
type Model struct {
	// Data
	columns  []Column
	rows     []Row
	metadata map[string]any

	// Caches for optimizations
	visibleRowCacheUpdated bool
	visibleRowCache        []Row

	// Shown when data is missing from a row
	missingDataIndicator any

	// Interaction
	focused bool
	keyMap  KeyMap

	// Taken from: 'Bubbles/List'
	// Additional key mappings for the short and full help views. This allows
	// you to add additional key mappings to the help menu without
	// re-implementing the help component. Of course, you can also disable the
	// list's help component and implement a new one if you need more
	// flexibility.
	// You have to supply a keybinding like this:
	// key.NewBinding( key.WithKeys("shift+left"), key.WithHelp("shift+←", "scroll left"))
	// It needs both 'WithKeys' and 'WithHelp'
	additionalShortHelpKeys func() []key.Binding
	additionalFullHelpKeys  func() []key.Binding

	selectableRows bool
	rowCursorIndex int

	// Events
	lastUpdateUserEvents []UserEvent

	// Styles
	baseStyle      lipgloss.Style
	highlightStyle lipgloss.Style
	headerStyle    lipgloss.Style
	rowStyleFunc   func(RowStyleFuncInput) lipgloss.Style
	border         Border
	selectedText   string
	unselectedText string

	// Header
	headerVisible bool

	// Footers
	footerVisible bool
	staticFooter  string

	// Pagination
	pageSize           int
	currentPage        int
	paginationWrapping bool

	// Sorting, where a stable sort is applied from first element to last so
	// that elements are grouped by the later elements.
	sortOrder []SortColumn

	// Filter
	filtered        bool
	filterTextInput textinput.Model
	filterFunc      FilterFunc

	// For flex columns
	targetTotalWidth int

	// The maximum total width for overflow/scrolling
	maxTotalWidth int

	// Internal cached calculations for reference, may be higher than
	// maxTotalWidth.  If this is the case, we need to adjust the view
	totalWidth int

	// How far to scroll to the right, in columns
	horizontalScrollOffsetCol int

	// How many columns to freeze when scrolling horizontally
	horizontalScrollFreezeColumnsCount int

	// Calculated maximum column we can scroll to before the last is displayed
	maxHorizontalColumnIndex int

	// Minimum total height of the table
	minimumHeight int

	// Target total height of the table in terminal lines, including borders,
	// header, and footer. When set, the table fits as many rows as possible
	// per page. Mutually exclusive with pageSize.
	targetHeight int

	// Cached page boundary indices when targetHeight is set.
	// pageStartIndices[i] is the index of the first visible row on page i.
	// Nil means the cache needs rebuilding.
	pageStartIndices []int

	// Internal cached calculation, the height of the header and footer
	// including borders. Used to determine how many padding rows to add.
	metaHeight int

	// If true, the table will be multiline
	multiline bool

	// If true, draw a horizontal separator line between each data row
	rowSeparator bool

	// If true, render the outer border (top line, bottom line, left/right cell borders)
	outerBorder bool
}

// New creates a new table ready for further modifications.
func New(columns []Column) Model {
	filterInput := textinput.New()
	filterInput.Prompt = "/"
	// Use plain styles without foreground colors so that the filter text
	// renders without ANSI color codes, keeping the output predictable.
	plainStyles := textinput.DefaultDarkStyles()
	plainStyles.Focused.Prompt = lipgloss.NewStyle()
	plainStyles.Focused.Text = lipgloss.NewStyle()
	plainStyles.Blurred.Prompt = lipgloss.NewStyle()
	plainStyles.Blurred.Text = lipgloss.NewStyle()
	plainStyles.Cursor.Color = lipgloss.NoColor{}
	filterInput.SetStyles(plainStyles)
	model := Model{
		columns:        make([]Column, len(columns)),
		metadata:       make(map[string]any),
		highlightStyle: defaultHighlightStyle,
		border:         borderDefault,
		outerBorder:    true,
		headerVisible:  true,
		footerVisible:  true,
		keyMap:         DefaultKeyMap(),

		selectedText:   "[x]",
		unselectedText: "[ ]",

		filterTextInput: filterInput,
		filterFunc:      filterFuncContains,
		baseStyle:       lipgloss.NewStyle().Align(lipgloss.Right),

		paginationWrapping: true,
	}

	// Do a full deep copy to avoid unexpected edits
	copy(model.columns, columns)

	model.recalculateWidth()

	return model
}

// Init initializes the table per the Bubble Tea architecture.
func (m Model) Init() tea.Cmd {
	return nil
}
