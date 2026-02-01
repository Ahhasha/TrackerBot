package start_test

import (
	"context"
	"testing"

	repoStart "github.com/Ahhasha/Tracker-bot/internal/repository/start"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockDBTX struct {
	execFunc     func(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	queryFunc    func(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error)
	queryRowFunc func(ctx context.Context, sql string, arguments ...any) pgx.Row
}

func (m *mockDBTX) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	if m.execFunc != nil {
		return m.execFunc(ctx, sql, arguments...)
	}
	return pgconn.CommandTag{}, nil
}

func (m *mockDBTX) Query(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error) {
	if m.queryFunc != nil {
		return m.queryFunc(ctx, sql, arguments...)
	}
	return nil, nil
}

func (m *mockDBTX) QueryRow(ctx context.Context, sql string, arguments ...any) pgx.Row {
	if m.queryRowFunc != nil {
		return m.queryRowFunc(ctx, sql, arguments...)
	}
	return nil
}

type mockRow struct {
	scanFunc func(dest ...any) error
}

func (m *mockRow) Scan(dest ...any) error {
	if m.scanFunc != nil {
		return m.scanFunc(dest...)
	}
	return nil
}

func TestRepo_UpsertUser_NewUser(t *testing.T) {
	ctx := context.Background()
	repo := repoStart.NewRepo()

	db := &mockDBTX{
		queryRowFunc: func(ctx context.Context, sql string, arguments ...any) pgx.Row {
			assert.Contains(t, sql, "INSERT INTO users", "Нужен INSERT запрос")
			assert.Contains(t, sql, "tg_id, username", "Должны вставляться tg_id и username")
			assert.Contains(t, sql, "VALUES ($1, $2)", "Должно быть 2 параметра")
			assert.Contains(t, sql, "ON CONFLICT (tg_id) DO NOTHING", "Должен быть ON CONFLICT")
			assert.Contains(t, sql, "RETURNING id", "Должен возвращать id")
			assert.Len(t, arguments, 2, "Должно быть 2 аргумента")
			assert.Equal(t, int64(123), arguments[0], "Первый аргумент должен быть tg_id=123")
			assert.Equal(t, "test_user", arguments[1], "Второй аргумент должен быть username=test_user")
			return &mockRow{
				scanFunc: func(dest ...any) error {
					idPtr, ok := dest[0].(*int64)
					require.True(t, ok, "Первый аргумент Scan должен быть *int64")
					*idPtr = 42
					return nil
				},
			}
		},
	}

	userID, created, err := repo.UpsertUser(ctx, db, 123, "test_user")
	require.NoError(t, err)
	assert.Equal(t, int64(42), userID, "Должен вернуть ID=42")
	assert.True(t, created, "Должен вернуть created=true для нового пользователя")
}

func TestRepo_UpsertUser_ExistingUser(t *testing.T) {
	ctx := context.Background()
	repo := repoStart.NewRepo()

	callCount := 0
	db := &mockDBTX{
		queryRowFunc: func(ctx context.Context, sql string, arguments ...any) pgx.Row {
			callCount++

			if callCount == 1 {
				assert.Contains(t, sql, "INSERT INTO users", "Нужен INSERT запрос для поиска сущ. пз.")
				return &mockRow{
					scanFunc: func(dest ...any) error {
						return pgx.ErrNoRows
					},
				}
			}

			assert.Contains(t, sql, "UPDATE users", "Нужен INSERT запрос")
			assert.Contains(t, sql, "SET username = $2", "Второй аргумент должен быть username")
			assert.Contains(t, sql, "WHERE tg_id = $1", "Первый аргумент должен быть tg_id")
			assert.Contains(t, sql, "RETURNING id", "Должен возщвращать id")

			assert.Len(t, arguments, 2)
			assert.Equal(t, int64(123), arguments[0], "Первый аргумент должен быть tg_id=123")
			assert.Equal(t, "updated_user", arguments[1], "Второй аргумент должен быть username=updated_user")

			return &mockRow{
				scanFunc: func(dest ...any) error {
					idPtr := dest[0].(*int64)
					*idPtr = 42
					return nil
				},
			}
		},
	}

	userID, created, err := repo.UpsertUser(ctx, db, 123, "updated_user")

	require.NoError(t, err)
	assert.Equal(t, int64(42), userID)
	assert.False(t, created)
	assert.Equal(t, 2, callCount)
}

func TestRepo_UpsertUser_DatabaseError(t *testing.T) {
	ctx := context.Background()
	repo := repoStart.NewRepo()

	db := &mockDBTX{
		queryRowFunc: func(ctx context.Context, sql string, arguments ...any) pgx.Row {
			assert.Contains(t, sql, "INSERT INTO users")

			return &mockRow{
				scanFunc: func(dest ...any) error {
					return assert.AnError
				},
			}
		},
	}

	userID, created, err := repo.UpsertUser(ctx, db, 123, "test_user")

	require.Error(t, err)

	assert.ErrorIs(t, err, assert.AnError, "Ошибка должна прокинуться из репозитория")

	assert.Equal(t, int64(0), userID, "При ошибке должен вернуть 0")

	assert.False(t, created, "При ошибке должен вернуть false")
}

func TestRepo_UpsertUser_UpdateError(t *testing.T) {
	ctx := context.Background()
	repo := repoStart.NewRepo()

	callCount := 0
	db := &mockDBTX{
		queryRowFunc: func(ctx context.Context, sql string, arguments ...any) pgx.Row {
			callCount++

			if callCount == 1 {
				return &mockRow{
					scanFunc: func(dest ...any) error {
						return pgx.ErrNoRows
					},
				}
			}

			return &mockRow{
				scanFunc: func(dest ...any) error {
					return assert.AnError
				},
			}
		},
	}

	userID, created, err := repo.UpsertUser(ctx, db, 123, "test_user")

	require.Error(t, err)
	assert.Equal(t, 2, callCount)
	assert.Equal(t, int64(0), userID)
	assert.False(t, created)
}

func TestRepo_CreateDefaultCategories_Success(t *testing.T) {
	ctx := context.Background()
	repo := repoStart.NewRepo()

	execCalls := 0
	categories := []string{"Еда", "Транспорт", "Прочее", "Развлечения"}

	db := &mockDBTX{
		execFunc: func(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
			assert.Contains(t, sql, "INSERT INTO categories (user_id, name)")
			assert.Contains(t, sql, "VALUES ($1, $2)")
			assert.Contains(t, sql, "ON CONFLICT (user_id, name) DO NOTHING")
			assert.Len(t, arguments, 2, "Должно быть 2 аргумента")
			assert.Equal(t, int64(42), arguments[0], "user_id должен быть 42")
			expectedCategory := categories[execCalls]
			assert.Equal(t, expectedCategory, arguments[1], "Категория %d должна быть %s, а получил %s", execCalls, expectedCategory, arguments[1])
			execCalls++
			return pgconn.CommandTag{}, nil
		},
	}
	categories, err := repo.CreateDefaultCategories(ctx, db, 42)

	require.NoError(t, err, "Не должно быть ошибки")
	assert.Len(t, categories, 4, "Должно вернуть 4 категории")
	assert.Equal(t, categories, categories, "Должны вернуться те же категории")
	assert.Equal(t, 4, execCalls, "Должно быть 4 вызова Exec")
}

func TestRepo_CreateDefaultCategories_AlreadyExist(t *testing.T) {
	ctx := context.Background()
	repo := repoStart.NewRepo()

	execCalls := 0
	categories := []string{"Еда", "Транспорт", "Прочее", "Развлечения"}
	db := &mockDBTX{
		execFunc: func(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
			execCalls++

			assert.Contains(t, sql, "ON CONFLICT (user_id, name) DO NOTHING",
				"Должен игнорировать дубли")

			assert.Equal(t, int64(77), arguments[0], "user_id должен быть 77")

			assert.Equal(t, categories[execCalls-1], arguments[1],
				"Категория %d должна быть %s", execCalls, categories[execCalls-1])

			return pgconn.CommandTag{}, nil
		},
	}

	categories, err := repo.CreateDefaultCategories(ctx, db, 77)

	require.NoError(t, err, "ON CONFLICT DO NOTHING не должен возвращать ошибку при дублях")
	assert.Len(t, categories, 4, "Должно вернуть 4 категории если они существуют")
	assert.Equal(t, 4, execCalls, "Должно попытаться создать все 4 категории если они есть")
}

func TestRepo_CreateDefaultCategories_DatabaseError(t *testing.T) {
	ctx := context.Background()
	repo := repoStart.NewRepo()

	execCalls := 0
	db := &mockDBTX{
		execFunc: func(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
			execCalls++

			if execCalls <= 2 {
				assert.Contains(t, sql, "INSERT INTO categories")
				assert.Equal(t, int64(99), arguments[0])

				if execCalls == 1 {
					assert.Equal(t, "Еда", arguments[1], "Первая категория должна быть Еда")
				} else {
					assert.Equal(t, "Транспорт", arguments[1], "Вторая категория должна быть Транспорт")
				}

				return pgconn.CommandTag{}, nil
			}

			assert.Equal(t, "Прочее", arguments[1], "Третья категория должна быть Прочее")

			return pgconn.CommandTag{}, assert.AnError
		},
	}

	categories, err := repo.CreateDefaultCategories(ctx, db, 99)

	require.Error(t, err, "Должна вернуться ошибка БД")
	assert.ErrorContains(t, err, "Прочее", "Ошибка должна указывать какая категория не создалась")
	assert.Nil(t, categories, "При ошибке должен вернуть nil")
	assert.Equal(t, 3, execCalls, "Должно быть 3 вызова Exec до ошибки")
}
