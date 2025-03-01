package server

import (
	_ "libs/src/docs"
	handler_api "libs/src/internal/handlers/api"
	handler_middlewares "libs/src/internal/handlers/middlewares"
	swagger "github.com/swaggo/gin-swagger"
	files "github.com/swaggo/files"

	"net/http"
	"time"
	"github.com/gin-gonic/gin"
)

var middlewares = []gin.HandlerFunc{
	handler_middlewares.AuthMiddleware,
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
		}
	}

	server := newServer(router)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
