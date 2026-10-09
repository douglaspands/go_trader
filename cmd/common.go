package cmd

import "github.com/jedib0t/go-pretty/v6/table"

func render(t table.Writer, csv bool) {
	if csv {
		t.RenderCSV()
	} else {
		t.Render()
	}
}
