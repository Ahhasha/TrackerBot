package add

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Ahhasha/Tracker-bot/internal/contracts"
	"github.com/Ahhasha/Tracker-bot/internal/contracts/add"
	"github.com/Ahhasha/Tracker-bot/internal/model"
	"github.com/jackc/pgx/v5"
)

type service struct {
	tx       contracts.TxManager
	expRepo  add.ExpenseRepository
	catRepo  add.CategoryRepository
	logger   *slog.Logger
	userRepo add.UserRepository
}

func NewService(tx contracts.TxManager, expRepo add.ExpenseRepository, catRepo add.CategoryRepository, logger *slog.Logger, userRepo add.UserRepository) add.AddService {
	return &service{
		tx:       tx,
		expRepo:  expRepo,
		catRepo:  catRepo,
		logger:   logger,
		userRepo: userRepo,
	}
}

func (s *service) AddExpense(ctx context.Context, tgUserID int64, req add.AddRequest) (int64, error) {
	const op = "service.add.AddExpense"

	if strings.TrimSpace(req.Category) == "" {
		return 0, model.ErrInvalidCategory
	}

	var expenseID int64

	err := s.tx.Do(ctx, func(db contracts.DBTX) error {
		internalUserID, err := s.userRepo.GetIDByTgID(ctx, db, tgUserID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return model.ErrUserNotRegistered
			}
			return fmt.Errorf("%s: get internal user id: %w", op, err)
		}

		category, err := s.catRepo.GetByName(ctx, db, internalUserID, req.Category)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return model.ErrCategoryNotFound
			}
			return fmt.Errorf("%s: get category: %w", op, err)
		}

		expense := model.Expense{
			UserID:      internalUserID,
			Amount:      req.Amount,
			CategoryID:  category.ID,
			Description: req.Description,
		}

		if err := expense.Validate(); err != nil {
			return err
		}

		expenseID, err = s.expRepo.Create(ctx, db, expense)
		if err != nil {
			return fmt.Errorf("%s: create expense: %w", op, err)
		}

		return nil
	})

	if err != nil {
		return 0, err
	}
	return expenseID, nil
}
