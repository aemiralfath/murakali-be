package middleware

import (
	"github.com/gin-gonic/gin"
	"time"
)

func (mw *MWManager) Latency() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		latency := time.Since(t).Milliseconds()
		if latency > 400 {
			mw.logger.Warnf("request latency: %dms", latency)
		} else {
			mw.logger.Infof("request latency: %dms", latency)
		}
	}
}
