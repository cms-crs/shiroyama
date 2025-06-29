package utils

import (
	"net/http"

	"api-gateway/internal/models"
	"github.com/gin-gonic/gin"
)

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    data,
	})
}

func CreatedResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, models.Response{
		Success: true,
		Data:    data,
	})
}

func NoContentResponse(c *gin.Context) {
	c.JSON(http.StatusNoContent, models.Response{
		Success: true,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, models.Response{
		Success: false,
		Error: &models.ErrorInfo{
			Code:    http.StatusText(statusCode),
			Message: message,
		},
	})
}

func ErrorResponseWithCode(c *gin.Context, statusCode int, code, message string) {
	c.JSON(statusCode, models.Response{
		Success: false,
		Error: &models.ErrorInfo{
			Code:    code,
			Message: message,
		},
	})
}

func ErrorResponseWithDetails(c *gin.Context, statusCode int, code, message, details string) {
	c.JSON(statusCode, models.Response{
		Success: false,
		Error: &models.ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

func ValidationErrorResponse(c *gin.Context, err error) {
	ErrorResponseWithDetails(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", err.Error())
}

func ListResponse(c *gin.Context, data interface{}, pagination *models.PaginationResponse) {
	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data: models.ListResponse{
			Data:       data,
			Pagination: pagination,
		},
	})
}
