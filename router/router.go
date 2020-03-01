package router

import (
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/model"
)

type Service struct {
	Applications *ApplicationService
}

func SetRouting(e *echo.Echo) {
	service := Service{
		Applications: NewApplicationService(model.NewApplicationRepository()),
	}

	apiApplications := e.Group("/applications")
	{
		apiApplications.GET("", service.GetApplicationList)
		apiApplications.POST("", service.PostApplication)
		apiApplications.GET("/:applicationId", service.GetApplication)
		apiApplications.PATCH("/:applicationId", service.PatchApplication)
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
