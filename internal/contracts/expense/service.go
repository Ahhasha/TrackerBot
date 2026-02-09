package expense

import (
	"context"

	"github.com/Ahhasha/Tracker-bot/internal/model"
)

type ExpenseRequest struct {
	Amount      int64
	Category    string
	Description string
}

type Service interface {
	AddExpense(ctx context.Context, userID int64, req ExpenseRequest) (int64, error)
	Today(ctx context.Context, tgUserID int64) (model.PeriodReport, error)
	Week(ctx context.Context, tgUserID int64) (model.PeriodReport, error)
	Month(ctx context.Context, tgUserID int64) (model.PeriodReport, error)
}
