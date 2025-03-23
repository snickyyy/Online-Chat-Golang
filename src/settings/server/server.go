package server

import (
	_ "libs/src/docs"
	handler_api "libs/src/internal/handlers/api"
	handler_middlewares "libs/src/internal/handlers/middlewares"

	files "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var middlewares = []gin.HandlerFunc{
	handler_middlewares.DependenciesMiddleware,
	handler_middlewares.AuthMiddleware,
	handler_middlewares.ErrorHandler,
}

func newServer(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:           ":8000",
		Handler:        handler,
		ReadTimeout:    2 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    12 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

// swag init --parseDependency --parseInternal --output ./docs

// @title           Online-Chat API
// @version         1.0
// @description     Documentation for the Online-Chat API
// @termsOfService  http://swagger.io/terms/

// @contact.name   Support
// @contact.url    http://www.blabla.com/support
// @contact.email  support@blabla.com

// @license.name  MIT
// @license.url   http://opensource.org/licenses/MIT

// @host      127.0.0.1:8000
// @BasePath  /api/v1
func RunServer() {
	router := gin.Default()

	router.Use(middlewares...)

	router.GET("/docs/*any", swagger.WrapHandler(files.Handler))

	router.GET("/", handler_api.Index)

	accounts := router.Group("/accounts")
	{
		auth := accounts.Group("/auth")
		{
			auth.POST("/register", handler_api.Register)
			auth.GET("/confirm-account/:token", handler_api.ConfirmAccount)
			auth.POST("/login", handler_api.Login)
			auth.DELETE("/logout", handler_api.Logout)
		}
		profile := accounts.Group("/profile")
		{
			profile.GET("/:username", handler_api.UserProfile)
			profile.PATCH("/edit", handler_api.ChangeUserProfile)
			profile.PUT("/reset-password", handler_api.ResetPassword)
			profile.PUT("/reset-password/confirm/:token", handler_api.ConfirmResetPassword)
		}
	}

	server := newServer(router)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
