package statelock

import "time"

// State represents the application state
type State string

const (
	// StateOK allows all API operations
	StateOK State = "OK"
	// StateSoftLock allows only GET operations
	StateSoftLock State = "SOFT_LOCK"
	// StateHardLock blocks all API operations
	StateHardLock State = "HARD_LOCK"
)

// Config holds the configuration for the state lock module
type Config struct {
	// ScheduleInterval is the interval for checking state (e.g., "@every 30s", "0 */5 * * * *")
	ScheduleInterval string
	// StateEndpoint is the URL to fetch the current state
	StateEndpoint string
	// TimeoutDuration for HTTP requests
	TimeoutDuration time.Duration
	// OnStateChange is a callback when state changes
	OnStateChange func(oldState, newState State)
}

// StateResponse represents the expected response from state endpoint
type StateResponse struct {
	State   State  `json:"state"`
	Message string `json:"message,omitempty"`
}

// DefaultConfig returns a config with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		ScheduleInterval: "@every 30s",
		TimeoutDuration:  10 * time.Second,
	}
}
