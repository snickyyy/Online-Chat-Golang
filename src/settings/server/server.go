package server

import (
	handler_api "libs/src/internal/handlers/api"
	handler_middlewares "libs/src/internal/handlers/middlewares"
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

func RunServer() {
	router := gin.Default()

	router.Use(middlewares...)

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
