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
	"github.com/traPtitech/Jomon/model"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	sessionUserKey           = "user"
	sessionOwnerKey          = "group_owner"
	sessionRequestCreatorKey = "request_creator"
	sessionFileCreatorKey    = "request_creator"
)

// AccessLoggingMiddleware ですべてのエラーを出力する
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
		sess, err := session.Get(h.SessionName, c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		_, ok := sess.Values[sessionUserKey].(User)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "you are not logged in")
		}

		return next(c)
	}
}

func (h Handlers) CheckAdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get(h.SessionName, c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		user, err := getUserInfo(sess)
		if err != nil {
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
		sess, err := session.Get(h.SessionName, c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		user, err := getUserInfo(sess)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		creator := sess.Values[sessionRequestCreatorKey].(uuid.UUID)
		if creator != user.ID {
			return echo.NewHTTPError(http.StatusForbidden, "you are not creator")
		}

		return next(c)
	}
}

func (h Handlers) CheckAdminOrRequestCreatorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get(h.SessionName, c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		user, err := getUserInfo(sess)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		creator := sess.Values[sessionRequestCreatorKey].(uuid.UUID)
		if creator != user.ID && !user.Admin {
			return echo.NewHTTPError(http.StatusForbidden, "you are not admin or creator")
		}

		return next(c)
	}
}

func (h Handlers) CheckAdminOrGroupOwnerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get(h.SessionName, c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		user, err := getUserInfo(sess)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		owners, ok := sess.Values[sessionOwnerKey].([]*model.Owner)
		if !ok {
			return echo.ErrInternalServerError
		}

		for _, owner := range owners {
			if owner.ID == user.ID {
				if user.Admin {
					return next(c)
				}
			}
		}

		return echo.ErrForbidden
	}
}

func (h Handlers) CheckFileCreatorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get(h.SessionName, c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		user, err := getUserInfo(sess)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		creator := sess.Values[sessionFileCreatorKey].(uuid.UUID)
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

		creator := sess.Values[sessionFileCreatorKey].(uuid.UUID)
		if creator != user.ID && !user.Admin {
			return echo.NewHTTPError(http.StatusForbidden, "you are not admin or creator")
		}

		return next(c)
	}
}

func (h Handlers) RetrieveGroupOwner(repo model.Repository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, err := session.Get(h.SessionName, c)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			id, err := uuid.Parse(c.Param("groupID"))
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}

			ctx := c.Request().Context()
			owners, err := repo.GetOwners(ctx, id)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}

			sess.Values[sessionOwnerKey] = owners

			if err = sess.Save(c.Request(), c.Response()); err != nil {
				return echo.ErrInternalServerError
			}

			return next(c)
		}
	}
}

func (h Handlers) RetrieveRequestCreator(repo model.Repository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, err := session.Get(h.SessionName, c)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			id, err := uuid.Parse(c.Param("requestID"))
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}

			ctx := c.Request().Context()
			request, err := repo.GetRequest(ctx, id)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}

			sess.Values[sessionRequestCreatorKey] = request.CreatedBy

			if err = sess.Save(c.Request(), c.Response()); err != nil {
				return echo.ErrInternalServerError
			}

			return next(c)
		}
	}
}

func (h Handlers) RetrieveFileCreator(repo model.Repository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, err := session.Get(h.SessionName, c)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			id, err := uuid.Parse(c.Param("fileID"))
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}

			ctx := c.Request().Context()
			file, err := repo.GetFile(ctx, id)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}

			sess.Values[sessionFileCreatorKey] = file.CreatedBy

			if err = sess.Save(c.Request(), c.Response()); err != nil {
				return echo.ErrInternalServerError
			}

			return next(c)
		}
	}
}

func getUserInfo(sess *sessions.Session) (*User, error) {
	user, ok := sess.Values[sessionUserKey].(*User)
	if !ok {
		return nil, errors.New("user not found")
	}

	return user, nil
}
