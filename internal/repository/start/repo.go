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

func (r *Repo) UpsertUser(ctx context.Context, tgID int64, username string) (int64, bool, error) {
	const op = "repo.start.UpsertUser"

	const q = `
		INSERT INTO users (tg_id, username)
		VALUES ($1, $2)
		ON CONFLICT (tg_id) DO UPDATE
		SET username = EXCLUDED.username
		RETURNING id, (xmax = 0) AS created;
	`

	var id int64
	var created bool

	if err := r.db.QueryRow(ctx, q, tgID, username).Scan(&id, &created); err != nil {
		return 0, false, fmt.Errorf("%s: queryrow scan: %w", op, err)
	}

	return id, created, nil
}

func (r *Repo) CreateDefaultCategory(ctx context.Context, userID int64) error {
	const op = "repo.start.CreateDefaultCategory"
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
			return fmt.Errorf("%s: Exec %q: %w", op, name, err)
		}
	}
	return nil
}
