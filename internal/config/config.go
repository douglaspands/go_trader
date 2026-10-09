package config

import "time"

type Config interface {
	GetVersion() string
	GetScrapingTimeout() time.Duration
}

type config struct {
	version         string
	scrapingTimeout time.Duration
}

func (c *config) GetVersion() string {
	return c.version
}

func (c *config) GetScrapingTimeout() time.Duration {
	return c.scrapingTimeout
}

var version string = "development"

func NewConfig() Config {
	return &config{
		version:         version,
		scrapingTimeout: 60 * time.Second,
	}
}
