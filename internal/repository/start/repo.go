package registration

import (
	"context"
	"errors"
	"fmt"

	"github.com/Ahhasha/Tracker-bot/internal/model"
	"github.com/jackc/pgx/v5"
)

type Repo struct {
	tx pgx.Tx
}

func NewStart(tx pgx.Tx) *Repo {
	return &Repo{tx: tx}
}

func (r *Repo) FindOrCreateStart(ctx context.Context, tgID int64, username string) (*model.User, error) {
	u, err := r.getByTelegramID(ctx, tgID)
	if err != nil {
		return nil, err
	}
	if u != nil {
		return u, nil
	}

	u, err = r.create(ctx, tgID, username)
	if err != nil {
		return nil, err
	}

	if err := r.createDefaultCategory(ctx, u.ID); err != nil {
		return nil, err
	}
	return u, nil
}

func (r *Repo) getByTelegramID(ctx context.Context, tgID int64) (*model.User, error) {
	const q = `
SELECT id, tg_id, username, created_at
FROM users
WHERE tg_id = $1
`

	var u model.User
	err := r.tx.QueryRow(ctx, q, tgID).Scan(&u.ID, &u.TgID, &u.Username, &u.CreatedAt)
	if err == nil {
		return &u, nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return nil, fmt.Errorf("select user by tg_id: %w", err)
}

func (r *Repo) create(ctx context.Context, tgID int64, username string) (*model.User, error) {
	const q = `
INSERT INTO users (tg_id, username)
VALUES ($1, $2)
RETURNING id, tg_id, username, created_at
`

	var u model.User
	if err := r.tx.QueryRow(ctx, q, tgID, username).Scan(&u.ID, &u.TgID, &u.Username, &u.CreatedAt); err != nil {
		return nil, fmt.Errorf("insert user: %w", err)
	}

	return &u, nil
}

func (r *Repo) createDefaultCategory(ctx context.Context, userID int64) error {
	defaults := []string{
		"Еда",
		"Транспорт",
		"Жильё",
		"Развлечения",
	}

	const q = `
	INSERT INTO categories (user_id, name)
	VALUES ($1, $2)
	ON CONFLICT (user_id, name) DO NOTHING
	`

	for _, name := range defaults {
		if _, err := r.tx.Exec(ctx, q, userID, name); err != nil {
			return fmt.Errorf("insert default category %q: %w", name, err)
		}
	}
	return nil
}
