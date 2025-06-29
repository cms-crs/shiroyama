package middleware

import (
	"net/http"

	"api-gateway/internal/utils"
	"api-gateway/pkg/logger"
	"github.com/gin-gonic/gin"
)

func Recovery(log logger.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Error("Panic recovered",
			"panic", recovered,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"ip", c.ClientIP(),
		)

		utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
	})
}
