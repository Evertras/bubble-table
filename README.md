# Bubble-table

<p>
  <a href="https://github.com/Evertras/bubble-table/releases"><img src="https://img.shields.io/github/release/Evertras/bubble-table.svg" alt="Latest Release"></a>
  <a href="https://pkg.go.dev/github.com/evertras/bubble-table/table?tab=doc"><img src="https://godoc.org/github.com/golang/gddo?status.svg" alt="GoDoc"></a>
  <a href='https://coveralls.io/github/Evertras/bubble-table?branch=main'><img src='https://coveralls.io/repos/github/Evertras/bubble-table/badge.svg?branch=main&hash=abc' alt='Coverage Status'/></a>
</p>

A customizable, interactive table component for the
[Bubble Tea framework](https://github.com/charmbracelet/bubbletea).

![Styled table](https://user-images.githubusercontent.com/5923958/156778142-cc1a32e1-1b1e-4a65-b699-187f39f0f946.png)

[View above sample source code](./examples/pokemon/main.go)

## Contributing

Contributions welcome, please [check the contributions doc](./CONTRIBUTING.md)
for a few helpful tips!

## Features

For a code reference of all available features, please see the [full feature example](./examples/features/main.go).
If you want to get started with a simple default table, [check the simplest example](./examples/simplest/main.go).

Displays a table with a header, rows, footer, and borders.

Columns can be fixed-width [or flexible width](./examples/flex/main.go).

Border shape is customizable with a basic thick square default.  The color can
be modified by applying a base style with `lipgloss.NewStyle().BorderForeground(...)`.

Styles can be applied globally and to columns, rows, and individual cells.
The base style is applied first, then column, then row, then cell when
determining overrides.  The default base style is a basic right-alignment.
[See the main feature example](./examples/features/main.go) to see styles and
how they override each other.

Can be focused to highlight a row and navigate with up/down (and j/k).  These
keys can be customized with a KeyMap.

Can make rows selectable, and fetch the current selections.

Pagination can be set with a given page size, which automatically generates a
simple footer to show the current page and total pages.

Built-in filtering can be enabled by setting any columns as filterable, using
a text box in the footer and `/` (customizable by keybind) to start filtering.

Columns can be sorted in either ascending or descending order.  Multiple columns
can be specified in a row.  If multiple columns are specified, first the table
is sorted by the first specified column, then each group within that column is
sorted in smaller and smaller groups.  [See the sorting example](examples/sorting/main.go)
for more information.  If a column contains numbers (either ints or floats),
the numbers will be sorted by numeric value.  Otherwise rendered string values
will be compared.

If a feature is confusing to use or could use a better example, please feel free
to open an issue.

## Defining table data

A table is defined by a list of `Column` values that define the columns in the
table.  Each `Column` is associated with a unique string key.

A table contains a list of `Row`s.  Each `Row` contains a `RowData` object which
is simply a map of string column IDs to arbitrary `interface{}` data values.
When the table is rendered, each `Row` is checked for each `Column` key.  If the
key exists in the `Row`'s `RowData`, it is rendered with `fmt.Sprintf("%v")`.
If it does not exist, nothing is rendered.

Extra data in the `RowData` object is ignored.  This can be helpful to simply
dump data into `RowData` and create columns that select what is interesting to
view, or to generate different columns based on view options on the fly (see the
[metadata example](./examples/metadata/main.go) for an example of using this).

An example is given below.  For more detailed examples, see
[the examples directory](./examples).

```golang
// This makes it easier/safer to match against values, but isn't necessary
const (
  // This value isn't visible anywhere, so a simple lowercase is fine
  columnKeyID = "id"

  // It's just a string, so it can be whatever, really!  They only must be unique
  columnKeyName = "何?!"
)

// Note that there's nothing special about "ID" or "Name", these are completely
// arbitrary columns
columns := []table.Column{
  table.NewColumn(columnKeyID, "ID", 5),
  table.NewColumn(columnKeyName, "Name", 10),
}

rows := []table.Row{
  // This row contains both an ID and a name
  table.NewRow(table.RowData{
    columnKeyID:          "abc",
    columnKeyName:        "Hello",
  }),

  table.NewRow(table.RowData{
    columnKeyID:          "123",
    columnKeyName:        "Oh no",
    // This field exists in the row data but won't be visible
    "somethingelse": "Super bold!",
  }),

  table.NewRow(table.RowData{
    columnKeyID:          "def",
    // This row is missing the Name column, so it will simply be blank
  }),

  // We can also apply styling to the row or to individual cells

  // This row has individual styling to make it bold
  table.NewRow(table.RowData{
    columnKeyID:          "bold",
    columnKeyName:        "Bolded",
  }).WithStyle(lipgloss.NewStyle().Bold(true),

  // This row also has individual styling to make it bold
  table.NewRow(table.RowData{
    columnKeyID:          "alert",
    // This cell has styling applied on top of the bold
    columnKeyName:        table.NewStyledCell("Alert", lipgloss.NewStyle().Foreground(lipgloss.Color("#f88"))),
  }).WithStyle(lipgloss.NewStyle().Bold(true),
}
```

## Demos

Code examples are located in [the examples directory](./examples).  Run commands
are added to the [Makefile](Makefile) for convenience but they should be as
simple as `go run ./examples/features/main.go`, etc.

To run the examples, clone this repo and run:

```bash
# Run the full feature demo
make

# Run dimensions example to see multiple sizes of simple tables in action
make example-dimensions

# Or run any of them directly
go run ./examples/pagination/main.go
```

