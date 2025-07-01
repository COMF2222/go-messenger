package router

import (
	"github.com/COMF2222/go-messenger/internal/handler"
	"github.com/COMF2222/go-messenger/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Deps struct {
	AuthHandler    *handler.AuthHandler
	MessageHandler *handler.MessageHandler
}

func SetupRouter(deps Deps) *gin.Engine {
	r := gin.Default()

	// Публичные маршруты
	r.POST("/register", deps.AuthHandler.Register)
	r.POST("/login", deps.AuthHandler.Login)

	// Защищённые маршруты
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	auth.POST("/messages", deps.MessageHandler.SendMessage)
	auth.GET("/messages", deps.MessageHandler.GetMessages)

	return r
}
