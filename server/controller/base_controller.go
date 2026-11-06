package controller

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/your-org/go-web-api-framework/pkg/validator"
	"github.com/your-org/go-web-api-framework/types"
)

// BaseController provides common HTTP handling utilities.
type BaseController struct{}

// Success writes a 200 JSON response with the supplied data.
func (bc *BaseController) Success(c *gin.Context, data interface{}) {
	types.SuccessResponse(c, data)
}

// SuccessWithMessage writes a 200 JSON response with a custom message and data.
func (bc *BaseController) SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	types.SuccessResponseWithMessage(c, message, data)
}

// Created writes a 201 JSON response with Success: true and the supplied data.
func (bc *BaseController) Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, types.Response[interface{}]{
		Success: true,
		Data:    data,
	})
}

// BadRequest writes a 400 JSON error response.
func (bc *BaseController) BadRequest(c *gin.Context, message string, err error) {
	appErr := types.NewAppError(types.ErrCodeBadRequest, message, err)
	types.ErrorResponse(c, appErr)
}

// NotFound writes a 404 JSON error response.
func (bc *BaseController) NotFound(c *gin.Context, message string) {
	appErr := types.NewAppError(types.ErrCodeNotFound, message, nil)
	types.ErrorResponse(c, appErr)
}

// Unauthorized writes a 401 JSON error response.
func (bc *BaseController) Unauthorized(c *gin.Context, message string) {
	appErr := types.NewAppError(types.ErrCodeUnauthorized, message, nil)
	types.ErrorResponse(c, appErr)
}

// Forbidden writes a 403 JSON error response.
func (bc *BaseController) Forbidden(c *gin.Context, message string) {
	appErr := types.NewAppError(types.ErrCodeForbidden, message, nil)
	types.ErrorResponse(c, appErr)
}

// InternalError writes a 500 JSON error response.
func (bc *BaseController) InternalError(c *gin.Context, message string, err error) {
	appErr := types.NewAppError(types.ErrCodeInternalServer, message, err)
	types.ErrorResponse(c, appErr)
}

// ParamInt extracts an integer path parameter.
// Returns 0 if the value is not a valid integer.
func (bc *BaseController) ParamInt(c *gin.Context, key string) int {
	val, _ := strconv.Atoi(c.Param(key))
	return val
}

// QueryInt extracts an integer query parameter.
// Returns defaultValue if the parameter is missing or not a valid integer.
func (bc *BaseController) QueryInt(c *gin.Context, key string, defaultValue int) int {
	val := c.Query(key)
	if val == "" {
		return defaultValue
	}
	result, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}
	return result
}

// BindJSON decodes a JSON request body into the supplied object.
func (bc *BaseController) BindJSON(c *gin.Context, obj interface{}) error {
	return c.ShouldBindJSON(obj)
}

// BindQuery decodes query-string parameters into the supplied object.
func (bc *BaseController) BindQuery(c *gin.Context, obj interface{}) error {
	return c.ShouldBindQuery(obj)
}

// ValidateAndBind binds the request body and validates the resulting struct.
func (bc *BaseController) ValidateAndBind(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBind(obj); err != nil {
		return err
	}
	v := validator.New()
	return v.ValidateStruct(obj)
}

// Paginate writes a paginated JSON response.
func (bc *BaseController) Paginate(c *gin.Context, items interface{}, total int64, page, perPage int) {
	// Convert the items slice to []any so it can be passed to PaginatedResponse[T=any].
	v := reflect.ValueOf(items)
	n := v.Len()
	anySlice := make([]any, n)
	for i := 0; i < n; i++ {
		anySlice[i] = v.Index(i).Interface()
	}
	types.PaginatedResponse(c, anySlice, total, page, perPage)
}

// NoContent writes a 204 No Content response.
func (bc *BaseController) NoContent(c *gin.Context) {
	c.AbortWithStatus(http.StatusNoContent)
}
