package router

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/traPtitech/Jomon/model"
)

type Service struct {
	Administrators model.AdministratorRepository
	Applications   model.ApplicationRepository
	Comments       model.CommentRepository
	Images         model.ApplicationsImageRepository
	Users          model.UserRepository
	TraQAuth       model.TraQAuthRepository
	Webhook        model.WebhookRepository
}

func (s *Service) AuthUserMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c, err := s.AuthUser(c)
		if c == nil || err != nil {
			return err
		}
		return next(c)
	}
}

func SetRouting(e *echo.Echo, service Service) {
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:  "client/dist",
		HTML5: true,
	}))
	e.Use(middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{
		Skipper: func(c echo.Context) bool {
			if strings.HasPrefix(c.Request().URL.String(), "/api/applications") && c.Request().Method == "POST" {
				return false
			}
			return true
		},
		Handler: service.Webhook.WebhookEventHandler,
	}))

	api := e.Group("/api")
	{
		apiApplications := api.Group("/applications", service.AuthUserMiddleware)
		{
			apiApplications.GET("", service.GetApplicationList)
			apiApplications.POST("", service.PostApplication)
			apiApplications.GET("/:applicationId", service.GetApplication)
			apiApplications.PATCH("/:applicationId", service.PatchApplication)
		}

		apiImages := api.Group("/images", service.AuthUserMiddleware)
		{
			apiImages.GET("/:imageId", service.GetImages)
			apiImages.DELETE("/:imageId", service.DeleteImages)
		}

		apiComments := api.Group("/applications/:applicationId/comments", service.AuthUserMiddleware)
		{
			apiComments.POST("", service.PostComments)
			apiComments.PUT("/:commentId", service.PutComments)
			apiComments.DELETE("/:commentId", service.DeleteComments)
		}

		apiStatus := api.Group("/applications/:applicationId/states", service.AuthUserMiddleware)
		{
			apiStatus.PUT("", service.PutStates)
			apiStatus.PUT("/repaid/:repaidToId", service.PutRepaidStates)
		}

		apiUsers := api.Group("/users", service.AuthUserMiddleware)
		{
			apiUsers.GET("", service.GetUsers)
			apiUsers.GET("/me", service.GetMyUser)
			apiUsers.GET("/admins", service.GetAdminUsers)
			apiUsers.PUT("/admins", service.PutAdminUsers)
		}

		apiAuth := api.Group("/auth")
		{
			apiAuth.GET("/callback", service.AuthCallback)
			apiAuth.GET("/genpkce", service.GeneratePKCE)
		}
	}

}
