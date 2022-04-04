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

	retrieveGroupOwner := h.RetrieveGroupOwner(h.Repository)
	retrieveRequestCreater := h.RetrieveRequestCreater(h.Repository)

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
			apiRequestIDs := api.Group("/:requestID", retrieveRequestCreater)
			{
				apiRequestIDs.GET("", h.GetRequest)
				apiRequestIDs.PUT("", h.PutRequest, middleware.BodyDump(service.WebhookEventHandler), h.CheckRequestCreaterMiddleware)
				apiRequestIDs.POST("/comments", h.PostComment, middleware.BodyDump(service.WebhookEventHandler))
				apiRequestIDs.PUT("/status", h.PutStatus, h.CheckAdminOrRequestCreaterMiddleware)
			}
		}

		apiComments := api.Group("/transactions", h.CheckLoginMiddleware)
		{
			apiComments.GET("", h.GetTransactions)
			apiComments.POST("", h.PostTransaction, h.CheckAdminMiddleware)
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
			apiGroupIDs := api.Group("/:groupID", retrieveGroupOwner)
			{
				apiGroupIDs.PUT("", h.PutGroup, h.CheckAdminOrGroupOwnerMiddleware)
				apiGroupIDs.DELETE("", h.DeleteGroup, h.CheckAdminOrGroupOwnerMiddleware)
				apiGroupIDs.GET("/members", h.GetMembers)
				apiGroupIDs.POST("/members", h.PostMember, h.CheckAdminOrGroupOwnerMiddleware)
				apiGroupIDs.DELETE("/members", h.DeleteMember, h.CheckAdminOrGroupOwnerMiddleware)
				apiGroupIDs.GET("/owners", h.GetOwners)
				apiGroupIDs.POST("/owners", h.PostOwner, h.CheckAdminOrGroupOwnerMiddleware)
				apiGroupIDs.DELETE("/owners", h.DeleteOwner, h.CheckAdminOrGroupOwnerMiddleware)
			}
		}

		apiUsers := api.Group("/users", h.CheckLoginMiddleware)
		{
			apiUsers.GET("", h.GetUsers)
			apiUsers.PUT("", h.UpdateUserInfo, h.CheckAdminMiddleware)
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
