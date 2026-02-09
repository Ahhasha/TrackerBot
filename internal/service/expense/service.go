package expense

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

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
	now      func() time.Time
}

func NewService(tx contracts.TxManager, expRepo expense.ExpenseRepository, catRepo expense.CategoryRepository, userRepo expense.UserRepository, now func() time.Time) expense.Service {
	return &service{
		tx:       tx,
		expRepo:  expRepo,
		catRepo:  catRepo,
		userRepo: userRepo,
		now:      now,
	}
}

type period string

const (
	periodToday period = "today"
	periodWeek  period = "week"
	periodMonth period = "month"
)

func (s *service) AddExpense(ctx context.Context, tgUserID int64, req expense.ExpenseRequest) (int64, error) {
	const op = "service.expense.AddExpense"

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

func (s *service) Today(ctx context.Context, tgUserID int64) (model.PeriodReport, error) {
	return s.period(ctx, tgUserID, periodToday)
}

func (s *service) Week(ctx context.Context, tgUserID int64) (model.PeriodReport, error) {
	return s.period(ctx, tgUserID, periodWeek)
}

func (s *service) Month(ctx context.Context, tgUserID int64) (model.PeriodReport, error) {
	return s.period(ctx, tgUserID, periodMonth)
}

func (s *service) period(ctx context.Context, tgUserID int64, p period) (model.PeriodReport, error) {
	const op = "service.expense.period"
	var out model.PeriodReport

	err := s.tx.Do(ctx, func(db contracts.DBTX) error {
		internalUserID, err := s.userRepo.GetIDByTgID(ctx, db, tgUserID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return model.ErrUserNotRegistered
			}
			return fmt.Errorf("%s: get internal user id: %w", op, err)
		}

		var rows []model.ExpenseWithCategory
		switch p {
		case periodToday:
			rows, err = s.expRepo.GetTodayWithCategory(ctx, db, internalUserID)
		case periodWeek:
			rows, err = s.expRepo.GetWeekWithCategory(ctx, db, internalUserID)
		case periodMonth:
			rows, err = s.expRepo.GetMonthWithCategory(ctx, db, internalUserID)
		default:
			rows, err = s.expRepo.GetTodayWithCategory(ctx, db, internalUserID)
		}

		if err != nil {
			return fmt.Errorf("%s: get expenses: %w", op, err)
		}

		now := s.now()
		anchor, _ := periodRange(p, now)
		out = buildPeriodReport(anchor, rows)

		return nil
	})
	if err != nil {
		return model.PeriodReport{}, err
	}

	return out, nil
}

func buildPeriodReport(anchor time.Time, rows []model.ExpenseWithCategory) model.PeriodReport {

	byCat := make(map[string]*model.CategoryReport)
	var total int64

	for _, r := range rows {
		total += r.Amount

		cr, ok := byCat[r.Category]
		if !ok {
			cr = &model.CategoryReport{Name: r.Category}
			byCat[r.Category] = cr
		}

		cr.Total += r.Amount

		cr.Items = append(cr.Items, model.ExpenseItem{
			Amount:      r.Amount,
			Description: r.Description,
		})
	}

	cats := make([]model.CategoryReport, 0, len(byCat))
	for _, v := range byCat {
		cats = append(cats, *v)
	}

	sort.Slice(cats, func(i, j int) bool {
		if cats[i].Total == cats[j].Total {
			return cats[i].Name < cats[j].Name
		}
		return cats[i].Total > cats[j].Total
	})

	return model.PeriodReport{
		Date:       anchor,
		Categories: cats,
		Total:      total,
	}
}

func periodRange(p period, now time.Time) (start, end time.Time) {
	switch p {
	case periodToday:
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 0, 1)
		return start, end

	case periodWeek:
		weekday := now.Weekday()
		if weekday == time.Sunday {
			weekday = 7
		}
		daysSinceMonday := int(weekday) - 1
		start = time.Date(now.Year(), now.Month(), now.Day()-daysSinceMonday, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 0, 7)
		return start, end

	case periodMonth:
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 1, 0)
		return start, end

	default:
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 0, 1)
		return start, end
	}
}
