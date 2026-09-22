package middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"prangibar-go/config"

	"github.com/gin-gonic/gin"
)

type rateClient struct {
	count    int
	lastSeen time.Time
}

var (
	rateClients = make(map[string]*rateClient)
	rateMu      sync.Mutex
	rateInit    = false
)

func getRealIP(c *gin.Context) string {
	if xff := c.Request.Header.Get("X-Forwarded-For"); xff != "" {
		if idx := strings.Index(xff, ","); idx != -1 {
			xff = xff[:idx]
		}
		xff = strings.TrimSpace(xff)
		if xff != "" {
			return xff
		}
	}
	if xri := c.Request.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}
	return c.ClientIP()
}

func isPrivateIP(ip string) bool {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return true
	}
	privateRanges := []*net.IPNet{
		{IP: net.ParseIP("10.0.0.0"), Mask: net.CIDRMask(8, 32)},
		{IP: net.ParseIP("172.16.0.0"), Mask: net.CIDRMask(12, 32)},
		{IP: net.ParseIP("192.168.0.0"), Mask: net.CIDRMask(16, 32)},
		{IP: net.ParseIP("127.0.0.0"), Mask: net.CIDRMask(8, 32)},
	}
	for _, r := range privateRanges {
		if r.Contains(parsed) {
			return true
		}
	}
	return false
}

func ensureRateLimitInit() {
	if !rateInit {
		rateInit = true
		go rateCleanupRoutine()
	}
}

func rateCleanupRoutine() {
	for {
		time.Sleep(5 * time.Minute)
		window := time.Duration(config.App.Security.RateLimit.WindowSeconds) * time.Second
		rateMu.Lock()
		for ip, c := range rateClients {
			if time.Since(c.lastSeen) > window {
				delete(rateClients, ip)
			}
		}
		rateMu.Unlock()
	}
}

func RateLimit() gin.HandlerFunc {
	ensureRateLimitInit()

	return func(c *gin.Context) {
		ip := getRealIP(c)

		if isPrivateIP(ip) {
			c.Next()
			return
		}

		limit := config.App.Security.RateLimit.RequestsPerMinute
		window := time.Duration(config.App.Security.RateLimit.WindowSeconds) * time.Second

		rateMu.Lock()
		v, exists := rateClients[ip]
		if !exists {
			rateClients[ip] = &rateClient{count: 1, lastSeen: time.Now()}
			rateMu.Unlock()
			c.Next()
			return
		}

		if time.Since(v.lastSeen) > window {
			v.count = 1
			v.lastSeen = time.Now()
			rateMu.Unlock()
			c.Next()
			return
		}

		if v.count >= limit {
			rateMu.Unlock()
			c.Header("Retry-After", string(rune(config.App.Security.RateLimit.WindowSeconds)))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Terlalu banyak request, silakan coba lagi nanti",
			})
			c.Abort()
			return
		}

		v.count++
		v.lastSeen = time.Now()
		rateMu.Unlock()
		c.Next()
	}
}
