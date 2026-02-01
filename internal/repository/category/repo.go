package category

import (
	"context"
	"fmt"

	"github.com/Ahhasha/Tracker-bot/internal/contracts"
	"github.com/Ahhasha/Tracker-bot/internal/model"
)

type Repo struct{}

func NewRepo() *Repo {
	return &Repo{}
}

func (r *Repo) GetByName(ctx context.Context, db contracts.DBTX, userID int64, name string) (model.Category, error) {
	op := "repo.category.GetByName"

	const q = `
		SELECT id, user_id, name
		FROM categories
		WHERE user_id = $1 AND LOWER(name) = LOWER($2)
	`

	var cat model.Category
	err := db.QueryRow(ctx, q, userID, name).Scan(&cat.ID, &cat.UserID, &cat.Name)
	if err != nil {
		return model.Category{}, fmt.Errorf("%s: query: %w", op, err)
	}

	return cat, nil
}
