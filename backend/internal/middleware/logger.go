package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Logger provides request logging
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Generate request ID if not present
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)

		// Process request
		c.Next()

		// Log after request
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path

		log.Printf("[%s] %s %s %d %v %s",
			requestID[:8],
			method,
			path,
			statusCode,
			latency,
			clientIP,
		)
	}
}

// RequestID ensures each request has a unique ID
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()
	}
}

// GetRequestID helper to get request ID from context
func GetRequestID(c *gin.Context) string {
	requestID, _ := c.Get("request_id")
	if requestID == nil {
		return ""
	}
	return requestID.(string)
}

// GetClientIP helper to get client IP
func GetClientIP(c *gin.Context) string {
	// Check X-Forwarded-For header first (for proxied requests)
	xff := c.GetHeader("X-Forwarded-For")
	if xff != "" {
		return xff
	}

	// Check X-Real-IP header
	xri := c.GetHeader("X-Real-IP")
	if xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	return c.ClientIP()
}
