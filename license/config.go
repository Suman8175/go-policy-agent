package license

import "time"

type Config struct {
	BaseURL         string
	RefreshInterval time.Duration
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() Config {
	return Config{
		BaseURL:         "http://localhost:8080",
		RefreshInterval: 30 * time.Minute,
	}
}
