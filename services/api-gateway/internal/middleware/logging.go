package middleware

import (
	"time"

	"api-gateway/pkg/logger"
	"github.com/gin-gonic/gin"
)

func Logging(log logger.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		log.Info("HTTP Request",
			"timestamp", param.TimeStamp.Format(time.RFC3339),
			"status", param.StatusCode,
			"latency", param.Latency,
			"client_ip", param.ClientIP,
			"method", param.Method,
			"path", param.Path,
			"error", param.ErrorMessage,
		)
		return ""
	})
}
