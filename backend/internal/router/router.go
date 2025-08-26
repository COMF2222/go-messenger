package router

import (
	"github.com/COMF2222/go-messenger/internal/handler"
	"github.com/COMF2222/go-messenger/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Deps struct {
	AuthHandler    *handler.AuthHandler
	MessageHandler *handler.MessageHandler
	WsHandler      *handler.WSHandler
}

func SetupRouter(deps Deps) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		// Публичные маршруты
		api.POST("/register", deps.AuthHandler.Register)
		api.POST("/login", deps.AuthHandler.Login)

		api.GET("/online/:id", handler.GetOnlineStatus)
		api.GET("/ws", deps.WsHandler.ServerWS)

		// Защищённые маршруты
		auth := api.Group("/")
		auth.Use(middleware.AuthMiddleware())
		auth.POST("/messages", deps.MessageHandler.SendMessage)
		auth.GET("/messages", deps.MessageHandler.GetMessages)
		auth.GET("/interlocutors", deps.MessageHandler.GetInterlocutors)
	}

	return r
}
