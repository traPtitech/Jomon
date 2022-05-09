package router

import (
	"encoding/gob"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (h Handlers) AccessLoggingMiddleware(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			if err != nil {
				c.Error(err)
			}
			stop := time.Now()

			req := c.Request()
			res := c.Response()
			fields := []zapcore.Field{
				zap.String("requestMethod", req.Method),
				zap.Int("status", res.Status),
				zap.String("userAgent", req.UserAgent()),
				zap.String("remoteIp", c.RealIP()),
				zap.String("referer", req.Referer()),
				zap.String("protocol", req.Proto),
				zap.String("requestUrl", req.URL.String()),
				zap.String("requestSize", req.Header.Get(echo.HeaderContentLength)),
				zap.String("responseSize", strconv.FormatInt(res.Size, 10)),
				zap.String("latency", strconv.FormatFloat(stop.Sub(start).Seconds(), 'f', 9, 64)+"s"),
			}
			httpCode := res.Status
			switch {
			case httpCode >= 500:
				fields = append(fields, zap.Error(err))
				logger.Error("server error", fields...)
			case httpCode >= 400:
				fields = append(fields, zap.Error(err))
				logger.Warn("client error", fields...)
			case httpCode >= 300:
				logger.Info("redirect", fields...)
			default:
				logger.Info("success", fields...)
			}
			return nil
		}
	}
}

func (h Handlers) CheckLoginMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		gob.Register(&User{})
		sess, err := h.SessionStore.Get(c.Request(), h.SessionName)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		_, ok := sess.Values[sessionUserKey].(*User)
		if !ok {
			return c.Redirect(http.StatusSeeOther, "/api/auth/genpkce")
		}

		return next(c)
	}
}
