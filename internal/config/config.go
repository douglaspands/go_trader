package config

type Config struct {
	ScrapingTimeoutTtl int8
}

var config *Config

func GetConfig() *Config {
	if config == nil {
		config = &Config{ScrapingTimeoutTtl: 60}
	}
	return config
}
