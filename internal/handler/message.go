package handler

import (
	"net/http"
	"strconv"

	"github.com/COMF2222/go-messenger/internal/model"
	"github.com/COMF2222/go-messenger/internal/service"
	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	messageService *service.MessageService
}

func NewMessageHandler(s *service.MessageService) *MessageHandler {
	return &MessageHandler{messageService: s}
}

type MessageInput struct {
	ReceiverID int    `json:"receiver_id"`
	Text       string `json:"text"`
}

// POST /messages
func (h *MessageHandler) SendMessage(c *gin.Context) {
	userID := c.GetInt("user_id")
	var input MessageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат"})
	}

	msg := &model.Message{
		SenderID:   userID,
		ReceiverID: input.ReceiverID,
		Text:       input.Text,
	}
	if err := h.messageService.SendMessage(c.Request.Context(), msg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось отправить сообщение"})
		return
	}
	c.JSON(http.StatusOK, msg)
}

// GET /messages?with=2
func (h *MessageHandler) GetMessages(c *gin.Context) {
	userID := c.GetInt("user_id")
	otherIDStr := c.Query("with")
	otherID, err := strconv.Atoi(otherIDStr)
	if err != nil || otherID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ?with"})
		return
	}

	msgs, err := h.messageService.GetConversation(c.Request.Context(), userID, otherID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении сообщений"})
		return
	}

	c.JSON(http.StatusOK, msgs)
}

// GET /interlocutors
func (h *MessageHandler) GetInterlocutors(c *gin.Context) {
	userID := c.GetInt("user_id")

	users, err := h.messageService.GetInterlocutorsUsers(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не удалось получить список собеседников"})
		return
	}

	c.JSON(http.StatusOK, users)
}
