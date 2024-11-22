package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type clientData struct {
	timestamps []time.Time
	mu         sync.Mutex
}

var (
	clientRateLimiters = make(map[int64]*clientData)
	rateLimiterMutex   sync.Mutex
)

func RateLimitMiddleware(maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetInt64("user_id")

		rateLimiterMutex.Lock()
		clientLimiter, exists := clientRateLimiters[userId]
		if !exists {
			clientLimiter = &clientData{timestamps: []time.Time{}}
			clientRateLimiters[userId] = clientLimiter
		}
		rateLimiterMutex.Unlock()

		clientLimiter.mu.Lock()
		defer clientLimiter.mu.Unlock()

		now := time.Now()
		threshold := now.Add(-window)
		validTimestamps := []time.Time{}
		for _, timestamp := range clientLimiter.timestamps {
			if timestamp.After(threshold) {
				validTimestamps = append(validTimestamps, timestamp)
			}
		}
		clientLimiter.timestamps = validTimestamps

		if len(clientLimiter.timestamps) >= maxRequests {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			return
		}

		clientLimiter.timestamps = append(clientLimiter.timestamps, now)

		rateLimiterMutex.Lock()
		clientRateLimiters[userId] = clientLimiter
		rateLimiterMutex.Unlock()

		c.Next()
	}
}
