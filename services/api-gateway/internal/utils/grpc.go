package utils

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

func GRPCErrorToHTTP(err error) (int, string) {
	if err == nil {
		return http.StatusOK, ""
	}

	st, ok := status.FromError(err)
	if !ok {
		return http.StatusInternalServerError, "Internal server error"
	}

	switch st.Code() {
	case codes.OK:
		return http.StatusOK, ""
	case codes.InvalidArgument:
		return http.StatusBadRequest, st.Message()
	case codes.NotFound:
		return http.StatusNotFound, "Resource not found"
	case codes.AlreadyExists:
		return http.StatusConflict, "Resource already exists"
	case codes.PermissionDenied:
		return http.StatusForbidden, "Permission denied"
	case codes.Unauthenticated:
		return http.StatusUnauthorized, "Authentication required"
	case codes.FailedPrecondition:
		return http.StatusPreconditionFailed, st.Message()
	case codes.Aborted:
		return http.StatusConflict, "Request aborted"
	case codes.OutOfRange:
		return http.StatusRequestedRangeNotSatisfiable, st.Message()
	case codes.Unimplemented:
		return http.StatusNotImplemented, "Not implemented"
	case codes.Internal:
		return http.StatusInternalServerError, "Internal server error"
	case codes.Unavailable:
		return http.StatusServiceUnavailable, "Service unavailable"
	case codes.DataLoss:
		return http.StatusInternalServerError, "Data loss"
	case codes.DeadlineExceeded:
		return http.StatusRequestTimeout, "Request timeout"
	case codes.Canceled:
		return http.StatusRequestTimeout, "Request canceled"
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests, "Rate limit exceeded"
	default:
		return http.StatusInternalServerError, "Internal server error"
	}
}

func HandleGRPCError(c *gin.Context, err error, operation string) {
	if err == nil {
		return
	}

	//statusCode, message := GRPCErrorToHTTP(err)

	st, _ := status.FromError(err)

	switch st.Code() {
	case codes.NotFound, codes.InvalidArgument, codes.AlreadyExists:
	case codes.Internal, codes.Unavailable, codes.DataLoss:
	default:
	}

	//ErrorResponse(c, statusCode, message)
}

func WithRetry(ctx context.Context, operation func() error, maxRetries int, delay time.Duration) error {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			waitTime := delay * time.Duration(1<<uint(attempt-1))
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(waitTime):
			}
		}

		err := operation()
		if err == nil {
			return nil
		}

		lastErr = err

		if !shouldRetry(err) {
			return err
		}
	}

	return fmt.Errorf("operation failed after %d attempts: %w", maxRetries, lastErr)
}

func shouldRetry(err error) bool {
	st, ok := status.FromError(err)
	if !ok {
		return false
	}

	switch st.Code() {
	case codes.Unavailable, codes.DeadlineExceeded, codes.Internal, codes.ResourceExhausted:
		return true
	default:
		return false
	}
}

func CreateContextWithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout <= 0 {
		timeout = 30 * time.Second // Значение по умолчанию
	}
	return context.WithTimeout(parent, timeout)
}

func IsRetryableError(err error) bool {
	return shouldRetry(err)
}

func GetGRPCTimeout(operationType string) time.Duration {
	switch operationType {
	case "create", "update", "delete":
		return 30 * time.Second
	case "get", "list":
		return 15 * time.Second
	case "search":
		return 45 * time.Second
	case "upload":
		return 2 * time.Minute
	default:
		return 30 * time.Second
	}
}

type GRPCHealthStatus struct {
	Service   string    `json:"service"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Error     string    `json:"error,omitempty"`
}

func CheckGRPCHealth(ctx context.Context, serviceName string, healthCheck func(ctx context.Context) error) GRPCHealthStatus {
	status := GRPCHealthStatus{
		Service:   serviceName,
		Timestamp: time.Now(),
	}

	healthCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := healthCheck(healthCtx); err != nil {
		status.Status = "unhealthy"
		status.Error = err.Error()
	} else {
		status.Status = "healthy"
	}

	return status
}

type GRPCMetrics struct {
	Operation string        `json:"operation"`
	Duration  time.Duration `json:"duration"`
	Success   bool          `json:"success"`
	Error     string        `json:"error,omitempty"`
}

func MeasureGRPCOperation(operation string, fn func() error) GRPCMetrics {
	start := time.Now()
	err := fn()
	duration := time.Since(start)

	metrics := GRPCMetrics{
		Operation: operation,
		Duration:  duration,
		Success:   err == nil,
	}

	if err != nil {
		metrics.Error = err.Error()
	}

	return metrics
}
