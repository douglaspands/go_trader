package common_test

import (
	"bytes"
	"strings"
	"testing"
	"trader/internal/common"

	"github.com/jedib0t/go-pretty/v6/table"
)

func renderSample(noColor bool) string {
	var out bytes.Buffer
	t := common.NewTableWriter(noColor, &out)
	t.AppendHeader(table.Row{"FIELD", "VALUE"})
	t.AppendRow(table.Row{"Ticker", "PETR4"})
	t.Render()
	return out.String()
}

func TestNewTableWriterNoColorOk(t *testing.T) {
	// THEN
	result := renderSample(true)

	// WHEN
	if !strings.Contains(result, "FIELD") || !strings.Contains(result, "PETR4") {
		t.Errorf("expected the table in the given writer and received %q", result)
	}
	if strings.Contains(result, "\x1b[") {
		t.Errorf("expected no ANSI escape codes and received %q", result)
	}
	if strings.Contains(result, "+") || strings.Contains(result, "|") {
		t.Errorf("expected no borders and received %q", result)
	}
}

func TestNewTableWriterColorOk(t *testing.T) {
	// THEN
	result := renderSample(false)

	// WHEN
	if !strings.Contains(result, "FIELD") || !strings.Contains(result, "PETR4") {
		t.Errorf("expected the table in the given writer and received %q", result)
	}
	if !strings.Contains(result, "\x1b[") {
		t.Errorf("expected ANSI escape codes and received %q", result)
	}
}
