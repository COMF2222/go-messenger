package service

import (
	"context"

	"github.com/COMF2222/go-messenger/internal/model"
	"github.com/COMF2222/go-messenger/internal/repository"
)

type MessageService struct {
	repo *repository.MessageRepository
}

func NewMessageService(repo *repository.MessageRepository) *MessageService {
	return &MessageService{repo: repo}
}

func (s *MessageService) SendMessage(ctx context.Context, msg *model.Message) error {
	msg.Status = "sent"
	return s.repo.SaveMessage(ctx, msg)
}

func (s *MessageService) GetConversation(ctx context.Context, userID1, userID2 int) ([]*model.Message, error) {
	return s.repo.GetMessageBetween(ctx, userID1, userID2)
}
