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
	if config.GetScrapingTimeout() != 60*time.Second {
		t.Errorf(`expected at "%s" and received at %s`, 60*time.Second, config.GetScrapingTimeout())
	}
	if config.GetMaxConcurrentRequests() != 4 {
		t.Errorf(`expected at "%d" and received at %d`, 4, config.GetMaxConcurrentRequests())
	}
}
