package category

import (
	"context"
	"errors"
	"fmt"

	"github.com/Ahhasha/Tracker-bot/internal/contracts"
	"github.com/Ahhasha/Tracker-bot/internal/contracts/add"
	"github.com/Ahhasha/Tracker-bot/internal/model"
	"github.com/jackc/pgx/v5"
)

type postgresRepo struct{}

func NewPostgresRepo() add.CategoryRepository {
	return &postgresRepo{}
}

func (r *postgresRepo) GetByName(ctx context.Context, db contracts.DBTX, userID int64, name string) (model.Category, error) {
	const op = "repo.category.GetByName"

	const q = `
		SELECT id, user_id, name, created_at
		FROM categories
		WHERE user_id = $1 AND LOWER(name) = LOWER($2)
	`

	var cat model.Category
	err := db.QueryRow(ctx, q, userID, name).Scan(&cat.ID, &cat.UserID, &cat.Name, &cat.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Category{}, fmt.Errorf("%s: not found: %w", op, err)
		}
		return model.Category{}, fmt.Errorf("%s: query: %w", op, err)
	}

	return cat, nil
}
