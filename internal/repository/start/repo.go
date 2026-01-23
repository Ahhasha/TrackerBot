package start

import (
	"context"
	"fmt"

	"github.com/Ahhasha/Tracker-bot/internal/contracts"
)

type Repo struct {
	db contracts.DBTX
}

func NewRepo(db contracts.DBTX) *Repo {
	return &Repo{db: db}
}

func (r *Repo) UserExists(ctx context.Context, tgID int64) (bool, error) {
	const q = `SELECT EXISTS(SELECT 1 FROM users WHERE tg_id = $1)`

	var exists bool
	err := r.db.QueryRow(ctx, q, tgID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check user exists %d: %w", tgID, err)
	}
	return exists, nil
}

func (r *Repo) Create(ctx context.Context, tgID int64, username string) (int64, error) {
	const q = `
INSERT INTO users (tg_id, username)
VALUES ($1, $2)
RETURNING id, tg_id, username, created_at
`

	var userID int64
	err := r.db.QueryRow(ctx, q, tgID, username).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("create user %d, %s: %w", tgID, username, err)
	}
	return userID, nil
}

func (r *Repo) CreateDefaultCategory(ctx context.Context, userID int64) error {
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
		if _, err := r.db.Exec(ctx, q, userID, name); err != nil {
			return fmt.Errorf("insert default category %q: %w", name, err)
		}
	}
	return nil
}
