package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/COMF2222/go-messenger/internal/model"
	"github.com/COMF2222/go-messenger/internal/service"
	"github.com/COMF2222/go-messenger/internal/ws"
	"github.com/gin-gonic/gin"
)

type WSHandler struct {
	Hub            *ws.Hub
	MessageService *service.MessageService
}

func NewWSHandler(hub *ws.Hub, msgService *service.MessageService) *WSHandler {
	return &WSHandler{
		Hub:            hub,
		MessageService: msgService,
	}
}

type IncomingMessage struct {
	ToUserID int    `json:"to_user_id"`
	Text     string `json:"text"`
}

func (h *WSHandler) ServerWS(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "отсутствует user_id в query"})
		return
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "невалидный user_id"})
		return
	}

	conn, err := ws.Upgrade(c.Writer, c.Request)
	if err != nil {
		log.Println("Ошибка upgrade:", err)
		return
	}

	client := &ws.Client{
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte),
	}

	h.Hub.Register <- client

	go h.writePump(client)
	h.readPump(client)
}

func (h *WSHandler) readPump(client *ws.Client) {
	defer func() {
		h.Hub.Unregister <- client
		client.Conn.Close()
	}()

	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println("readPump ошибка:", err)
			break
		}

		var incoming IncomingMessage
		if err := json.Unmarshal(msg, &incoming); err != nil {
			log.Println("Невалидный формат сообщения:", err)
			continue
		}

		message := &model.Message{
			SenderID:   client.UserID,
			ReceiverID: incoming.ToUserID,
			Text:       incoming.Text,
		}

		if err := h.MessageService.SendMessage(context.Background(), message); err != nil {
			log.Println("Ошибка сохранения сообщения:", err)
			continue
		}

		wsMsg := ws.WSMessage{
			FromUserID: client.UserID,
			ToUserID:   incoming.ToUserID,
			Message:    msg,
		}
		h.Hub.Broadcast <- wsMsg
	}
}

func (h *WSHandler) writePump(client *ws.Client) {
	for msg := range client.Send {
		err := client.Conn.WriteMessage(1, msg)
		if err != nil {
			break
		}
	}
}
