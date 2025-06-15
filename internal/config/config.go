package config

import "time"

type Config interface {
	GetVersion() string
	GetScrapingTimeout() time.Duration
}

type config struct {
	version            string
	scrapingTimeoutTtl int64
}

func (c *config) GetVersion() string {
	return c.version
}

func (c *config) GetScrapingTimeout() time.Duration {
	return time.Duration(c.scrapingTimeoutTtl)
}

var version string = "development"

func NewConfig() Config {
	return &config{
		version:            version,
		scrapingTimeoutTtl: 60,
	}
}
