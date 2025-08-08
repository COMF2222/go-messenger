package repository

import (
	"context"

	"github.com/COMF2222/go-messenger/internal/model"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetUsersByIDs(ctx context.Context, ids []int) ([]*model.User, error)
}

type MessageRepositoryInterface interface {
	SaveMessage(ctx context.Context, msg *model.Message) error
	GetMessageBetween(ctx context.Context, user1ID, user2ID int) ([]*model.Message, error)
	GetInterlocutors(ctx context.Context, userID int) ([]int, error)
}
