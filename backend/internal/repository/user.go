package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) DeleteUser(ctx context.Context, userID any) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var email string
	if err := tx.QueryRowContext(ctx, `SELECT email FROM users WHERE id = $1`, userID).Scan(&email); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	statements := []struct {
		query string
		args  []any
	}{
		{query: `DELETE FROM sessions WHERE user_id = $1`, args: []any{userID}},
		{query: `DELETE FROM accounts WHERE user_id = $1`, args: []any{userID}},
		{query: `DELETE FROM user_groups WHERE user_id = $1`, args: []any{userID}},
		{query: `DELETE FROM verifications WHERE subject = $1`, args: []any{"email_verification::" + email}},
		{query: `UPDATE memberships SET transaction_id = NULL WHERE user_id = $1`, args: []any{userID}},
		{query: `DELETE FROM transactions WHERE user_id = $1`, args: []any{userID}},
		{query: `DELETE FROM memberships WHERE user_id = $1`, args: []any{userID}},
		{query: `DELETE FROM users WHERE id = $1`, args: []any{userID}},
	}

	for _, statement := range statements {
		if _, err := tx.ExecContext(ctx, statement.query, statement.args...); err != nil {
			return fmt.Errorf("delete user step failed: %w", err)
		}
	}

	return tx.Commit()
}
