package category

import (
	"context"
	"fmt"

	"github.com/Ahhasha/Tracker-bot/internal/contracts"
	"github.com/Ahhasha/Tracker-bot/internal/contracts/add"
	"github.com/Ahhasha/Tracker-bot/internal/model"
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
		return model.Category{}, fmt.Errorf("%s: scan: %w", op, err)
	}

	return cat, nil
}
