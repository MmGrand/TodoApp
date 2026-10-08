package core_logger

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

const (
	FormatConsole = "console"
	FormatJSON    = "json"
)

type Config struct {
	Level  string `envconfig:"LEVEL" default:"DEBUG"`
	Folder string `envconfig:"FOLDER" required:"true"`
	Format string `envconfig:"FORMAT" default:"console"`

	MaxSizeMB  int `envconfig:"MAX_SIZE_MB" default:"100"`
	MaxBackups int `envconfig:"MAX_BACKUPS" default:"10"`
	MaxAgeDays int `envconfig:"MAX_AGE_DAYS" default:"30"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("LOGGER", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	if config.Format != FormatConsole && config.Format != FormatJSON {
		return Config{}, fmt.Errorf(
			"unknown log format %q, expected %q or %q",
			config.Format,
			FormatConsole,
			FormatJSON,
		)
	}

	return config, nil
}
