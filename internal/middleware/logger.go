package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger(logger *slog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		status := ctx.Writer.Status()

		attributes := []any{
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.Request.URL.Path),
			slog.String("route", ctx.FullPath()),
			slog.Int("status", status),
			slog.Float64("duration_ms", float64(time.Since(start).Nanoseconds())/1e6),
		}

		if errs := ctx.Errors.Errors(); len(errs) > 0 {
			attributes = append(attributes, slog.Any("errors", errs))
		}

		req := ctx.Request.Context()

		switch {
		case status >= http.StatusInternalServerError:
			logger.ErrorContext(req, "request", attributes...)
		case status >= http.StatusBadRequest:
			logger.WarnContext(req, "request", attributes...)
		default:
			logger.InfoContext(req, "request", attributes...)
		}

	}
}
