package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/COMF2222/go-messenger/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, created_at;`
	return r.db.QueryRow(ctx, query, user.Username, user.Email, user.PasswordHash).Scan(&user.ID, &user.CreatedAt)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, username, email, password_hash, created_at
		FROM users
		WHERE email = $1`
	row := r.db.QueryRow(ctx, query, email)

	var user model.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (r *UserRepository) GetUsersByIDs(ctx context.Context, ids []int) ([]*model.User, error) {
	if len(ids) == 0 {
		return []*model.User{}, nil
	}

	params := make([]string, len(ids))
	args := make([]interface{}, len((ids)))
	for i, id := range ids {
		params[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT id, username
		FROM users
		WHERE id IN (%s)
	`, strings.Join(params, ", "))

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		fmt.Println("GetUsersByIDs query err:", err)
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Username); err != nil {
			fmt.Println("GetUsersByIDs scan err:", err)
			return nil, err
		}
		users = append(users, &u)
	}

	return users, nil
}
