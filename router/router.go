package router

import (
	"github.com/labstack/echo/v4"
)

func SetRouting(e *echo.Echo) {
	api := e.Group("/api")
	{
		apiApplications := api.Group("/applications")
		{
			apiApplications.GET("", GetApplications)
			apiApplications.POST("", PostApplications)
			apiApplications.PATCH("/:applicationId", PatchApplications)
		}
		apiImages := api.Group("/images")
		{
			apiImages.GET("/:imageId", GetImages)
			apiImages.DELETE("/:imageId", DeleteImages)
		}

		apiComments := api.Group("/applications/:applicationId/comments")
		{
			apiComments.POST("", PostComments)
			apiComments.PUT("/:commentId", PutComments)
			apiComments.DELETE("/:commentId", DeleteComments)
		}

		apiStatus := api.Group("/application/:applicationId/status")
		{
			apiStatus.PUT("", PutStatus)
			apiStatus.PUT("/repaid/:repaidTold", PutRepaidStatus)
		}

		apiUsers := api.Group("/users")
		{
			apiUsers.GET("", GetUsers)
			apiUsers.GET("/me", GetMyUser)
			apiUsers.GET("/admins", GetAdminUsers)
			apiUsers.PUT("/admins", PutAdminUsers)
		}
	}
}
