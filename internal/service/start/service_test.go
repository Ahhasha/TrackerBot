package start

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/Ahhasha/Tracker-bot/internal/contracts"
	"github.com/stretchr/testify/require"
)

type fakeTxManager struct {
	calls int
	err   error
}

type fakeRepo struct {
	upsertFn func(ctx context.Context, db contracts.DBTX, tgID int64, username string) (int64, bool, error)
	catsFn   func(ctx context.Context, db contracts.DBTX, userID int64) ([]string, error)

	upsertCalls int
	catsCalls   int
}

func (f *fakeRepo) UpsertUser(ctx context.Context, db contracts.DBTX, tgID int64, username string) (int64, bool, error) {
	f.upsertCalls++
	return f.upsertFn(ctx, db, tgID, username)
}

func (f *fakeRepo) CreateDefaultCategories(ctx context.Context, db contracts.DBTX, userID int64) ([]string, error) {
	f.catsCalls++
	return f.catsFn(ctx, db, userID)
}

func (f *fakeTxManager) Do(ctx context.Context, fn func(db contracts.DBTX) error) error {
	f.calls++

	if f.err != nil {
		return f.err
	}
	return fn(nil)
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))
}

func TestService_Register_CreatedUser_CreatesDefaultCategories(t *testing.T) {
	ctx := context.Background()

	tx := &fakeTxManager{}

	repo := &fakeRepo{
		upsertFn: func(ctx context.Context, db contracts.DBTX, tgID int64, username string) (int64, bool, error) {
			return 42, true, nil
		},
		catsFn: func(ctx context.Context, db contracts.DBTX, userID int64) ([]string, error) {
			return []string{"Еда", "Транспорт", "Развлечения", "Прочее"}, nil
		},
	}

	svc := NewService(tx, repo, testLogger())

	res, err := svc.Register(ctx, 1001, "alice")
	require.NoError(t, err)

	require.True(t, res.Created)
	require.Equal(t, int64(42), res.UserID)

	require.Len(t, res.CategoriesCreated, 4)

	require.Equal(t, 1, tx.calls)
	require.Equal(t, 1, repo.upsertCalls)
	require.Equal(t, 1, repo.catsCalls)
}

func TestService_Register_ExistingUser_DoesNotCreateDefaultCategories(t *testing.T) {
	ctx := context.Background()

	tx := &fakeTxManager{}

	repo := &fakeRepo{
		upsertFn: func(ctx context.Context, db contracts.DBTX, tgID int64, username string) (int64, bool, error) {
			return 42, false, nil
		},
		catsFn: func(ctx context.Context, db contracts.DBTX, userID int64) ([]string, error) {
			t.Fatalf("CreateDefaultCategories NOT CALLED")
			return nil, nil
		},
	}

	svc := NewService(tx, repo, testLogger())

	res, err := svc.Register(ctx, 1001, "alice")
	require.NoError(t, err)

	require.False(t, res.Created)
	require.Equal(t, int64(42), res.UserID)

	require.Len(t, res.CategoriesCreated, 0)

	require.Equal(t, 1, tx.calls)
	require.Equal(t, 1, repo.upsertCalls)
	require.Equal(t, 0, repo.catsCalls)
}

func TestService_Register_UpsertFails_ReturnsError(t *testing.T) {
	ctx := context.Background()
	expErr := errors.New("db down")
	tx := &fakeTxManager{}

	repo := &fakeRepo{
		upsertFn: func(ctx context.Context, db contracts.DBTX, tgID int64, username string) (int64, bool, error) {
			return 0, false, expErr
		},
		catsFn: func(ctx context.Context, db contracts.DBTX, userID int64) ([]string, error) {
			t.Fatalf("CreateDefaultCategories NOT CALLED")
			return nil, nil
		},
	}

	svc := NewService(tx, repo, testLogger())

	_, err := svc.Register(ctx, 1001, "alice")
	require.Error(t, err)
	require.True(t, errors.Is(err, expErr))

	require.Equal(t, 1, tx.calls)
	require.Equal(t, 1, repo.upsertCalls)
	require.Equal(t, 0, repo.catsCalls)
}

func TestService_Register_TxManagerFails_ReturnsError(t *testing.T) {
	ctx := context.Background()
	expErr := errors.New("cannot begin tx")

	tx := &fakeTxManager{
		err: expErr,
	}

	repo := &fakeRepo{
		upsertFn: func(ctx context.Context, db contracts.DBTX, tgID int64, username string) (int64, bool, error) {
			t.Fatalf("UpsertUser NOT called")
			return 0, false, nil
		},
		catsFn: func(ctx context.Context, db contracts.DBTX, userID int64) ([]string, error) {
			t.Fatalf("CreateDefaultCategories NOT called")
			return nil, nil
		},
	}

	svc := NewService(tx, repo, testLogger())

	_, err := svc.Register(ctx, 1001, "alice")
	require.Error(t, err)
	require.True(t, errors.Is(err, expErr))

	require.Equal(t, 1, tx.calls)
	require.Equal(t, 0, repo.upsertCalls)
	require.Equal(t, 0, repo.catsCalls)
}
