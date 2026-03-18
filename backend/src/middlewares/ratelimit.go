package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type IPRateLimit struct {
	Limiters map[string]*ipEntry
	Mu sync.Mutex
	Rps rate.Limit
	Burst int
}

type ipEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func NewIPRateLimit(rps rate.Limit, burst int) *IPRateLimit {
	rl := &IPRateLimit{
		Limiters: make(map[string]*ipEntry),
		Rps: rps,
		Burst: burst,
	}
	go rl.cleanuploop()
	return rl
}

func (i *IPRateLimit) GetLimiter(ip string) *rate.Limiter {
	i.Mu.Lock()
	defer i.Mu.Unlock()

	entry, exists := i.Limiters[ip]
	if !exists {
		entry= &ipEntry{limiter: rate.NewLimiter(i.Rps, i.Burst)}
		i.Limiters[ip] = entry
	}

	entry.lastSeen = time.Now()
	return entry.limiter
}

func (i *IPRateLimit) cleanuploop() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		i.Mu.Lock() 
		for ip, entry := range i.Limiters {
			if time.Since(entry.lastSeen) > 5*time.Minute {
				delete(i.Limiters, ip)
			}
		}
		i.Mu.Unlock()
	}
}

func (i *IPRateLimit) RateLimiting() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := i.GetLimiter(ip)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}

		c.Next()
	}
}
