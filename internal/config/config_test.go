package config_test

import (
	"testing"
	"time"
	"trader/internal/config"
)

func TestGetConfigOk(t *testing.T) {
	// THEN
	config := config.NewConfig()

	// WHEN
	if config.GetVersion() != "development" {
		t.Errorf(`expected at "%s" and received at %s`, "development", config.GetVersion())
	}
	if config.GetScrapingTimeout() != time.Duration(60) {
		t.Errorf(`expected at "%s" and received at %s`, time.Duration(60), config.GetScrapingTimeout())
	}
}
