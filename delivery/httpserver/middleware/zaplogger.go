package middleware

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"mymodule/logger"
	"net/http"
	"time"
)

func ZapLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		req := c.Request()

		// 1. Execute the handler FIRST to achieve handler status
		err := next(c) // ← This is where status code gets set
		// 2. NOW log the response details
		res := c.Response()
		finalStatus := res.Status
		if err != nil {
			// If Echo sets an error response, use the correct status code

			var httpErr *echo.HTTPError
			if errors.As(err, &httpErr) {
				finalStatus = httpErr.Code
			}
		}
		fields := []zap.Field{
			zap.String("remote_ip", c.RealIP()),
			zap.String("latency", time.Since(start).String()),
			zap.String("host", req.Host),
			zap.String("uri", req.RequestURI),
			zap.String("method", req.Method),
			zap.Int("status", finalStatus), // ← Now gets the real status
			zap.String("user_agent", req.UserAgent()),
		}
		fmt.Println("ssss", res.Status)
		switch {
		case finalStatus >= http.StatusInternalServerError:
			logger.Info("server error", fields...)
		case finalStatus >= http.StatusBadRequest:
			logger.Warn("client error", fields...)
		default:
			logger.Info("ok", fields...)
		}

		return err // Return the handler's error (if any)
	}
}
