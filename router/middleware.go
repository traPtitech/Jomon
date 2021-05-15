package router

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/ent"
	log "github.com/traPtitech/Jomon/logging"
	"go.uber.org/zap"
)

const (
	contextUserKey        = "user"
	contextAccessTokenKey = "access_token"
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
			tmp := &log.HTTPPayload{
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

func (h Handlers) AuthUser(c echo.Context) (echo.Context, error) {
	sess, err := session.Get(sessionKey, c)
	if err != nil {
		return nil, c.NoContent(http.StatusInternalServerError)
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   sessionDuration,
		HttpOnly: true,
	}

	accTok, ok := sess.Values[sessionAccessTokenKey].(string)
	if !ok || accTok == "" {
		return nil, c.NoContent(http.StatusUnauthorized)
	}
	c.Set(contextAccessTokenKey, accTok)

	user, ok := sess.Values[sessionUserKey].(ent.User)
	if !ok {
		user, err = h.Service.Users.GetMyUser(accTok)
		sess.Values[sessionUserKey] = user
		if err := sess.Save(c.Request(), c.Response()); err != nil {
			return nil, c.NoContent(http.StatusInternalServerError)
		}

		if err != nil {
			return nil, c.NoContent(http.StatusInternalServerError)
		}
	}

	admins, err := h.Service.Administrators.GetAdministratorList()
	if err != nil {
		return nil, c.NoContent(http.StatusInternalServerError)
	}
	user.GiveIsUserAdmin(admins)

	c.Set(contextUserKey, user)

	return c, nil
}
