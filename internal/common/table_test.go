package common_test

import (
	"testing"
	"trader/internal/common"

	"github.com/jedib0t/go-pretty/v6/table"
)

func TestNewTableWriterNoColorOk(t *testing.T) {
	// THEN
	result := common.NewTableWriter(true)

	// WHEN
	switch result.(type) {
	case table.Writer:
		break
	default:
		t.Errorf(`variable type is not table.Writer`)
	}
}

func TestNewTableWriterColorOk(t *testing.T) {
	// THEN
	result := common.NewTableWriter(false)

	// WHEN
	switch result.(type) {
	case table.Writer:
		break
	default:
		t.Errorf(`variable type is not table.Writer`)
	}
}
