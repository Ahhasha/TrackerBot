package add

import (
	"context"
)

type AddRequest struct {
	Amount      int64
	Category    string
	Description string
}

type AddService interface {
	AddExpense(ctx context.Context, userID int64, req AddRequest) (int64, error)
}
