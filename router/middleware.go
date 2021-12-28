package router

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/logging"
	"go.uber.org/zap"
)

func (h Handlers) AccessLoggingMiddleware(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			if err := next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()

			req := c.Request()
			res := c.Response()
			tmp := &logging.HTTPPayload{
				RequestMethod: req.Method,
				Status:        res.Status,
				UserAgent:     req.UserAgent(),
				RemoteIP:      c.RealIP(),
				Referer:       req.Referer(),
				Protocol:      req.Proto,
				RequestURL:    req.URL.String(),
				RequestSize:   req.Header.Get(echo.HeaderContentLength),
				ResponseSize:  strconv.FormatInt(res.Size, 10),
				Latency:       strconv.FormatFloat(stop.Sub(start).Seconds(), 'f', 9, 64) + "s",
			}
			httpCode := res.Status
			switch {
			case httpCode >= 500:
				errorRuntime, ok := c.Get("Error").(error)
				if ok {
					tmp.Error = errorRuntime.Error()
				} else {
					tmp.Error = "no data"
				}
				logger.Info("server error", zap.Object("field", tmp))
			case httpCode >= 400:
				errorRuntime, ok := c.Get("Error").(error)
				if ok {
					tmp.Error = errorRuntime.Error()
				} else {
					tmp.Error = "no data"
				}
				logger.Info("client error", zap.Object("field", tmp))
			case httpCode >= 300:
				logger.Info("redirect", zap.Object("field", tmp))
			case httpCode >= 200:
				logger.Info("success", zap.Object("field", tmp))
			}
			return nil
		}
	}
}

func (h Handlers) CheckLoginMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := h.SessionStore.Get(c.Request(), h.SessionName)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		_, ok := sess.Values[sessionUserKey].([]byte)
		if !ok {
			return c.Redirect(http.StatusSeeOther, "/api/auth/genpkce")
		}

		return next(c)
	}
}
