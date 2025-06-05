package common

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

func NewTableWriter() table.Writer {
	style := table.StyleDefault
	style.Options = table.OptionsNoBordersAndSeparators
	// width, _, _ := term.GetSize(int(os.Stdout.Fd()))
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(style)
	// t.SetAllowedRowLength(width)
	return t
}
