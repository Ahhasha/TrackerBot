package expense

import (
	"context"
	"fmt"
	"time"

	"github.com/Ahhasha/Tracker-bot/internal/contracts"
	"github.com/Ahhasha/Tracker-bot/internal/contracts/add"
	"github.com/Ahhasha/Tracker-bot/internal/model"
)

type postgresRepo struct{}

func NewPostgresRepo() add.ExpenseRepository {
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

func (r *postgresRepo) GetToday(ctx context.Context, db contracts.DBTX, userID int64) ([]model.Expense, error) {
	const op = "repo.expense.GetToday"

	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 0, 1)

	expenses, err := r.getByPeriod(ctx, db, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return expenses, nil
}

func (r *postgresRepo) GetWeek(ctx context.Context, db contracts.DBTX, userID int64) ([]model.Expense, error) {
	const op = "repo.expense.GetWeek"

	now := time.Now()
	weekday := now.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	daysSinceMonday := int(weekday) - 1
	start := time.Date(now.Year(), now.Month(), now.Day()-daysSinceMonday, 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 0, 7)

	expenses, err := r.getByPeriod(ctx, db, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return expenses, nil
}

func (r *postgresRepo) GetMonth(ctx context.Context, db contracts.DBTX, userID int64) ([]model.Expense, error) {
	const op = "repo.expense.GetMonth"

	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 1, 0)

	expenses, err := r.getByPeriod(ctx, db, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return expenses, nil
}

func (r *postgresRepo) getByPeriod(ctx context.Context, db contracts.DBTX, userID int64, start, end time.Time) ([]model.Expense, error) {
	const op = "repo.expense.getByPeriod"

	const q = `
		SELECT id, user_id, amount, category_id, description, created_at
		FROM expenses
		WHERE user_id = $1
		AND created_at >= $2
		AND created_at < $3
		ORDER BY created_at DESC
	`

	rows, err := db.Query(ctx, q, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("%s: query: %w", op, err)
	}
	defer rows.Close()

	var expenses []model.Expense
	for rows.Next() {
		var e model.Expense
		if err := rows.Scan(&e.ID, &e.UserID, &e.Amount, &e.CategoryID, &e.Description, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s: scan: %w", op, err)
		}
		expenses = append(expenses, e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error: %w", op, err)
	}

	return expenses, nil
}
