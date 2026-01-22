package statelock

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/robfig/cron/v3"
)

// Manager handles state checking and enforcement
type Manager struct {
	config       *Config
	currentState State
	mu           sync.RWMutex
	scheduler    *cron.Cron
	httpClient   *http.Client
}

// NewManager creates a new state lock manager
func NewManager(config *Config) (*Manager, error) {
	if config == nil {
		config = DefaultConfig()
	}

	if config.StateEndpoint == "" {
		return nil, fmt.Errorf("state endpoint is required")
	}

	m := &Manager{
		config:       config,
		currentState: StateOK, // Start with OK state
		scheduler:    cron.New(),
		httpClient: &http.Client{
			Timeout: config.TimeoutDuration,
		},
	}

	// Schedule the state check
	_, err := m.scheduler.AddFunc(config.ScheduleInterval, m.checkState)
	if err != nil {
		return nil, fmt.Errorf("failed to schedule state check: %w", err)
	}

	return m, nil
}

// Start begins the scheduled state checking
func (m *Manager) Start() {
	log.Println("State lock manager started")
	m.checkState() // Run immediately on start
	m.scheduler.Start()
}

// Stop halts the scheduled state checking
func (m *Manager) Stop() {
	log.Println("State lock manager stopping")
	m.scheduler.Stop()
}

// GetCurrentState returns the current state (thread-safe)
func (m *Manager) GetCurrentState() State {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.currentState
}

// checkState fetches the state from the configured endpoint
func (m *Manager) checkState() {
	ctx, cancel := context.WithTimeout(context.Background(), m.config.TimeoutDuration)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, m.config.StateEndpoint, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}

	resp, err := m.httpClient.Do(req)
	if err != nil {
		log.Printf("Error fetching state: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("State endpoint returned status: %d", resp.StatusCode)
		return
	}

	var stateResp StateResponse
	if err := json.NewDecoder(resp.Body).Decode(&stateResp); err != nil {
		log.Printf("Error decoding state response: %v", err)
		return
	}

	m.updateState(stateResp.State)
}

// updateState updates the current state and triggers callback if changed
func (m *Manager) updateState(newState State) {
	m.mu.Lock()
	oldState := m.currentState
	m.currentState = newState
	m.mu.Unlock()

	if oldState != newState {
		log.Printf("State changed: %s -> %s", oldState, newState)
		if m.config.OnStateChange != nil {
			go m.config.OnStateChange(oldState, newState)
		}
	}
}

// IsAllowed checks if a given HTTP method is allowed in the current state
func (m *Manager) IsAllowed(method string) bool {
	state := m.GetCurrentState()

	switch state {
	case StateOK:
		return true
	case StateSoftLock:
		return method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions
	case StateHardLock:
		return false
	default:
		return false
	}
}
