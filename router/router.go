package router

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/service"
)

type Handlers struct {
	Repo         model.Repository
	Service      service.Service
	SessionName  string
	SessionStore sessions.Store
}

func (h Handlers) Setup(e *echo.Echo) {
	api := e.Group("/api")
	{
		apiApplications := api.Group("/applications", AuthUserMiddleware)
		{
			apiApplications.GET("", GetApplicationList)
			apiApplications.POST("", PostApplication)
			apiApplications.GET("/:applicationId", GetApplication)
			apiApplications.PATCH("/:applicationId", PatchApplication)
		}

		apiImages := api.Group("/images", AuthUserMiddleware)
		{
			apiImages.GET("/:imageId", GetImages)
			apiImages.DELETE("/:imageId", DeleteImages)
		}

		apiComments := api.Group("/applications/:applicationId/comments", AuthUserMiddleware)
		{
			apiComments.POST("", PostComments)
			apiComments.PUT("/:commentId", PutComments)
			apiComments.DELETE("/:commentId", DeleteComments)
		}

		apiStatus := api.Group("/applications/:applicationId/states", AuthUserMiddleware)
		{
			apiStatus.PUT("", PutStates)
			apiStatus.PUT("/repaid/:repaidToId", PutRepaidStates)
		}

		apiUsers := api.Group("/users", AuthUserMiddleware)
		{
			apiUsers.GET("", GetUsers)
			apiUsers.GET("/me", GetMyUser)
			apiUsers.GET("/admins", GetAdminUsers)
			apiUsers.PUT("/admins", PutAdminUsers)
		}

		apiAuth := api.Group("/auth")
		{
			apiAuth.GET("/callback", AuthCallback)
			apiAuth.GET("/genpkce", GeneratePKCE)
		}
	}

}
