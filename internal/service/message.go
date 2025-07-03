package service

import (
	"context"

	"github.com/COMF2222/go-messenger/internal/model"
	"github.com/COMF2222/go-messenger/internal/repository"
)

type MessageService struct {
	repo     *repository.MessageRepository
	userRepo *repository.UserRepository
}

func NewMessageService(repo *repository.MessageRepository, userRepo *repository.UserRepository) *MessageService {
	return &MessageService{repo: repo, userRepo: userRepo}
}

func (s *MessageService) SendMessage(ctx context.Context, msg *model.Message) error {
	msg.Status = "sent"
	return s.repo.SaveMessage(ctx, msg)
}

func (s *MessageService) GetConversation(ctx context.Context, userID1, userID2 int) ([]*model.Message, error) {
	return s.repo.GetMessageBetween(ctx, userID1, userID2)
}

func (s *MessageService) GetInterlocutors(ctx context.Context, userID int) ([]int, error) {
	return s.repo.GetInterlocutors(ctx, userID)
}

func (s *MessageService) GetInterlocutorsUsers(ctx context.Context, userID int) ([]*model.User, error) {
	ids, err := s.repo.GetInterlocutors(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.userRepo.GetUsersByIDs(ctx, ids)
}
