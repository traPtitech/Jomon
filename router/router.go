package router

import (
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/model"
)

type Service struct {
	Administrators model.AdministratorRepository
	Applications   model.ApplicationRepository
	Comments       model.CommentRepository
	Images         model.ApplicationsImageRepository
	Users          model.UserRepository
}

func SetRouting(e *echo.Echo, service Service) {
	e.Use(service.AuthUser)

	apiApplications := e.Group("/applications")
	{
		apiApplications.GET("", service.GetApplicationList)
		apiApplications.POST("", service.PostApplication)
		apiApplications.GET("/:applicationId", service.GetApplication)
		apiApplications.PATCH("/:applicationId", service.PatchApplication)
	}

	apiImages := e.Group("/images")
	{
		apiImages.GET("/:imageId", service.GetImages)
		apiImages.DELETE("/:imageId", service.DeleteImages)
	}

	apiComments := e.Group("/applications/:applicationId/comments")
	{
		apiComments.POST("", service.PostComments)
		apiComments.PUT("/:commentId", service.PutComments)
		apiComments.DELETE("/:commentId", service.DeleteComments)
	}

	apiStatus := e.Group("/application/:applicationId/states")
	{
		apiStatus.PUT("", PutStates)
		apiStatus.PUT("/repaid/:repaidTold", PutRepaidStates)
	}

	apiUsers := e.Group("/users")
	{
		apiUsers.GET("", service.GetUsers)
		apiUsers.GET("/me", service.GetMyUser)
		apiUsers.GET("/admins", service.GetAdminUsers)
		apiUsers.PUT("/admins", service.PutAdminUsers)
	}
}
