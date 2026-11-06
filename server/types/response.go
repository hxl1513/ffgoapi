package types

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response is a generic envelope returned by successful or failed API calls.
type Response[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

// PaginatedResult is a generic envelope for paginated list endpoints.
type PaginatedResult[T any] struct {
	Success bool       `json:"success"`
	Data    []T       `json:"data"`
	Total   int64     `json:"total"`
	Page    int       `json:"page"`
	PerPage int       `json:"per_page"`
}

// SuccessResponse writes a 200 JSON response with the supplied data.
func SuccessResponse[T any](c *gin.Context, data T) {
	c.JSON(http.StatusOK, Response[T]{
		Success: true,
		Data:    data,
	})
}

// SuccessResponseWithMessage writes a 200 JSON response with a custom message and data.
func SuccessResponseWithMessage[T any](c *gin.Context, message string, data T) {
	c.JSON(http.StatusOK, Response[T]{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse writes a JSON error response using the HTTP status mapped from
// the AppError code.
func ErrorResponse(c *gin.Context, err *AppError) {
	status := httpStatus(err.Code)
	c.JSON(status, Response[any]{
		Success: false,
		Message: err.Message,
	})
}

// PaginatedResponse writes a 200 JSON response for paginated list endpoints.
func PaginatedResponse[T any](c *gin.Context, data []T, total int64, page, perPage int) {
	c.JSON(http.StatusOK, PaginatedResult[T]{
		Success: true,
		Data:    data,
		Total:   total,
		Page:    page,
		PerPage: perPage,
	})
}

// httpStatus maps an application error code to an HTTP status code.
func httpStatus(code int) int {
	switch code {
	case ErrCodeBadRequest:
		return http.StatusBadRequest
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeValidation:
		return http.StatusUnprocessableEntity
	case ErrCodeInternalServer:
		return http.StatusInternalServerError
	case ErrCodeDatabase:
		return http.StatusInternalServerError
	case ErrCodeCache:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
