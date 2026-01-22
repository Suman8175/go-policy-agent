package statelock

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GinMiddleware returns a Gin middleware that enforces state-based access control
func (m *Manager) GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !m.IsAllowed(c.Request.Method) {
			state := m.GetCurrentState()

			var message string
			var statusCode int

			switch state {
			case StateSoftLock:
				message = "Service is in SOFT_LOCK mode. Only read operations are allowed."
				statusCode = http.StatusForbidden
			case StateHardLock:
				message = "Service is in HARD_LOCK mode. All operations are blocked."
				statusCode = http.StatusServiceUnavailable
			default:
				message = "Operation not allowed in current state."
				statusCode = http.StatusForbidden
			}

			c.JSON(statusCode, gin.H{
				"error": message,
				"state": state,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// HealthCheckMiddleware provides a health check endpoint that bypasses state checking
func (m *Manager) HealthCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"state":  m.GetCurrentState(),
		})
	}
}
