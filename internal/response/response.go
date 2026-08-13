package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Data    any      `json:"data,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

func send(ctx *gin.Context, statusCode int, success bool, message string, data any, errors []string) {
	ctx.JSON(statusCode, Response{
		Success: success,
		Message: message,
		Data:    data,
		Errors:  errors,
	})
}

func firstOr(fallback string, values ...string) string {
	if len(values) > 0 && values[0] != "" {
		return values[0]
	}
	return fallback
}

func Ok(ctx *gin.Context, data any, message ...string) {
	send(ctx, http.StatusOK, true, firstOr("operation completed successfully!", message...), data, nil)
}

func Created(ctx *gin.Context, data any, message ...string) {
	send(ctx, http.StatusCreated, true, firstOr("resource created successfully!", message...), data, nil)
}

func BadRequest(ctx *gin.Context, message string, errors ...string) {
	send(ctx, http.StatusBadRequest, false, message, nil, errors)
}

func Unauthorized(ctx *gin.Context, message ...string) {
	send(ctx, http.StatusUnauthorized, false, firstOr("unauthorized.", message...), nil, nil)
}

func NotFound(ctx *gin.Context, message ...string) {
	send(ctx, http.StatusNotFound, false, firstOr("resource not found.", message...), nil, nil)
}

func Conflict(ctx *gin.Context, message string) {
	send(ctx, http.StatusConflict, false, message, nil, nil)
}

func InternalServerError(ctx *gin.Context, message ...string) {
	send(ctx, http.StatusInternalServerError, false, firstOr("internal server error.", message...), nil, nil)
}
