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
	start, end := todayRange(time.Now())
	return r.getByPeriod(ctx, db, userID, start, end)
}

func (r *postgresRepo) GetWeek(ctx context.Context, db contracts.DBTX, userID int64) ([]model.Expense, error) {
	start, end := weekRange(time.Now())
	return r.getByPeriod(ctx, db, userID, start, end)
}

func (r *postgresRepo) GetMonth(ctx context.Context, db contracts.DBTX, userID int64) ([]model.Expense, error) {
	start, end := monthRange(time.Now())
	return r.getByPeriod(ctx, db, userID, start, end)
}

func (r *postgresRepo) GetTodayWithCategory(ctx context.Context, db contracts.DBTX, userID int64) ([]model.ExpenseWithCategory, error) {
	start, end := todayRange(time.Now())
	return r.getByPeriodWithCategory(ctx, db, userID, start, end)
}

func (r *postgresRepo) GetWeekWithCategory(ctx context.Context, db contracts.DBTX, userID int64) ([]model.ExpenseWithCategory, error) {
	start, end := weekRange(time.Now())
	return r.getByPeriodWithCategory(ctx, db, userID, start, end)
}

func (r *postgresRepo) GetMonthWithCategory(ctx context.Context, db contracts.DBTX, userID int64) ([]model.ExpenseWithCategory, error) {
	start, end := monthRange(time.Now())
	return r.getByPeriodWithCategory(ctx, db, userID, start, end)
}

func (r *postgresRepo) getByPeriod(ctx context.Context, db contracts.DBTX, userID int64, start, end time.Time) ([]model.Expense, error) {
	const q = `
		SELECT id, user_id, amount, category_id, description, created_at
		FROM expenses
		WHERE user_id = $1
		AND created_at >= $2
		AND created_at <  $3
		ORDER BY created_at DESC
	`

	rows, err := db.Query(ctx, q, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("repo.expense.getByPeriod: query: %w", err)
	}
	defer rows.Close()

	var expenses []model.Expense
	for rows.Next() {
		var e model.Expense
		if err := rows.Scan(&e.ID, &e.UserID, &e.Amount, &e.CategoryID, &e.Description, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("repo.expense.getByPeriod: scan: %w", err)
		}
		expenses = append(expenses, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repo.expense.getByPeriod: rows: %w", err)
	}

	return expenses, nil
}

func (r *postgresRepo) getByPeriodWithCategory(ctx context.Context, db contracts.DBTX, userID int64, start, end time.Time) ([]model.ExpenseWithCategory, error) {
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
		return nil, fmt.Errorf("repo.expense.getByPeriodWithCategory: query: %w", err)
	}
	defer rows.Close()

	var result []model.ExpenseWithCategory
	for rows.Next() {
		var row model.ExpenseWithCategory
		if err := rows.Scan(&row.Amount, &row.Description, &row.Category, &row.CreatedAt); err != nil {
			return nil, fmt.Errorf("repo.expense.getByPeriodWithCategory: scan: %w", err)
		}
		result = append(result, row)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repo.expense.getByPeriodWithCategory: rows: %w", err)
	}

	return result, nil
}

func todayRange(now time.Time) (time.Time, time.Time) {
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 0, 1)
	return start, end
}

func weekRange(now time.Time) (time.Time, time.Time) {
	weekday := now.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	daysSinceMonday := int(weekday) - 1
	start := time.Date(now.Year(), now.Month(), now.Day()-daysSinceMonday, 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 0, 7)
	return start, end
}

func monthRange(now time.Time) (time.Time, time.Time) {
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 1, 0)
	return start, end
}
