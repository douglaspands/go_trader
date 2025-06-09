package common

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

func NewTableWriter() table.Writer {
	style := table.StyleColoredBlackOnBlueWhite
	// style.Options = table.OptionsNoBordersAndSeparators
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(style)
	return t
}
