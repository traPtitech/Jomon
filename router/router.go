package router

import (
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/service"
	"go.uber.org/zap"
)

type Handlers struct {
	Repository   model.Repository
	Logger       *zap.Logger
	Service      service.Services
	SessionName  string
	SessionStore sessions.Store
}

func SetRouting(e *echo.Echo, h Handlers) {
	e.Debug = (os.Getenv("IS_DEBUG_MODE") != "")
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

		apiRequests := api.Group("/requests", h.AuthUserMiddleware)
		{
			apiRequests.GET("", h.GetRequests)
			apiRequests.POST("", h.PostRequest)
			apiRequests.GET("/:requestID", h.GetRequest)
			apiRequests.PUT("/:requestID", h.PutRequest)
			apiRequests.POST("/:requestID/comments", h.PostComment)
			apiRequests.PUT("/:requestID/comments/:commentID", h.PutComment)
			apiRequests.DELETE("/:requestID/comments/:commentID", h.DeleteComment)
			apiRequests.PUT("/:requestID/status", h.PutStatus)
		}

		apiComments := api.Group("/transactions", h.AuthUserMiddleware)
		{
			apiComments.GET("", h.GetTransactions)
			apiComments.POST("", h.PostTransaction)
			apiComments.GET("/:transactionID", h.GetTransaction)
			apiComments.PUT("/:transactionID", h.PutTransaction)
		}

		apiFiles := api.Group("/files", h.AuthUserMiddleware)
		{
			apiFiles.POST("", h.PostFile)
			apiFiles.GET("/:fileID", h.GetFile)
			apiFiles.DELETE("/:fileID", h.DeleteFile)
		}

		apiTags := api.Group("/tags", h.AuthUserMiddleware)
		{
			apiTags.GET("", h.GetTags)
			apiTags.POST("", h.PostTag)
			apiTags.PUT("/:tagID", h.PutTag)
			apiTags.DELETE("/:tagID", h.DeleteTag)
		}

		apiGroups := api.Group("/groups", h.AuthUserMiddleware)
		{
			apiGroups.GET("", h.GetGroups)
			apiGroups.POST("", h.PostGroup)
			apiGroups.PUT("/:groupID", h.PutGroup)
			apiGroups.DELETE("/:groupID", h.DeleteGroup)
			apiGroups.GET("/:groupID/members", h.GetMembers)
			apiGroups.POST("/:groupID/members", h.PostMember)
			apiGroups.DELETE("/:groupID/members", h.DeleteMember)
			apiGroups.GET("/:groupID/owners", h.GetOwners)
			apiGroups.POST("/:groupID/owners", h.PostOwner)
			apiGroups.DELETE("/:groupID/owners", h.DeleteOwner)

		}

		apiUsers := api.Group("/users", h.AuthUserMiddleware)
		{
			apiUsers.GET("", h.GetUsers)
			apiUsers.PUT("", h.PutUsers)
			apiUsers.GET("/me", h.GetMe)
		}
	}
}
