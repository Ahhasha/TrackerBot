package add

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Ahhasha/Tracker-bot/internal/contracts"
	"github.com/Ahhasha/Tracker-bot/internal/contracts/add"
	"github.com/Ahhasha/Tracker-bot/internal/model"
)

type Service struct {
	tx      contracts.TxManager
	expRepo add.ExpenseRepository
	catRepo add.CategoryRepository
	logger  *slog.Logger
}

func NewService(tx contracts.TxManager, expRepo add.ExpenseRepository, catRepo add.CategoryRepository, logger *slog.Logger) *Service {
	return &Service{
		tx:      tx,
		expRepo: expRepo,
		catRepo: catRepo,
		logger:  logger,
	}
}

func (s *Service) AddExpense(ctx context.Context, userID int64, req add.AddRequest) (int64, error) {
	const op = "service.add.AddExpense"

	if req.Amount < 0 {
		return 0, model.ErrInvalidAmount
	}

	if strings.TrimSpace(req.Category) == "" {
		return 0, model.ErrInvalidCategory
	}

	var expenseID int64

	err := s.tx.Do(ctx, func(db contracts.DBTX) error {
		category, err := s.catRepo.GetByName(ctx, db, userID, req.Category)
		if err != nil {
			return fmt.Errorf("%s: get category: %w", op, err)
		}

		expense := model.Expense{
			UserID:      userID,
			Amount:      req.Amount,
			CategoryID:  category.ID,
			Description: req.Description,
		}

		if err := expense.Validate(); err != nil {
			return fmt.Errorf("%s: validate expense: %w", op, err)
		}

		expenseID, err = s.expRepo.Create(ctx, db, expense)
		if err != nil {
			return fmt.Errorf("%s: create expense: %w", op, err)
		}

		return nil
	})

	if err != nil {
		s.logger.Error("fail to add expense", "op", op, "user_id", userID, "error", err)
	}
	s.logger.Info("expense add", "expense_id", expenseID, "user_id", userID, "amount", req.Amount)

	return expenseID, nil
}
