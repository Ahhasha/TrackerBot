package expense

import (
	"context"
	"fmt"
	"time"

	"github.com/Ahhasha/Tracker-bot/internal/contracts"
	expense "github.com/Ahhasha/Tracker-bot/internal/contracts/expense"
	"github.com/Ahhasha/Tracker-bot/internal/model"
)

type postgresRepo struct{}

func NewPostgresRepo() expense.ExpenseRepository {
	return &postgresRepo{}
}

func (r *postgresRepo) Create(ctx context.Context, db contracts.DBTX, expense model.Expense) (int64, error) {
	const op = "repo.expense.Create"

	const q = `
		INSERT INTO expenses (user_id, amount, category_id, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var id int64
	err := db.QueryRow(ctx, q, expense.UserID, expense.Amount, expense.CategoryID, expense.Description).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: scan: %w", op, err)
	}

	return id, nil
}

func (r *postgresRepo) GetPeriodWithCategory(ctx context.Context, db contracts.DBTX, userID int64, start time.Time, end time.Time) ([]model.ExpenseWithCategory, error) {
	const op = "repo.expense.GetByPeriodWithCategory"

	const q = `
		SELECT e.amount, e.description, c.name, e.created_at
		FROM expenses e JOIN categories c ON c.id = e.category_id
		WHERE e.user_id = $1
		AND e.created_at >= $2
		AND e.created_at <  $3
		ORDER BY e.created_at DESC
	`

	rows, err := db.Query(ctx, q, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("%s: query: %w", op, err)
	}
	defer rows.Close()

	result := make([]model.ExpenseWithCategory, 0)
	for rows.Next() {
		var row model.ExpenseWithCategory
		if err := rows.Scan(&row.Amount, &row.Description, &row.Category, &row.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s: scan: %w", op, err)
		}
		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows: %w", op, err)
	}

	return result, nil
}
