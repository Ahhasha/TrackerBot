package expense_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Ahhasha/Tracker-bot/internal/model"
	expRepo "github.com/Ahhasha/Tracker-bot/internal/repository/expense"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func TestRepo_ExpenseCreate_HappyPath(t *testing.T) {
	ctx := context.Background()

	db, err := pgxmock.NewConn()
	require.NoError(t, err)
	defer db.Close(ctx)

	repo := expRepo.NewPostgresRepo()

	in := model.Expense{
		UserID:      10,
		Amount:      500,
		CategoryID:  3,
		Description: "обед",
	}

	q := `INSERT INTO expenses .* RETURNING id`

	db.ExpectQuery(q).
		WithArgs(in.UserID, in.Amount, in.CategoryID, in.Description).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(int64(42)))

	id, err := repo.Create(ctx, db, in)
	require.NoError(t, err)
	require.Equal(t, int64(42), id)

	require.NoError(t, db.ExpectationsWereMet())
}

func TestRepo_ExpenseCreate_ScanError(t *testing.T) {
	ctx := context.Background()

	db, err := pgxmock.NewConn()
	require.NoError(t, err)
	defer db.Close(ctx)

	repo := expRepo.NewPostgresRepo()

	in := model.Expense{UserID: 1, Amount: 1, CategoryID: 1, Description: "x"}

	q := `INSERT INTO expenses .* RETURNING id`

	dbErr := errors.New("db cancelled")

	db.ExpectQuery(q).
		WithArgs(in.UserID, in.Amount, in.CategoryID, in.Description).
		WillReturnError(dbErr)

	id, err := repo.Create(ctx, db, in)
	require.Error(t, err)
	require.Equal(t, int64(0), id)
	require.ErrorIs(t, err, dbErr)

	require.NoError(t, db.ExpectationsWereMet())
}

func TestRepo_GetPeriodWithCategory_HappyPath(t *testing.T) {
	ctx := context.Background()

	db, err := pgxmock.NewConn()
	require.NoError(t, err)
	defer db.Close(ctx)

	repo := expRepo.NewPostgresRepo()

	userID := int64(7)
	start := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 2, 2, 0, 0, 0, 0, time.UTC)

	q := `SELECT .* FROM expenses .* JOIN categories .* ORDER BY .*`

	rows := pgxmock.NewRows([]string{"amount", "description", "name", "created_at"}).
		AddRow(int64(500), "кофе", "Еда", time.Date(2026, 2, 1, 12, 0, 0, 0, time.UTC)).
		AddRow(int64(300), "метро", "Транспорт", time.Date(2026, 2, 1, 9, 0, 0, 0, time.UTC))

	db.ExpectQuery(q).
		WithArgs(userID, start, end).
		WillReturnRows(rows)

	got, err := repo.GetPeriodWithCategory(ctx, db, userID, start, end)
	require.NoError(t, err)
	require.Len(t, got, 2)

	require.Equal(t, int64(500), got[0].Amount)
	require.Equal(t, "Еда", got[0].Category)

	require.NoError(t, db.ExpectationsWereMet())
}

func TestRepo_GetPeriodWithCategory_QueryError(t *testing.T) {
	ctx := context.Background()

	db, err := pgxmock.NewConn()
	require.NoError(t, err)
	defer db.Close(ctx)

	repo := expRepo.NewPostgresRepo()
	dbErr := errors.New("query failed")
	q := `SELECT .* FROM expenses .* JOIN categories .*`

	db.ExpectQuery(q).
		WithArgs(int64(7), pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnError(dbErr)

	got, err := repo.GetPeriodWithCategory(ctx, db, 7, time.Now(), time.Now().Add(time.Hour))
	require.Error(t, err)
	require.Nil(t, got)
	require.ErrorIs(t, err, dbErr)

	require.NoError(t, db.ExpectationsWereMet())
}

func TestRepo_GetPeriodWithCategory_ScanError(t *testing.T) {
	ctx := context.Background()

	db, err := pgxmock.NewConn()
	require.NoError(t, err)
	defer db.Close(ctx)

	repo := expRepo.NewPostgresRepo()

	userID := int64(7)
	start := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 2, 2, 0, 0, 0, 0, time.UTC)

	q := `SELECT .* FROM expenses .* JOIN categories .*`

	dbErr := errors.New("scan failed")

	rows := pgxmock.NewRows([]string{"amount", "description", "name", "created_at"}).
		AddRow(int64(500), "кофе", "Еда", time.Date(2026, 2, 1, 12, 0, 0, 0, time.UTC)).
		RowError(0, dbErr)

	db.ExpectQuery(q).
		WithArgs(userID, start, end).
		WillReturnRows(rows)

	got, err := repo.GetPeriodWithCategory(ctx, db, userID, start, end)

	require.Error(t, err)
	require.Nil(t, got)
	require.ErrorIs(t, err, dbErr)

	require.NoError(t, db.ExpectationsWereMet())
}

func TestRepo_GetPeriodWithCategory_RowsErr(t *testing.T) {
	ctx := context.Background()

	db, err := pgxmock.NewConn()
	require.NoError(t, err)
	defer db.Close(ctx)

	repo := expRepo.NewPostgresRepo()

	userID := int64(7)
	start := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 2, 2, 0, 0, 0, 0, time.UTC)

	q := `SELECT .* FROM expenses .* JOIN categories .*`

	dbErr := errors.New("rows iteration failed")

	rows := pgxmock.NewRows([]string{"amount", "description", "name", "created_at"}).
		AddRow(int64(500), "кофе", "Еда", time.Date(2026, 2, 1, 12, 0, 0, 0, time.UTC)).
		CloseError(dbErr)

	db.ExpectQuery(q).
		WithArgs(userID, start, end).
		WillReturnRows(rows)

	got, err := repo.GetPeriodWithCategory(ctx, db, userID, start, end)

	require.Error(t, err)
	require.Nil(t, got)
	require.ErrorIs(t, err, dbErr)

	require.NoError(t, db.ExpectationsWereMet())
}
