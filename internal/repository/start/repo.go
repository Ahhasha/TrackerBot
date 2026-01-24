package start

import (
	"context"
	"errors"
	"fmt"

	"github.com/Ahhasha/Tracker-bot/internal/contracts"
	"github.com/jackc/pgx/v5"
)

type Repo struct{}

func NewRepo(db contracts.DBTX) *Repo {
	return &Repo{}
}

func (r *Repo) UpsertUser(ctx context.Context, db contracts.DBTX, tgID int64, username string) (int64, bool, error) {
	const op = "repo.start.UpsertUser"

	const insertQ = `
		INSERT INTO users (tg_id, username)
		VALUES ($1, $2)
		ON CONFLICT (tg_id) DO NOTHING
		RETURNING id;
	`

	var id int64
	err := db.QueryRow(ctx, insertQ, tgID, username).Scan(&id)
	if err == nil {
		return id, true, nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return 0, false, fmt.Errorf("%s: insert returning id: %w", op, err)
	}

	const updateQ = `
		UPDATE users
		SET username = $2
		WHERE tg_id = $1
		RETURNING id;
	`

	if err := db.QueryRow(ctx, updateQ, tgID, username).Scan(&id); err != nil {
		return 0, false, fmt.Errorf("%s: update returning id: %w", op, err)
	}

	return id, false, nil
}

func (r *Repo) CreateDefaultCategories(ctx context.Context, db contracts.DBTX, userID int64) ([]string, error) {
	const op = "repo.start.CreateDefaultCategory"
	defaults := []string{"Еда", "Транспорт", "Прочее", "Развлечения"}

	const q = `
	INSERT INTO categories (user_id, name)
	VALUES ($1, $2)
	ON CONFLICT (user_id, name) DO NOTHING
	`

	for _, name := range defaults {
		if _, err := db.Exec(ctx, q, userID, name); err != nil {
			return nil, fmt.Errorf("%s: insert %q: %w", op, name, err)
		}
	}
	return defaults, nil
}
