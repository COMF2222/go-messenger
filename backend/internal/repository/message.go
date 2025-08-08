package repository

import (
	"context"
	"fmt"

	"github.com/COMF2222/go-messenger/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageRepository struct {
	db *pgxpool.Pool
}

func NewMessageRepository(db *pgxpool.Pool) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) SaveMessage(ctx context.Context, msg *model.Message) error {
	query := `
		INSERT INTO messages (sender_id, receiver_id, text, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at;
	`
	err := r.db.QueryRow(ctx, query, msg.SenderID, msg.ReceiverID, msg.Text, msg.Status).Scan(&msg.ID, &msg.CreatedAt)
	if err != nil {
		return fmt.Errorf("SaveMessage: %w", err)
	}
	return nil
}

func (r *MessageRepository) GetMessageBetween(ctx context.Context, user1ID, user2ID int) ([]*model.Message, error) {
	query := `
		SELECT id, sender_id, receiver_id, text, status, created_at
		from messages
		WHERE (sender_id = $1 and receiver_id = $2)
			OR (sender_id = $2 and receiver_id = $1)
		ORDER BY created_at ASC;
	`
	rows, err := r.db.Query(ctx, query, user1ID, user2ID)
	if err != nil {
		return nil, fmt.Errorf("GetMessageBetween: %w", err)
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		var msg model.Message
		if err := rows.Scan(
			&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Text, &msg.Status, &msg.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		messages = append(messages, &msg)
	}

	return messages, nil
}

func (r *MessageRepository) GetInterlocutors(ctx context.Context, userID int) ([]int, error) {
	query := `
		SELECT DISTINCT
			CASE
				WHEN sender_id = $1 THEN receiver_id
				ELSE sender_id
			END AS interclocutor_id
		FROM messages
		WHERE sender_id = $1 OR receiver_id = $1;
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("GetInterlocutors query: %w", err)
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		ids = append(ids, id)
	}

	fmt.Println("Interlocutors:", ids)
	return ids, nil
}
