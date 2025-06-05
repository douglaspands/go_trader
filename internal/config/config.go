package config

type Config struct {
	Version            string
	ScrapingTimeoutTtl int8
}

var version string = "development"
var config *Config

func GetConfig() *Config {
	if config == nil {
		config = &Config{
			Version:            version,
			ScrapingTimeoutTtl: 60,
		}
	}
	return config
}
