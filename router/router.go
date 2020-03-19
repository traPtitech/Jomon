package router

import (
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/model"
)

type Service struct {
	Administrators model.AdministratorRepository
	Applications   model.ApplicationRepository
	Comments       model.CommentRepository
}

func SetRouting(e *echo.Echo) {
	service := &Service{
		Administrators: model.NewAdministratorRepository(),
		Applications:   model.NewApplicationRepository(),
		Comments:       model.NewCommentRepository(),
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
		apiUsers.GET("", GetUsers)
		apiUsers.GET("/me", GetMyUser)
		apiUsers.GET("/admins", GetAdminUsers)
		apiUsers.PUT("/admins", PutAdminUsers)
	}
}
