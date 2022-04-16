package router

import (
	"context"
	"encoding/gob"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/logging"
	"github.com/traPtitech/Jomon/model"
	"go.uber.org/zap"
)

const (
	sessionUserKey    = "user"
	sessionOwnerKey   = "group_owner"
	sessionCreatorKey = "request_creator"
)

// AccessLoggingMiddleware ですべてのエラーを出力する
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
		sess, err := session.Get(h.SessionName, c)
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

		creator := sess.Values[sessionCreatorKey].(uuid.UUID)
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

		creator := sess.Values[sessionCreatorKey].(uuid.UUID)
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

			ctx := context.Background()
			owners, err := repo.GetOwners(ctx, id)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}

			gob.Register([]*model.Owner{})

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

			ctx := context.Background()
			request, err := repo.GetRequest(ctx, id)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}

			gob.Register(&uuid.UUID{})

			sess.Values[sessionCreatorKey] = request.CreatedBy

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
