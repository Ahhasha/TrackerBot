package expense_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Ahhasha/Tracker-bot/internal/contracts"
	expCon "github.com/Ahhasha/Tracker-bot/internal/contracts/expense"
	"github.com/Ahhasha/Tracker-bot/internal/contracts/mocks"
	"github.com/Ahhasha/Tracker-bot/internal/model"
	"github.com/Ahhasha/Tracker-bot/internal/service/expense"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestService_AddExpense_HappyPath(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mocks.NewMockTxManager(ctrl)
	db := mocks.NewMockDBTX(ctrl)

	expRepo := mocks.NewMockExpenseRepository(ctrl)
	catRepo := mocks.NewMockCategoryRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)

	now := func() time.Time {
		return time.Date(2026, 2, 13, 10, 0, 0, 0, time.UTC)
	}

	var _ expCon.Service = expense.NewService(tx, expRepo, catRepo, userRepo, now)
	svc := expense.NewService(tx, expRepo, catRepo, userRepo, now)

	tgUserID := int64(777)
	internalUserID := int64(10)
	categoryID := int64(3)
	expectedExpenseID := int64(42)

	req := expCon.ExpenseRequest{
		Amount:      500,
		Category:    "Еда",
		Description: "обед",
	}
	tx.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(contracts.DBTX) error) error {
		return fn(db)
	})

	userRepo.EXPECT().GetIDByTgID(gomock.Any(), db, tgUserID).Return(internalUserID, nil)

	catRepo.EXPECT().GetByName(gomock.Any(), db, internalUserID, req.Category).Return(model.Category{ID: categoryID, Name: req.Category}, nil)

	expRepo.EXPECT().Create(gomock.Any(), db, gomock.Any()).DoAndReturn(func(_ context.Context, _ contracts.DBTX, e model.Expense) (int64, error) {
		require.Equal(t, internalUserID, e.UserID)
		require.Equal(t, req.Amount, e.Amount)
		require.Equal(t, categoryID, e.CategoryID)
		require.Equal(t, req.Description, e.Description)
		return expectedExpenseID, nil
	})

	gotID, err := svc.AddExpense(ctx, tgUserID, req)

	require.NoError(t, err)
	require.Equal(t, expectedExpenseID, gotID)
}

func TestService_AddExpense_InvalidCategory(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mocks.NewMockTxManager(ctrl)
	expRepo := mocks.NewMockExpenseRepository(ctrl)
	catRepo := mocks.NewMockCategoryRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)

	now := func() time.Time {
		return time.Date(2026, 2, 13, 10, 0, 0, 0, time.UTC)
	}
	svc := expense.NewService(tx, expRepo, catRepo, userRepo, now)

	tx.EXPECT().Do(gomock.Any(), gomock.Any()).Times(0)

	_, err := svc.AddExpense(ctx, 777, expCon.ExpenseRequest{
		Amount:      500,
		Category:    "   ",
		Description: "обед",
	})

	require.ErrorIs(t, err, model.ErrInvalidCategory)
}

func TestService_AddExpense_UserNotRegistered(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mocks.NewMockTxManager(ctrl)
	db := mocks.NewMockDBTX(ctrl)

	expRepo := mocks.NewMockExpenseRepository(ctrl)
	catRepo := mocks.NewMockCategoryRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)

	now := func() time.Time {
		return time.Date(2026, 2, 13, 10, 0, 0, 0, time.UTC)
	}
	svc := expense.NewService(tx, expRepo, catRepo, userRepo, now)

	tgUserID := int64(777)
	req := expCon.ExpenseRequest{
		Amount:      500,
		Category:    "Еда",
		Description: "обед",
	}

	tx.EXPECT().
		Do(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, fn func(contracts.DBTX) error) error {
		return fn(db)
	})

	userRepo.EXPECT().GetIDByTgID(gomock.Any(), db, tgUserID).Return(int64(0), pgx.ErrNoRows)

	gotID, err := svc.AddExpense(ctx, tgUserID, req)

	require.Error(t, err)
	require.ErrorIs(t, err, model.ErrUserNotRegistered)
	require.Equal(t, int64(0), gotID)
}

func TestService_AddExpense_UserRepoError_Wrapped(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mocks.NewMockTxManager(ctrl)
	db := mocks.NewMockDBTX(ctrl)

	expRepo := mocks.NewMockExpenseRepository(ctrl)
	catRepo := mocks.NewMockCategoryRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)

	now := func() time.Time {
		return time.Date(2026, 2, 13, 10, 0, 0, 0, time.UTC)
	}
	svc := expense.NewService(tx, expRepo, catRepo, userRepo, now)

	tgUserID := int64(777)
	req := expCon.ExpenseRequest{
		Amount:      500,
		Category:    "Еда",
		Description: "обед",
	}

	tx.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, fn func(contracts.DBTX) error) error {
		return fn(db)
	})

	dbErr := errors.New("db is down")
	userRepo.EXPECT().GetIDByTgID(gomock.Any(), db, tgUserID).Return(int64(0), dbErr)

	catRepo.EXPECT().GetByName(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	expRepo.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	gotID, err := svc.AddExpense(ctx, tgUserID, req)

	require.Error(t, err)
	require.Equal(t, int64(0), gotID)

	require.ErrorIs(t, err, dbErr)

	require.Contains(t, err.Error(), "service.expense.AddExpense")
	require.Contains(t, err.Error(), "get internal user id")
}

func TestService_AddExpense_CategoryNotFound(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mocks.NewMockTxManager(ctrl)
	db := mocks.NewMockDBTX(ctrl)

	expRepo := mocks.NewMockExpenseRepository(ctrl)
	catRepo := mocks.NewMockCategoryRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)

	now := func() time.Time {
		return time.Date(2026, 2, 13, 10, 0, 0, 0, time.UTC)
	}
	svc := expense.NewService(tx, expRepo, catRepo, userRepo, now)

	tgUserID := int64(777)
	internalUserID := int64(10)

	req := expCon.ExpenseRequest{
		Amount:      500,
		Category:    "Несуществующая",
		Description: "обед",
	}

	tx.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, fn func(contracts.DBTX) error) error {
		return fn(db)
	})

	userRepo.EXPECT().GetIDByTgID(gomock.Any(), db, tgUserID).Return(internalUserID, nil)

	catRepo.EXPECT().GetByName(gomock.Any(), db, internalUserID, req.Category).Return(model.Category{}, pgx.ErrNoRows)

	expRepo.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	gotID, err := svc.AddExpense(ctx, tgUserID, req)

	require.ErrorIs(t, err, model.ErrCategoryNotFound)
	require.Equal(t, int64(0), gotID)
}

func TestService_AddExpense_CategoryRepoError_Wrapped(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mocks.NewMockTxManager(ctrl)
	db := mocks.NewMockDBTX(ctrl)

	expRepo := mocks.NewMockExpenseRepository(ctrl)
	catRepo := mocks.NewMockCategoryRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)

	now := func() time.Time {
		return time.Date(2026, 2, 13, 10, 0, 0, 0, time.UTC)
	}
	svc := expense.NewService(tx, expRepo, catRepo, userRepo, now)

	tgUserID := int64(777)
	internalUserID := int64(10)

	req := expCon.ExpenseRequest{
		Amount:      500,
		Category:    "Еда",
		Description: "обед",
	}

	tx.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, fn func(contracts.DBTX) error) error {
		return fn(db)
	})

	userRepo.EXPECT().GetIDByTgID(gomock.Any(), db, tgUserID).Return(internalUserID, nil)

	dbErr := errors.New("categories query failed")
	catRepo.EXPECT().GetByName(gomock.Any(), db, internalUserID, req.Category).Return(model.Category{}, dbErr)

	expRepo.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	gotID, err := svc.AddExpense(ctx, tgUserID, req)

	require.Error(t, err)
	require.Equal(t, int64(0), gotID)
	require.ErrorIs(t, err, dbErr)
	require.Contains(t, err.Error(), "service.expense.AddExpense")
	require.Contains(t, err.Error(), "get category")
}

func TestService_AddExpense_ValidateError_ReturnIs(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mocks.NewMockTxManager(ctrl)
	db := mocks.NewMockDBTX(ctrl)

	expRepo := mocks.NewMockExpenseRepository(ctrl)
	catRepo := mocks.NewMockCategoryRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)

	now := func() time.Time {
		return time.Date(2026, 2, 13, 10, 0, 0, 0, time.UTC)
	}
	svc := expense.NewService(tx, expRepo, catRepo, userRepo, now)

	tgUserID := int64(777)
	internalUserID := int64(10)
	categoryID := int64(3)

	req := expCon.ExpenseRequest{
		Amount:      -1,
		Category:    "Еда",
		Description: "обед",
	}

	tx.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, fn func(contracts.DBTX) error) error {
		return fn(db)
	})

	userRepo.EXPECT().GetIDByTgID(gomock.Any(), db, tgUserID).Return(internalUserID, nil)

	catRepo.EXPECT().GetByName(gomock.Any(), db, internalUserID, req.Category).Return(model.Category{ID: categoryID, Name: req.Category}, nil)

	expRepo.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	gotID, err := svc.AddExpense(ctx, tgUserID, req)

	require.Error(t, err)
	require.Equal(t, int64(0), gotID)

	require.NotContains(t, err.Error(), "create expense")
}

func TestService_AddExpense_CreateError_Wrapped(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mocks.NewMockTxManager(ctrl)
	db := mocks.NewMockDBTX(ctrl)

	expRepo := mocks.NewMockExpenseRepository(ctrl)
	catRepo := mocks.NewMockCategoryRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)

	now := func() time.Time {
		return time.Date(2026, 2, 13, 10, 0, 0, 0, time.UTC)
	}
	svc := expense.NewService(tx, expRepo, catRepo, userRepo, now)

	tgUserID := int64(777)
	internalUserID := int64(10)
	categoryID := int64(3)

	req := expCon.ExpenseRequest{
		Amount:      500,
		Category:    "Еда",
		Description: "обед",
	}

	tx.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, fn func(contracts.DBTX) error) error {
		return fn(db)
	})

	userRepo.EXPECT().GetIDByTgID(gomock.Any(), db, tgUserID).Return(internalUserID, nil)

	catRepo.EXPECT().GetByName(gomock.Any(), db, internalUserID, req.Category).Return(model.Category{ID: categoryID, Name: req.Category}, nil)

	dbErr := errors.New("insert failed")
	expRepo.EXPECT().Create(gomock.Any(), db, gomock.Any()).Return(int64(0), dbErr)

	gotID, err := svc.AddExpense(ctx, tgUserID, req)

	require.Error(t, err)
	require.Equal(t, int64(0), gotID)
	require.ErrorIs(t, err, dbErr)
	require.Contains(t, err.Error(), "service.expense.AddExpense")
	require.Contains(t, err.Error(), "create expense")
}

func TestService_Today_HappyPath(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mocks.NewMockTxManager(ctrl)
	db := mocks.NewMockDBTX(ctrl)

	expRepo := mocks.NewMockExpenseRepository(ctrl)
	catRepo := mocks.NewMockCategoryRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)

	now := func() time.Time {
		return time.Date(2026, 2, 13, 10, 0, 0, 0, time.UTC)
	}

	svc := expense.NewService(tx, expRepo, catRepo, userRepo, now)

	tgUserID := int64(777)
	internalUserID := int64(10)

	expectedStart := time.Date(2026, 2, 13, 0, 0, 0, 0, time.UTC)
	expectedEnd := time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC)

	rows := []model.ExpenseWithCategory{
		{Amount: 500, Description: "кофе", Category: "Еда"},
		{Amount: 300, Description: "метро", Category: "Транспорт"},
	}

	tx.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, fn func(contracts.DBTX) error) error {
		return fn(db)
	})

	userRepo.EXPECT().GetIDByTgID(gomock.Any(), db, tgUserID).Return(internalUserID, nil)

	expRepo.EXPECT().GetPeriodWithCategory(gomock.Any(), db, internalUserID, expectedStart, expectedEnd).Return(rows, nil)

	report, err := svc.Today(ctx, tgUserID)

	require.NoError(t, err)

	require.Equal(t, int64(800), report.Total)
	require.Len(t, report.Categories, 2)

	require.Equal(t, "Еда", report.Categories[0].Name)
	require.Equal(t, int64(500), report.Categories[0].Total)

	require.Equal(t, "Транспорт", report.Categories[1].Name)
	require.Equal(t, int64(300), report.Categories[1].Total)
}

func TestService_Today_UserNotRegistered(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mocks.NewMockTxManager(ctrl)
	db := mocks.NewMockDBTX(ctrl)

	expRepo := mocks.NewMockExpenseRepository(ctrl)
	catRepo := mocks.NewMockCategoryRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)

	now := func() time.Time {
		return time.Now()
	}

	svc := expense.NewService(tx, expRepo, catRepo, userRepo, now)

	tx.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, fn func(contracts.DBTX) error) error {
		return fn(db)
	})

	userRepo.EXPECT().GetIDByTgID(gomock.Any(), db, int64(777)).Return(int64(0), pgx.ErrNoRows)

	expRepo.EXPECT().GetPeriodWithCategory(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	_, err := svc.Today(ctx, 777)

	require.ErrorIs(t, err, model.ErrUserNotRegistered)
}

func TestService_Today_GetExpensesError_Wrapped(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := mocks.NewMockTxManager(ctrl)
	db := mocks.NewMockDBTX(ctrl)

	expRepo := mocks.NewMockExpenseRepository(ctrl)
	catRepo := mocks.NewMockCategoryRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)

	now := func() time.Time {
		return time.Date(2026, 2, 13, 10, 0, 0, 0, time.UTC)
	}

	svc := expense.NewService(tx, expRepo, catRepo, userRepo, now)

	tgUserID := int64(777)
	internalUserID := int64(10)

	tx.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, fn func(contracts.DBTX) error) error {
		return fn(db)
	})

	userRepo.EXPECT().GetIDByTgID(gomock.Any(), db, tgUserID).Return(internalUserID, nil)

	dbErr := errors.New("select failed")

	expRepo.EXPECT().GetPeriodWithCategory(gomock.Any(), db, internalUserID, gomock.Any(), gomock.Any()).Return(nil, dbErr)

	_, err := svc.Today(ctx, tgUserID)

	require.Error(t, err)
	require.ErrorIs(t, err, dbErr)
	require.Contains(t, err.Error(), "service.expense.buildReport")
	require.Contains(t, err.Error(), "get expenses")
}
