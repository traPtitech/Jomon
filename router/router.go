package router

import (
	"encoding/gob"
	"os"

	"github.com/google/uuid"
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
	Repository  model.Repository
	Storage     storage.Storage
	Logger      *zap.Logger
	SessionName string
}

func NewServer(h Handlers) *echo.Echo {
	e := echo.New()
	e.Debug = os.Getenv("IS_DEBUG_MODE") != ""
	e.Use(h.AccessLoggingMiddleware(h.Logger))
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))))
	gob.Register(User{})
	gob.Register(uuid.UUID{})
	gob.Register([]*model.Owner{})

	retrieveGroupOwner := h.RetrieveGroupOwner()
	retrieveRequestCreator := h.RetrieveRequestCreator()
	retrieveFileCreator := h.RetrieveFileCreator()

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
			apiRequestIDs := apiRequests.Group("/:requestID", retrieveRequestCreator)
			{
				apiRequestIDs.GET("", h.GetRequest)
				apiRequestIDs.PUT("", h.PutRequest, middleware.BodyDump(service.WebhookEventHandler), h.CheckRequestCreatorMiddleware)
				apiRequestIDs.POST("/comments", h.PostComment, middleware.BodyDump(service.WebhookEventHandler))
				apiRequestIDs.PUT("/status", h.PutStatus, h.CheckAdminOrRequestCreatorMiddleware)
			}
		}

		apiTransactions := api.Group("/transactions", h.CheckLoginMiddleware)
		{
			apiTransactions.GET("", h.GetTransactions)
			apiTransactions.POST("", h.PostTransaction, middleware.BodyDump(service.WebhookEventHandler), h.CheckAdminMiddleware)
			apiTransactions.GET("/:transactionID", h.GetTransaction)
			apiTransactions.PUT("/:transactionID", h.PutTransaction, middleware.BodyDump(service.WebhookEventHandler), h.CheckAdminMiddleware)
		}

		apiFiles := api.Group("/files", h.CheckLoginMiddleware)
		{
			apiFiles.POST("", h.PostFile)
			apiFileIDs := apiFiles.Group("/:fileID", retrieveFileCreator)
			{
				apiFileIDs.GET("", h.GetFile)
				apiFileIDs.DELETE("", h.DeleteFile, h.CheckAdminOrFileCreatorMiddleware)
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

		apiGroups := api.Group("/groups", h.CheckLoginMiddleware)
		{
			apiGroups.GET("", h.GetGroups)
			apiGroups.POST("", h.PostGroup, h.CheckAdminMiddleware)
			apiGroupIDs := apiGroups.Group("/:groupID", retrieveGroupOwner)
			{
				apiGroupIDs.GET("", h.GetGroupDetail)
				apiGroupIDs.PUT("", h.PutGroup, h.CheckAdminOrGroupOwnerMiddleware)
				apiGroupIDs.DELETE("", h.DeleteGroup, h.CheckAdminOrGroupOwnerMiddleware)
				apiGroupIDs.POST("/members", h.PostMember, h.CheckAdminOrGroupOwnerMiddleware)
				apiGroupIDs.DELETE("/members", h.DeleteMember, h.CheckAdminOrGroupOwnerMiddleware)
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
			apiAdmins.POST("", h.PostAdmins)
			apiAdmins.DELETE("", h.DeleteAdmins)
		}
	}

	return e
}
