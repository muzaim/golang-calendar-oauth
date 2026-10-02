package routes

import (
	"net/http"

	"golang-test/config"
	"golang-test/internal/handler"
	"golang-test/internal/middleware"
	"golang-test/pkg/response"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	cfg *config.Config,
	authHandler *handler.AuthHandler,
	calendarHandler *handler.CalendarHandler,
) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		response.Success(c, http.StatusOK, "Server is healthy", gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	{
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/refresh", authHandler.RefreshToken)
			authGroup.POST("/logout", authHandler.Logout)
			authGroup.GET("/google", authHandler.GoogleLogin)
			authGroup.GET("/google/callback", authHandler.GoogleCallback)

			protectedAuth := authGroup.Group("")
			protectedAuth.Use(middleware.AuthMiddleware(cfg.JWTSecret))
			{
				protectedAuth.GET("/me", authHandler.Me)
				protectedAuth.GET("/me/events", authHandler.MeWithEvents)
			}
		}

		calendarGroup := api.Group("/calendar")
		calendarGroup.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			calendarGroup.GET("/events", calendarHandler.GetEvents)
			calendarGroup.GET("/events/:id", calendarHandler.GetEventByID)
			calendarGroup.POST("/events", calendarHandler.CreateEvent)
			calendarGroup.PUT("/events/:id", calendarHandler.UpdateEvent)
			calendarGroup.DELETE("/events/:id", calendarHandler.DeleteEvent)
		}
	}

	return r
}
