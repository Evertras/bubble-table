# Bubble-table

<p>
  <a href="https://github.com/Evertras/bubble-table/releases"><img src="https://img.shields.io/github/release/Evertras/bubble-table.svg" alt="Latest Release"></a>
  <a href="https://pkg.go.dev/github.com/evertras/bubble-table/table?tab=doc"><img src="https://godoc.org/github.com/golang/gddo?status.svg" alt="GoDoc"></a>
  <a href='https://coveralls.io/github/Evertras/bubble-table?branch=main'><img src='https://coveralls.io/repos/github/Evertras/bubble-table/badge.svg?branch=main' alt='Coverage Status'/></a>
</p>

A table component for the [Bubble Tea framework](https://github.com/charmbracelet/bubbletea).

![Table with all features](https://user-images.githubusercontent.com/5923958/154851497-78a28031-db82-4a95-a4c7-dae5b0d6cf00.png)

[View sample source code](./examples/features/main.go)

## Features

For a code reference, please see the [full feature example](./examples/features/main.go).

Displays a table with a header, rows, footer, and borders.

Border shape is customizable with a basic thick square default.

Styles can be applied to columns, rows, and individual cells.  Column style is
applied first, then row, then cell when determining overrides.

Can be focused to highlight a row and navigate with up/down (and j/k).  These
keys can be customized with a KeyMap.

Can make rows selectable, and fetch the current selections.

Pagination can be set with a given page size, which automatically generates a
simple footer to show the current page and total pages.

## Defining table data

Each `Column` is associated with a unique string key.  Each `Row` contains a
`RowData` object which is simply a map of strings to `interface{}`.  When the
table is rendered, each `Row` is checked for each `Column` key.  If the key
exists in the `Row`'s `RowData`, it is rendered with `fmt.Sprintf("%v")`.  If it
does not exist, nothing is rendered.

Extra data in the `RowData` object is ignored.  This can be helpful to simply
dump data into `RowData` and create columns that select what is interesting to
view, or to generate different columns based on view options on the fly.

A basic example is given below.  For more detailed examples, see
[the examples directory](./examples).

```golang
// This makes it easier/safer to match against values, but isn't necessary
const (
  // This value isn't visible anywhere, so a simple lowercase is fine
  columnKeyID = "id"

  // It's just a string, so it can be whatever, really!  They only must be unique
  columnKeyName = "何?!"
)

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

  // This row has individual styling to make it bold
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

## Contributing

Contributions welcome, but since this is being actively developed for use in
[Khan](https://github.com/evertras/khan) please check first by opening an issue
or commenting on an existing one!

