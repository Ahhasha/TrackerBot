package expense

import (
	"context"
	"fmt"
	"time"

	"github.com/Ahhasha/Tracker-bot/internal/contracts"
	"github.com/Ahhasha/Tracker-bot/internal/model"
)

type Repo struct{}

func NewRepo() *Repo {
	return &Repo{}
}

type CategoryStat struct {
	CategoryName     string
	TotalAmount      int64
	TransactionCount int64
}

func (r *Repo) Create(ctx context.Context, db contracts.DBTX, expense model.Expense) (int64, error) {
	const op = "repo.expense.Create"

	const q = `
		INSERT INTO expenses (user_id, amount, category_id, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var id int64
	err := db.QueryRow(ctx, q, expense.UserID, expense.Amount, expense.CategoryID, expense.Description).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *Repo) GetToday(ctx context.Context, db contracts.DBTX, userID int64) ([]model.Expense, error) {
	const op = "repo.expense.GetToday"

	start, end := getTodayRange()
	expenses, err := r.getByPeriod(ctx, db, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return expenses, nil
}

func (r *Repo) GetWeek(ctx context.Context, db contracts.DBTX, userID int64) ([]model.Expense, error) {
	const op = "repo.expense.GetWeek"

	start, end := getWeekRange()
	expenses, err := r.getByPeriod(ctx, db, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return expenses, nil
}

func (r *Repo) GetMonth(ctx context.Context, db contracts.DBTX, userID int64) ([]model.Expense, error) {
	const op = "repo.expense.GetMonth"

	start, end := getMonthRange()
	expenses, err := r.getByPeriod(ctx, db, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return expenses, nil
}

func (r *Repo) GetTotalSpent(ctx context.Context, db contracts.DBTX, userID int64) (int64, error) {
	const op = "repo.expense.GetTotalSpent"

	const q = `
		SELECT COALESCE(SUM(amount), 0)
		FROM expenses
		WHERE user_id = $1
	`

	var total int64
	err := db.QueryRow(ctx, q, userID).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return total, nil
}

func (r *Repo) GetTotalTransactions(ctx context.Context, db contracts.DBTX, userID int64) (int64, error) {
	const op = "repo.expense.GetTotalTransactions"

	const q = `
		SELECT COUNT(*)
		FROM expenses
		WHERE user_id = $1
	`

	var count int64
	err := db.QueryRow(ctx, q, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return count, nil
}

func (r *Repo) GetTopCategories(ctx context.Context, db contracts.DBTX, userID int64, limit int) ([]CategoryStat, error) {
	const op = "repo.expense.GetTopCategories"

	const q = `
		SELECT 
			c.name as category_name,
			SUM(e.amount) as total_amount,
			COUNT(e.id) as transaction_count
		FROM expenses e
		JOIN categories c ON e.category_id = c.id
		WHERE e.user_id = $1
		GROUP BY c.name
		ORDER BY total_amount DESC
		LIMIT $2
	`

	rows, err := db.Query(ctx, q, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: query: %w", op, err)
	}
	defer rows.Close()

	var stats []CategoryStat
	for rows.Next() {
		var s CategoryStat
		if err := rows.Scan(&s.CategoryName, &s.TotalAmount, &s.TransactionCount); err != nil {
			return nil, fmt.Errorf("%s: scan: %w", op, err)
		}
		stats = append(stats, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error: %w", op, err)
	}

	return stats, nil
}

func (r *Repo) getByPeriod(ctx context.Context, db contracts.DBTX, userID int64, start, end time.Time) ([]model.Expense, error) {
	const op = "repo.expense.getByPeriod"

	const q = `
		SELECT id, user_id, amount, category_id, description, created_at
		FROM expenses
		WHERE user_id = $1
		AND created_at BETWEEN $2 AND $3
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

func getTodayRange() (start, end time.Time) {
	now := time.Now()
	start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end = start.Add(24*time.Hour - 1*time.Second)
	return
}

func getWeekRange() (start, end time.Time) {
	now := time.Now()

	weekday := now.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	} else {
		weekday--
	}

	start = time.Date(now.Year(), now.Month(), now.Day()-int(weekday)+1, 0, 0, 0, 0, now.Location())
	end = start.Add(7*24*time.Hour - 1*time.Second)
	return
}

func getMonthRange() (start, end time.Time) {
	now := time.Now()

	start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	end = time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location()).
		Add(-1 * time.Second)

	return
}
