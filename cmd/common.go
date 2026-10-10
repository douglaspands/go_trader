package cmd

import (
	"errors"
	"fmt"
	"io"
	"math"
	"strings"
	"trader/internal/resource"
	"trader/internal/tools"

	"github.com/jedib0t/go-pretty/v6/table"
)

// maxAmount is the largest amount accepted by purchase-balance.
const maxAmount = 1e12

func validateAmount(amount float64) error {
	if math.IsNaN(amount) || math.IsInf(amount, 0) || amount <= 0 || amount > maxAmount {
		return fmt.Errorf("invalid amount %v: it must be a finite number greater than 0 and no larger than %.0f", amount, maxAmount)
	}
	return nil
}

// cell prepares text obtained from the provider to be written as a table or CSV cell.
func cell(text string, csv bool) string {
	text = tools.CleanText(text)
	if csv {
		text = tools.CSVSafe(text)
	}
	return text
}

// printWarnings writes one line per ticker that failed in a list that still has a result.
func printWarnings(w io.Writer, failures []*resource.TickerFailure) {
	for _, failure := range failures {
		fmt.Fprintf(w, "Warning: %s: %s\n", tools.CleanText(failure.Ticker), tools.CleanText(failure.Err.Error()))
	}
}

// cleanError removes control characters from each line of the message of err.
func cleanError(err error) error {
	lines := strings.Split(err.Error(), "\n")
	for i, line := range lines {
		lines[i] = tools.CleanText(line)
	}
	return errors.New(strings.Join(lines, "\n"))
}

func render(t table.Writer, csv bool) {
	if csv {
		t.RenderCSV()
	} else {
		t.Render()
	}
}
