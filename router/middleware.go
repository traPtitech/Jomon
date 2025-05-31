package router

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	loginUserKey             = "login_user"
	sessionUserKey           = "user"
	sessionRequestCreatorKey = "request_creator"
	sessionFileCreatorKey    = "request_creator"
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
		sess, err := session.Get(h.SessionName, c)
		if err != nil {
			logger.Error("failed to get session", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		u, ok := sess.Values[sessionUserKey].(User)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "you are not logged in")
		}
		user, err := h.Repository.GetUserByID(ctx, u.ID)
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
		logger := logging.GetLogger(c.Request().Context())
		sess, err := session.Get(h.SessionName, c)
		if err != nil {
			logger.Error("failed to get session", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		user, err := getUserInfo(sess)
		if err != nil {
			logger.Error("failed to get user info", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		if !user.Admin {
			return echo.NewHTTPError(http.StatusForbidden, "you are not admin")
		}

		return next(c)
	}
}

func (h Handlers) CheckRequestCreatorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger := logging.GetLogger(c.Request().Context())
		sess, err := session.Get(h.SessionName, c)
		if err != nil {
			logger.Error("failed to get session", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		user, err := getUserInfo(sess)
		if err != nil {
			logger.Error("failed to get user info", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		creator, ok := sess.Values[sessionRequestCreatorKey].(uuid.UUID)
		if !ok {
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				"session request creator key is not set")
		}
		if creator != user.ID {
			return echo.NewHTTPError(http.StatusForbidden, "you are not creator")
		}

		return next(c)
	}
}

func (h Handlers) CheckAdminOrRequestCreatorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger := logging.GetLogger(c.Request().Context())
		sess, err := session.Get(h.SessionName, c)
		if err != nil {
			logger.Error("failed to get session", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		user, err := getUserInfo(sess)
		if err != nil {
			logger.Error("failed to get user info", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		creator, ok := sess.Values[sessionRequestCreatorKey].(uuid.UUID)
		if !ok {
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				"session request creator key is not set")
		}
		if creator != user.ID && !user.Admin {
			return echo.NewHTTPError(http.StatusForbidden, "you are not admin or creator")
		}

		return next(c)
	}
}

func (h Handlers) CheckFileCreatorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger := logging.GetLogger(c.Request().Context())
		sess, err := session.Get(h.SessionName, c)
		if err != nil {
			logger.Error("failed to get session", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		user, err := getUserInfo(sess)
		if err != nil {
			logger.Error("failed to get user info", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		creator, ok := sess.Values[sessionFileCreatorKey].(uuid.UUID)
		if !ok {
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				"session file creator key is not set")
		}
		if creator != user.ID {
			return echo.NewHTTPError(http.StatusForbidden, "you are not creator")
		}

		return next(c)
	}
}

func (h Handlers) CheckAdminOrFileCreatorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get(h.SessionName, c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		user, err := getUserInfo(sess)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		creator, ok := sess.Values[sessionFileCreatorKey].(uuid.UUID)
		if !ok {
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				"session file creator key is not set")
		}
		if creator != user.ID && !user.Admin {
			return echo.NewHTTPError(http.StatusForbidden, "you are not admin or creator")
		}

		return next(c)
	}
}

func (h Handlers) RetrieveRequestCreator() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger := logging.GetLogger(c.Request().Context())
			sess, err := session.Get(h.SessionName, c)
			if err != nil {
				logger.Error("failed to get session", zap.Error(err))
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			id, err := uuid.Parse(c.Param("requestID"))
			if err != nil {
				logger.Info("could not parse request_id as UUID", zap.Error(err))
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}

			ctx := c.Request().Context()
			request, err := h.Repository.GetRequest(ctx, id)
			if err != nil {
				logger.Error("failed to get request from repository", zap.Error(err))
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}

			sess.Values[sessionRequestCreatorKey] = request.CreatedBy

			if err = sess.Save(c.Request(), c.Response()); err != nil {
				logger.Error("failed to save session", zap.Error(err))
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}

			return next(c)
		}
	}
}

func (h Handlers) RetrieveFileCreator() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger := logging.GetLogger(c.Request().Context())
			sess, err := session.Get(h.SessionName, c)
			if err != nil {
				logger.Error("failed to get session", zap.Error(err))
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			id, err := uuid.Parse(c.Param("fileID"))
			if err != nil {
				logger.Info("could not parse file_id as UUID", zap.Error(err))
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}

			ctx := c.Request().Context()
			file, err := h.Repository.GetFile(ctx, id)
			if err != nil {
				logger.Error("failed to get file from repository", zap.Error(err))
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}

			sess.Values[sessionFileCreatorKey] = file.CreatedBy

			if err = sess.Save(c.Request(), c.Response()); err != nil {
				logger.Error("failed to save session", zap.Error(err))
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}

			return next(c)
		}
	}
}

func getUserInfo(sess *sessions.Session) (*User, error) {
	user, ok := sess.Values[sessionUserKey].(User)
	if !ok {
		return nil, errors.New("user not found")
	}

	return &user, nil
}
