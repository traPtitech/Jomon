package router

import (
	"encoding/gob"
	"errors"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"github.com/traPtitech/Jomon/logging"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/router/wrapsession"
	"github.com/traPtitech/Jomon/service"
	"github.com/traPtitech/Jomon/storage"
)

type Handlers struct {
	WebhookService *service.WebhookService
	Repository     model.Repository
	Storage        storage.Storage
	SessionName    string
}

func (h Handlers) NewServer(logger *zap.Logger) *echo.Echo {
	e := echo.New()
	e.Debug = os.Getenv("IS_DEBUG_MODE") != ""
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		logger := logging.GetLogger(c.Request().Context())
		var httpErr *echo.HTTPError
		var getSessionErr *wrapsession.GetSessionError
		var saveSessionErr *wrapsession.SaveSessionError
		var retErr error
		if errors.As(err, &httpErr) {
			retErr = httpErr
		} else if errors.As(err, &getSessionErr) {
			inner := getSessionErr.Unwrap()
			logger.Error("failed to get session", zap.Error(inner))
			retErr = echo.ErrInternalServerError.WithInternal(inner)
		} else if errors.As(err, &saveSessionErr) {
			inner := saveSessionErr.Unwrap()
			logger.Error("failed to save session", zap.Error(inner))
			retErr = echo.ErrInternalServerError.WithInternal(inner)
		} else {
			retErr = err
		}
		c.Echo().DefaultHTTPErrorHandler(retErr, c)
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
			apiAuth.GET("/genpkce", h.GeneratePKCE)
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
			apiUsers.PUT("", h.UpdateUserInfo, h.CheckAccountManagerMiddleware)
			apiUsers.GET("/me", h.GetMe)
		}

		apiAccountManagers := api.Group(
			"/account-managers",
			h.CheckLoginMiddleware,
			h.CheckAccountManagerMiddleware)
		{
			apiAccountManagers.GET("", h.GetAccountManagers)
			apiAccountManagers.POST("", h.PostAccountManagers)
			apiAccountManagers.DELETE("", h.DeleteAccountManagers)
		}
	}

	return e
}
