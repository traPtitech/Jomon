package router

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/traPtitech/Jomon/internal/logging"
	"github.com/traPtitech/Jomon/internal/router/wrapsession"
	"github.com/traPtitech/Jomon/internal/service"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	loginUserKey = "login_user"
)

func (h Handlers) setLoggerMiddleware(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			req := c.Request()
			ctx := req.Context()
			reqID := req.Header.Get(echo.HeaderXRequestID)
			l := logger.With(zap.String("requestID", reqID))
			ctx = logging.SetLogger(ctx, l)
			c.SetRequest(req.WithContext(ctx))
			return next(c)
		}
	}
}

// AccessLoggingMiddleware ですべてのエラーを出力する
func (h Handlers) AccessLoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	// TODO(logging): https://echo.labstack.com/docs/middleware/logger を使う
	return func(c *echo.Context) error {
		start := time.Now()
		err := next(c)
		if err != nil {
			defaultHTTPErrorHandler(c, HTTPErrorHandlerInner(err))
		}
		stop := time.Now()

		req := c.Request()
		latency := strconv.FormatFloat(stop.Sub(start).Seconds(), 'f', 9, 64) + "s"
		fields := []zapcore.Field{
			zap.String("requestMethod", req.Method),
			zap.String("userAgent", req.UserAgent()),
			zap.String("remoteIp", c.RealIP()),
			zap.String("referer", req.Referer()),
			zap.String("protocol", req.Proto),
			zap.String("requestUrl", req.URL.String()),
			zap.String("requestSize", req.Header.Get(echo.HeaderContentLength)),
			zap.String("latency", latency),
		}
		rw, uErr := echo.UnwrapResponse(c.Response())
		if uErr != nil {
			logger := logging.GetLogger(req.Context())
			fields = append(
				fields,
				zap.Error(uErr),
				zap.String("status", "unknown"),
				zap.String("responseSize", "unknown"))
			logger.Error("failed to unwrap response for access logging", fields...)
			return nil
		}
		fields = append(fields,
			zap.Int("status", rw.Status),
			zap.String("responseSize", strconv.FormatInt(rw.Size, 10)))
		logger := logging.GetLogger(req.Context())
		httpCode := rw.Status
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

func (h Handlers) CheckLoginMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		ctx := c.Request().Context()

		id, err := wrapsession.WithSession(
			c, h.SessionName, func(w *wrapsession.W) (uuid.UUID, error) {
				v, ok := w.GetUserID()
				if !ok {
					err := echo.NewHTTPError(http.StatusUnauthorized, "you are not logged in")
					return uuid.Nil, err
				}
				return v, nil
			})
		if err != nil {
			return err
		}
		accessUser, err := h.Service.GetUserByID(ctx, id)
		if err != nil {
			return err
		}
		c.Set(loginUserKey, accessUser)

		return next(c)
	}
}

func (h Handlers) CheckAccountManagerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		loginUser, _ := c.Get(loginUserKey).(*service.User)
		if !loginUser.AccountManager {
			return echo.NewHTTPError(http.StatusForbidden, "you are not accountManager")
		}
		return next(c)
	}
}
