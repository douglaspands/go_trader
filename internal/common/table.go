package common

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

func NewTableWriter(noColor bool) table.Writer {
	var style table.Style
	if noColor {
		style = table.StyleDefault
		style.Options = table.Options{
			DrawBorder:      false,
			SeparateColumns: false,
			SeparateHeader:  true,
			SeparateRows:    false,
			SeparateFooter:  true,
		}
	} else {
		style = table.StyleColoredBlackOnBlueWhite
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(style)
	return t
}
