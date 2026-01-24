package license

import (
	"log"
	"time"
)

type Scheduler struct {
	client          *Client
	refreshInterval time.Duration
	stopChan        chan struct{}
}

func newScheduler(client *Client, refreshInterval time.Duration) *Scheduler {
	return &Scheduler{
		client:          client,
		refreshInterval: refreshInterval,
		stopChan:        make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	// Initial fetch
	s.refresh()

	// Start periodic refresh
	go func() {
		ticker := time.NewTicker(s.refreshInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.refresh()
			case <-s.stopChan:
				log.Println("License scheduler stopped")
				return
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	close(s.stopChan)
}

func (s *Scheduler) refresh() {
	log.Println("License refresh started")

	resp, err := s.client.FetchLicenseStatus()
	if err != nil {
		log.Printf("Failed to refresh license status: %v", err)
		// On error, set to HARD_LOCK for fail-secure behavior
		setStatus(HARD_LOCK)
		return
	}

	setStatus(resp.LicenseStatus)
	log.Printf("License status updated to: %s", resp.LicenseStatus)
}
