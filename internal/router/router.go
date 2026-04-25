package router

import (
	"encoding/gob"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/v5/session"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"go.uber.org/zap"

	"github.com/traPtitech/Jomon/internal/logging"
	"github.com/traPtitech/Jomon/internal/service"
	"github.com/traPtitech/Jomon/internal/webhook"
)

type Handlers struct {
	WebhookService *webhook.Service
	Service        *service.Service
	SessionName    string
}

var defaultHTTPErrorHandler = echo.DefaultHTTPErrorHandler(false)

func (h Handlers) NewServer(logger *zap.Logger) *echo.Echo {
	e := echo.New()
	// TODO(logging): e.Loggerに与えられたloggerを適用する
	e.HTTPErrorHandler = func(c *echo.Context, err error) {
		logger := logging.GetLogger(c.Request().Context())
		logger.Debug("handling error", zap.Error(err))
		he := HTTPErrorHandlerInner(err)
		defaultHTTPErrorHandler(c, he)
	}
	e.Use(middleware.RequestID())
	e.Use(h.setLoggerMiddleware(logger))
	e.Use(h.AccessLoggingMiddleware)
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))))
	gob.Register(User{})
	gob.Register(uuid.UUID{})

	api := e.Group("/api")
	{
		apiAuth := api.Group("/auth")
		{
			apiAuth.GET("/callback", h.AuthCallback)
			apiAuth.GET("/signin", h.BeginAuth)
		}

		apiApplications := api.Group("/applications", h.CheckLoginMiddleware)
		{
			apiApplications.GET("", h.GetApplications)
			apiApplications.POST(
				"",
				h.PostApplication,
				middleware.BodyDump(h.WebhookService.WebhookApplicationsEventHandler))
			apiApplicationIDs := apiApplications.Group("/:applicationID")
			{
				apiApplicationIDs.GET("", h.GetApplication)
				// FIXME: このままでは異常系のApplicationでもwebhookが呼ばれる
				// そのため、webhookの関数呼び出しをPutApplication内に移す
				apiApplicationIDs.PUT(
					"",
					h.PutApplication,
					middleware.BodyDump(h.WebhookService.WebhookApplicationsEventHandler))
				apiApplicationIDs.POST(
					"/comments",
					h.PostComment,
					middleware.BodyDump(h.WebhookService.WebhookApplicationsEventHandler))
				apiApplicationIDs.PUT("/status", h.PutStatus)
			}
		}

		apiFiles := api.Group("/files", h.CheckLoginMiddleware)
		{
			apiFiles.POST("", h.PostFile)
			apiFileIDs := apiFiles.Group("/:fileID")
			{
				apiFileIDs.GET("", h.GetFile)
				apiFileIDs.DELETE("", h.DeleteFile)
				apiFileIDs.GET("/meta", h.GetFileMeta)
			}
		}

		apiTags := api.Group("/tags", h.CheckLoginMiddleware)
		{
			apiTags.GET("", h.GetTags)
			apiTags.POST("", h.PostTag)
			apiTags.PUT("/:tagID", h.PutTag)
			apiTags.DELETE("/:tagID", h.DeleteTag)
		}

		apiUsers := api.Group("/users", h.CheckLoginMiddleware)
		{
			apiUsers.GET("", h.GetUsers)
			apiUsers.PUT("", h.UpdateUserInfo)
			apiUsers.GET("/me", h.GetMe)
		}

		apiAccountManagers := api.Group("/account-managers", h.CheckLoginMiddleware)
		{
			apiAccountManagers.GET("", h.GetAccountManagers)
			apiAccountManagers.POST("", h.PostAccountManagers)
			apiAccountManagers.DELETE("", h.DeleteAccountManagers)
		}
	}

	return e
}
