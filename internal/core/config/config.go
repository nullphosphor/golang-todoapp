package core_config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	TimeZone *time.Location
}

func NewConfigMust() *Config {
	tz := os.Getenv("TIME_ZONE")
	if tz == "" {
		tz = "UTC"
	}

	zone, err := time.LoadLocation(tz)
	if err != nil {
		err = fmt.Errorf("load time zone: %s: %w", tz, err)
		panic(err)
	}

	return &Config{
		TimeZone: zone,
	}
}
