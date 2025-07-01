package main

import (
	"log"
	"os"

	"github.com/COMF2222/go-messenger/internal/config"
	"github.com/COMF2222/go-messenger/internal/handler"
	"github.com/COMF2222/go-messenger/internal/repository"
	"github.com/COMF2222/go-messenger/internal/router"
	"github.com/COMF2222/go-messenger/internal/service"
)

func main() {
	cfg := config.LoadConfig()

	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal("Ошибка подключения к БД: ", err)
	}
	defer db.Close()

	// DI — зависимости
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg)
	authHandler := handler.NewAuthHandler(authService)

	messageRepo := repository.NewMessageRepository(db)
	messageService := service.NewMessageService(messageRepo)
	messageHandler := handler.NewMessageHandler(messageService)

	r := router.SetupRouter(router.Deps{
		AuthHandler:    authHandler,
		MessageHandler: messageHandler,
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Println("Сервер запущен на порту " + port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}
