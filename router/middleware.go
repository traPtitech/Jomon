package router

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/logging"
	"github.com/traPtitech/Jomon/router/wrapsession"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	loginUserKey   = "login_user"
	sessionUserKey = "user"
)

func (h Handlers) setLoggerMiddleware(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
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
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		if err != nil {
			c.Error(err)
		}
		stop := time.Now()

		req := c.Request()
		res := c.Response()
		logger := logging.GetLogger(req.Context())
		latency := strconv.FormatFloat(stop.Sub(start).Seconds(), 'f', 9, 64) + "s"
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
			zap.String("latency", latency),
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

func (h Handlers) CheckLoginMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		logger := logging.GetLogger(ctx)

		id, err := wrapsession.WithSession(c, h.SessionName, func(w *wrapsession.W) (uuid.UUID, error) {
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
		user, err := h.Repository.GetUserByID(ctx, id)
		if err != nil {
			if ent.IsNotFound(err) {
				logger.Info("user not found in repository", zap.Error(err))
				return echo.NewHTTPError(http.StatusUnauthorized, "you are not logged in")
			}
			logger.Error("failed to get user from repository", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		c.Set(loginUserKey, userFromModelUser(*user))

		return next(c)
	}
}

func (h Handlers) CheckAdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		loginUser, _ := c.Get(loginUserKey).(User)
		if !loginUser.Admin {
			return echo.NewHTTPError(http.StatusForbidden, "you are not admin")
		}
		return next(c)
	}
}
