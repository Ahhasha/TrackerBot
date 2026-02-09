package expense

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Ahhasha/Tracker-bot/internal/contracts"
	"github.com/Ahhasha/Tracker-bot/internal/contracts/expense"
	"github.com/Ahhasha/Tracker-bot/internal/model"
	"github.com/jackc/pgx/v5"
)

type service struct {
	tx       contracts.TxManager
	expRepo  expense.ExpenseRepository
	catRepo  expense.CategoryRepository
	userRepo expense.UserRepository
}

func NewService(tx contracts.TxManager, expRepo expense.ExpenseRepository, catRepo expense.CategoryRepository, userRepo expense.UserRepository) expense.ExpenseService {
	return &service{
		tx:       tx,
		expRepo:  expRepo,
		catRepo:  catRepo,
		userRepo: userRepo,
	}
}

func (s *service) AddExpense(ctx context.Context, tgUserID int64, req expense.ExpenseRequest) (int64, error) {
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

func (s *service) GetTodayWithCategory(ctx context.Context)
