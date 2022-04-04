package router

import (
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/service"
	"github.com/traPtitech/Jomon/storage"
)

type Handlers struct {
	Repository   model.Repository
	Storage      storage.Storage
	Logger       *zap.Logger
	SessionName  string
	SessionStore sessions.Store
}

func NewServer(h Handlers) *echo.Echo {
	e := echo.New()
	e.Debug = os.Getenv("IS_DEBUG_MODE") != ""
	e.Use(h.AccessLoggingMiddleware(h.Logger))
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(session.Middleware(h.SessionStore))

	api := e.Group("/api")
	{
		apiAuth := api.Group("/auth")
		{
			apiAuth.GET("/callback", h.AuthCallback)
			apiAuth.GET("/genpkce", h.GeneratePKCE)
		}

		apiRequests := api.Group("/requests", h.CheckLoginMiddleware)
		{
			apiRequests.GET("", h.GetRequests)
			apiRequests.POST("", h.PostRequest, middleware.BodyDump(service.WebhookEventHandler))
			apiRequests.GET("/:requestID", h.GetRequest)
			apiRequests.PUT("/:requestID", h.PutRequest, middleware.BodyDump(service.WebhookEventHandler))
			apiRequests.POST("/:requestID/comments", h.PostComment, middleware.BodyDump(service.WebhookEventHandler))
			apiRequests.PUT("/:requestID/status", h.PutStatus)
		}

		apiComments := api.Group("/transactions", h.CheckLoginMiddleware)
		{
			apiComments.GET("", h.GetTransactions)
			apiComments.POST("", h.PostTransaction)
			apiComments.GET("/:transactionID", h.GetTransaction)
			apiComments.PUT("/:transactionID", h.PutTransaction, h.CheckAdminMiddleware)
		}

		apiFiles := api.Group("/files", h.CheckLoginMiddleware)
		{
			apiFiles.POST("", h.PostFile)
			apiFiles.GET("/:fileID", h.GetFile)
			apiFiles.DELETE("/:fileID", h.DeleteFile)
			apiFiles.GET("/:fileID/meta", h.GetFileMeta)
		}

		apiTags := api.Group("/tags", h.CheckLoginMiddleware)
		{
			apiTags.GET("", h.GetTags)
			apiTags.POST("", h.PostTag)
			apiTags.PUT("/:tagID", h.PutTag)
			apiTags.DELETE("/:tagID", h.DeleteTag)
		}

		apiGroups := api.Group("/groups", h.CheckLoginMiddleware)
		{
			apiGroups.GET("", h.GetGroups)
			apiGroups.POST("", h.PostGroup, h.CheckAdminMiddleware)
			apiGroups.PUT("/:groupID", h.PutGroup)
			apiGroups.DELETE("/:groupID", h.DeleteGroup)
			apiGroups.GET("/:groupID/members", h.GetMembers)
			apiGroups.POST("/:groupID/members", h.PostMember)
			apiGroups.DELETE("/:groupID/members", h.DeleteMember)
			apiGroups.GET("/:groupID/owners", h.GetOwners)
			apiGroups.POST("/:groupID/owners", h.PostOwner)
			apiGroups.DELETE("/:groupID/owners", h.DeleteOwner)
		}

		apiUsers := api.Group("/users", h.CheckLoginMiddleware)
		{
			apiUsers.GET("", h.GetUsers)
			apiUsers.PUT("", h.UpdateUserInfo)
			apiUsers.GET("/me", h.GetMe)
		}

		apiAdmins := api.Group("/admins", h.CheckLoginMiddleware, h.CheckAdminMiddleware)
		{
			apiAdmins.GET("", h.GetAdmins)
			apiAdmins.POST("", h.PostAdmin)
			apiAdmins.DELETE("/:userID", h.DeleteAdmin)
		}
	}

	return e
}
