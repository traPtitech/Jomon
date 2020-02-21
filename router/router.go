package router

import (
	"github.com/labstack/echo/v4"
)

func SetRouting(e *echo.Echo) {

	apiApplications := e.Group("/applications")
	{
		apiApplications.GET("", GetApplicationList)
		apiApplications.POST("", PostApplication)
		apiApplications.GET("/:applicationId", GetApplication)
		apiApplications.PATCH("/:applicationId", PatchApplication)
	}

	apiImages := e.Group("/images")
	{
		apiImages.GET("/:imageId", GetImages)
		apiImages.DELETE("/:imageId", DeleteImages)
	}

	apiComments := e.Group("/applications/:applicationId/comments")
	{
		apiComments.POST("", PostComments)
		apiComments.PUT("/:commentId", PutComments)
		apiComments.DELETE("/:commentId", DeleteComments)
	}

	apiStatus := e.Group("/application/:applicationId/status")
	{
		apiStatus.PUT("", PutStatus)
		apiStatus.PUT("/repaid/:repaidTold", PutRepaidStatus)
	}

	apiUsers := e.Group("/users")
	{
		apiUsers.GET("", GetUsers)
		apiUsers.GET("/me", GetMyUser)
		apiUsers.GET("/admins", GetAdminUsers)
		apiUsers.PUT("/admins", PutAdminUsers)
	}
}
