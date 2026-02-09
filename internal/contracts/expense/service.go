package expense

import (
	"context"
)

type ExpenseRequest struct {
	Amount      int64
	Category    string
	Description string
}

type ExpenseService interface {
	AddExpense(ctx context.Context, userID int64, req ExpenseRequest) (int64, error)
}
