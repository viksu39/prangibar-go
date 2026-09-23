package middleware

import (
	"time"

	"prangibar-go/config"
	"prangibar-go/models"

	"github.com/gin-gonic/gin"
)

func APILogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		statusCode := c.Writer.Status()
		responseTime := time.Since(start).Milliseconds()

		ip := getRealIP(c)
		userAgent := c.Request.UserAgent()
		adminID, _ := c.Get("adminID")
		var adminIDPtr *int32
		if id, ok := adminID.(int32); ok {
			adminIDPtr = &id
		}

		log := models.ApiLog{
			Method:         c.Request.Method,
			Path:           c.Request.URL.Path,
			StatusCode:     statusCode,
			IP:             &ip,
			UserAgent:      &userAgent,
			AdminID:        adminIDPtr,
			ResponseTimeMs: int(responseTime),
		}

		go func() {
			config.DB.Create(&log)
		}()
	}
}
