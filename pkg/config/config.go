package config

type Config struct {
	LogLevel int `mapstructure:"log_level"`
}

var Conf Config
