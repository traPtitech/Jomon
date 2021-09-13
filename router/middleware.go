package router

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/traPtitech/Jomon/logging"
)

// AccessLoggingMiddleware ですべてのエラーを出力する
func (h Handlers) AccessLoggingMiddleware(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			elapsed := time.Since(start)

			req := c.Request()
			res := c.Response()
			payload := &logging.HTTPPayload{
				RequestMethod: req.Method,
				Status:        res.Status,
				UserAgent:     req.UserAgent(),
				RemoteIP:      c.RealIP(),
				Referer:       req.Referer(),
				Protocol:      req.Proto,
				RequestURL:    req.URL.String(),
				RequestSize:   req.Header.Get(echo.HeaderContentLength),
				ResponseSize:  strconv.FormatInt(res.Size, 10),
				Latency:       strconv.FormatFloat(elapsed.Seconds(), 'f', 9, 64) + "s",
				Error:         err.Error(),
			}
			switch {
			case res.Status >= 500:
				logger.Error("server error", zap.Object("field", payload))
			case res.Status >= 400:
				logger.Info("client error", zap.Object("field", payload))
			case res.Status >= 300:
				logger.Info("redirect", zap.Object("field", payload))
			case res.Status >= 200:
				logger.Info("success", zap.Object("field", payload))
			default:
				logger.Error("unknown", zap.Object("field", payload))
			}
			return nil
		}
	}
}

func (h Handlers) CheckLoginMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get(h.SessionName, c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		_, ok := sess.Values[sessionUserKey].(User)
		if !ok {
			return c.Redirect(http.StatusUnauthorized, "/api/auth/genpkce")
		}

		return next(c)
	}
}
