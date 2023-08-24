package httpgin

import (
	"dynamic-user-segmentation/pkg/logging"
	"github.com/gin-gonic/gin"
	"time"
)

func LoggerMiddleware(log logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		latency := time.Since(t)
		status := c.Writer.Status()
		log.Info("Latency:", latency, "\tMethod:", c.Request.Method, "\tPath:",
			c.Request.URL.Path, "\tStatus:", status)
	}
}
