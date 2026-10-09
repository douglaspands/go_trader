package common

import (
	"io"

	"github.com/jedib0t/go-pretty/v6/table"
)

func NewTableWriter(noColor bool, out io.Writer) table.Writer {
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
	t.SetOutputMirror(out)
	t.SetStyle(style)
	return t
}
