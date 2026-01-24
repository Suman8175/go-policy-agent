package license

import (
	"log"
	"sync"
)

var (
	scheduler     *Scheduler
	initOnce      sync.Once
	isInitialized bool
)

// Init initializes the license agent SDK with the given configuration
func Init(config Config) {
	initOnce.Do(func() {
		log.Printf("Initializing License Agent SDK with baseURL: %s, refreshInterval: %v",
			config.BaseURL, config.RefreshInterval)

		client := newClient(config.BaseURL)
		scheduler = newScheduler(client, config.RefreshInterval)
		scheduler.Start()

		isInitialized = true
		log.Println("License Agent SDK initialized successfully")
	})
}

// Shutdown gracefully stops the license refresh scheduler
func Shutdown() {
	if scheduler != nil {
		scheduler.Stop()
	}
}

// IsInitialized returns whether the SDK has been initialized
func IsInitialized() bool {
	return isInitialized
}
