package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Uranury/RBK_finalProject/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE username = $1", username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) FindByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetBalance(ctx context.Context, userID uuid.UUID) (float64, error) {
	var balance float64
	err := r.db.GetContext(ctx, &balance, "SELECT balance FROM users WHERE id = $1", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return balance, nil
}

func (r *repository) UpdateBalance(ctx context.Context, tx *sqlx.Tx, userID uuid.UUID, newBalance float64) error {
	_, err := tx.ExecContext(ctx, "UPDATE users SET balance = $1 WHERE id = $2", newBalance, userID)
	return err
}

func (r *repository) Create(ctx context.Context, user *models.User) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO users (id, name, email, password, balance, role, created_at, updated_at) 
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		user.ID, user.Name, user.Email, user.Password, user.Balance, user.Role, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *repository) Delete(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", userID)
	return err
}
