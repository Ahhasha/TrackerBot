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

type service struct {
	tx      contracts.TxManager
	expRepo add.ExpenseRepository
	catRepo add.CategoryRepository
	logger  *slog.Logger
}

func NewService(tx contracts.TxManager, expRepo add.ExpenseRepository, catRepo add.CategoryRepository, logger *slog.Logger) add.AddService {
	return &service{
		tx:      tx,
		expRepo: expRepo,
		catRepo: catRepo,
		logger:  logger,
	}
}

func (s *service) AddExpense(ctx context.Context, tgUserID int64, req add.AddRequest) (int64, error) {
	const op = "service.add.AddExpense"

	if strings.TrimSpace(req.Category) == "" {
		return 0, model.ErrInvalidCategory
	}

	var expenseID int64

	err := s.tx.Do(ctx, func(db contracts.DBTX) error {
		const getUserid = `SELECT id FROM users WHERE tg_id = $1`
		var internalUserID int64
		err := db.QueryRow(ctx, getUserid, tgUserID).Scan(&internalUserID)
		if err != nil {
			return fmt.Errorf("user with tg_id %d not found. Use /start", tgUserID)
		}

		category, err := s.catRepo.GetByName(ctx, db, internalUserID, req.Category)
		if err != nil {
			return fmt.Errorf("%s: get category: %w", op, err)
		}

		expense := model.Expense{
			UserID:      internalUserID,
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
		s.logger.Error("fail to add expense", "op", op, "tguser_id", tgUserID, "error", err)
		return 0, err
	}
	s.logger.Info("expense add", "expense_id", expenseID, "tguser_id", tgUserID, "amount", req.Amount)

	return expenseID, nil
}
