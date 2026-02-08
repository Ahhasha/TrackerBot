package add

import (
	"context"

	"github.com/Ahhasha/Tracker-bot/internal/contracts"
	"github.com/Ahhasha/Tracker-bot/internal/model"
)

type ExpenseRepository interface {
	Create(ctx context.Context, db contracts.DBTX, expense model.Expense) (int64, error)
	GetToday(ctx context.Context, db contracts.DBTX, userID int64) ([]model.Expense, error)
	GetWeek(ctx context.Context, db contracts.DBTX, userID int64) ([]model.Expense, error)
	GetMonth(ctx context.Context, db contracts.DBTX, userID int64) ([]model.Expense, error)
}

type CategoryRepository interface {
	GetByName(ctx context.Context, db contracts.DBTX, userID int64, name string) (model.Category, error)
}

type UserRepository interface {
	GetIDByTgID(ctx context.Context, db contracts.DBTX, tgUserID int64) (int64, error)
}
