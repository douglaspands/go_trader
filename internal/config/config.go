package config

import "time"

type Config interface {
	GetVersion() string
	GetScrapingTimeout() time.Duration
	GetMaxConcurrentRequests() int
}

type config struct {
	version               string
	scrapingTimeout       time.Duration
	maxConcurrentRequests int
}

func (c *config) GetVersion() string {
	return c.version
}

func (c *config) GetScrapingTimeout() time.Duration {
	return c.scrapingTimeout
}

func (c *config) GetMaxConcurrentRequests() int {
	return c.maxConcurrentRequests
}

var version string = "development"

func NewConfig() Config {
	return &config{
		version:               version,
		scrapingTimeout:       60 * time.Second,
		maxConcurrentRequests: 4,
	}
}
