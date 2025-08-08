package handler

import (
	"net/http"
	"strconv"

	"github.com/COMF2222/go-messenger/internal/session"
	"github.com/gin-gonic/gin"
)

func GetOnlineStatus(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID пользователя"})
		return
	}

	isOnline, err := session.IsUserOnline(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка Redis"})
	}

	c.JSON(http.StatusOK, gin.H{"online": isOnline})
}
