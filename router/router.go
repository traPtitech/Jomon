package router

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/service"
)

type Handlers struct {
	Repository   model.Repository
	Service      service.Service
	SessionName  string
	SessionStore sessions.Store
}

func SetRouting(e *echo.Echo, h Handlers) {
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
			apiTags.GET("/:tagID", h.GetTag)
			apiTags.PUT("/:tagID", h.PutTag)
			apiTags.DELETE("/:tagID", h.DeleteTag)
		}

		apiGroups := api.Group("/groups", h.AuthUserMiddleware)
		{
			apiGroups.GET("", h.GetGroups)
			apiGroups.POST("", h.PostGroup)
			apiGroups.GET("/:groupID", h.GetGroup)
			apiGroups.POST("/:groupID", h.PutGroup)
			apiGroups.PUT("/:groupID", h.PostGroupUser)
			apiGroups.DELETE("/:groupID", h.DeleteGroup)
		}

		apiAdmins := api.Group("/admins", h.AuthUserMiddleware)
		{
			apiAdmins.GET("", h.GetAdmins)
			apiAdmins.POST("", h.PostAdmin)
			apiAdmins.DELETE("/:userID", h.DeleteAdmin)
		}
	}
}
