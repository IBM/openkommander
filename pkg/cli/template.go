package cli

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

// RenderTable is a utility function to render a table with a dynamic header and rows.
// It provides a unified look for CLI commands.
func RenderTable(title string, headers []string, rows [][]interface{}) {
	if title != "" {
		fmt.Println("\n" + title)
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	headerRow := table.Row{}
	for _, header := range headers {
		headerRow = append(headerRow, header)
	}
	t.AppendHeader(headerRow)

	for _, row := range rows {
		t.AppendRow(table.Row(row))
	}

	t.SetStyle(table.StyleLight)
	t.Render()
}
