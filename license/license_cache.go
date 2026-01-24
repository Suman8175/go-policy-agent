package license

import "sync"

var (
	currentStatus LicenseStatus = HARD_LOCK
	cacheMutex    sync.RWMutex
)

// Get returns the current license status (thread-safe read)
func GetStatus() LicenseStatus {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()
	return currentStatus
}

// setStatus updates the license status (thread-safe write)
func setStatus(status LicenseStatus) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	currentStatus = status
}
