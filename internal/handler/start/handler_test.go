package start_test

import (
	"context"
	"io"
	"log/slog"
	"testing"

	startcont "github.com/Ahhasha/Tracker-bot/internal/contracts/start"
	startHandler "github.com/Ahhasha/Tracker-bot/internal/handler/start"
	"github.com/Ahhasha/Tracker-bot/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockregService struct {
	registerFunc func(ctx context.Context, tgID int64, username string) (startcont.RegisterResult, error)
}

func (m *mockregService) Register(ctx context.Context, tgID int64, username string) (startcont.RegisterResult, error) {
	if m.registerFunc != nil {
		return m.registerFunc(ctx, tgID, username)
	}
	return startcont.RegisterResult{}, nil
}

func TestStartHandler_Handle_NewUser(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	mockService := &mockregService{
		registerFunc: func(ctx context.Context, tgID int64, username string) (startcont.RegisterResult, error) {
			assert.Equal(t, int64(123), tgID, "tgID должен быть 123")

			assert.Equal(t, "test_user", username, "username должен быть test_user")

			return startcont.RegisterResult{
				UserID:            42,
				Created:           true,
				CategoriesCreated: []string{"Еда", "Транспорт", "Прочее", "Развлечения"},
			}, nil
		},
	}

	handler := startHandler.New(logger, mockService)

	cmd := &model.Command{
		Name:            model.CommandStart,
		ChatID:          456,
		UserID:          123,
		UserDisplayName: "test_user",
		RawArgs:         "",
	}

	result, err := handler.Handle(ctx, cmd)

	require.NoError(t, err, "Хендлер не должен возвращать ошибку")

	assert.Equal(t, int64(456), result.ChatID, "Ответ должен быть отправлен в тот же чат")
	assert.Contains(t, result.Text, "Добро пожаловать", "Должно быть приветствие")
	assert.Contains(t, result.Text, "Вы зарегистрированы", "Должно быть подтверждение")
	assert.Contains(t, result.Text, "категории", "Должен быть список категорий")
	assert.Contains(t, result.Text, "Еда", "Должна быть категория Еда")
	assert.Contains(t, result.Text, "Транспорт", "Должна быть категория Транспорт")
	assert.Contains(t, result.Text, "Развлечения", "Должна быть категория Развлечения")
	assert.Contains(t, result.Text, "Прочее", "Должна быть категория Прочее")
}

func TestStartHandler_Handle_ExistUser(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	mockService := &mockregService{
		registerFunc: func(ctx context.Context, tgID int64, username string) (startcont.RegisterResult, error) {
			assert.Equal(t, int64(456), tgID, "tgID должен быть 456")
			assert.Equal(t, "existing_user", username, "username должен быть existing_user")

			return startcont.RegisterResult{
				UserID:            42,
				Created:           false,
				CategoriesCreated: nil,
			}, nil
		},
	}

	handler := startHandler.New(logger, mockService)

	cmd := &model.Command{
		Name:            model.CommandStart,
		ChatID:          789,
		UserID:          456,
		UserDisplayName: "existing_user",
	}

	result, err := handler.Handle(ctx, cmd)

	require.NoError(t, err)
	assert.Equal(t, int64(789), result.ChatID)
	assert.Contains(t, result.Text, "Вы уже зарегистрированы!")
	assert.Contains(t, result.Text, "/help")
	assert.NotContains(t, result.Text, "категории", "Для существующего пользователя не должно быть списка категорий")
}

func TestStartHandler_Handle_ServiceError(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	mockService := &mockregService{
		registerFunc: func(ctx context.Context, tgID int64, username string) (startcont.RegisterResult, error) {
			return startcont.RegisterResult{}, assert.AnError
		},
	}

	handler := startHandler.New(logger, mockService)

	cmd := &model.Command{
		Name:            model.CommandStart,
		ChatID:          999,
		UserID:          777,
		UserDisplayName: "error_user",
	}

	result, err := handler.Handle(ctx, cmd)

	require.NoError(t, err, "Хендлер должен обработать ошибку сервиса, а не прокидывать её")

	assert.Equal(t, int64(999), result.ChatID)
	assert.Contains(t, result.Text, "Произошла ошибка", "Должно быть сообщение об ошибке")
}

func TestStartHandler_UserWithoutUsername(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	mockService := &mockregService{
		registerFunc: func(ctx context.Context, tgID int64, username string) (startcont.RegisterResult, error) {
			assert.Equal(t, int64(777), tgID, "tgID должен быть 777")
			assert.Equal(t, "", username, "username должен быть ПУСТОЙ строкой")

			return startcont.RegisterResult{
				UserID:            99,
				Created:           true,
				CategoriesCreated: []string{"Еда", "Транспорт", "Прочее", "Развлечения"},
			}, nil
		},
	}

	handler := startHandler.New(logger, mockService)

	cmd := &model.Command{
		Name:            model.CommandStart,
		ChatID:          888,
		UserID:          777,
		UserDisplayName: "",
		RawArgs:         "",
	}

	result, err := handler.Handle(ctx, cmd)

	require.NoError(t, err, "Хендлер должен принимать пользователя без username")
	assert.Equal(t, int64(888), result.ChatID, "Ответ в правильный чат")
	assert.Contains(t, result.Text, "Добро пожаловать", "Должно быть приветствие")
	assert.Contains(t, result.Text, "Вы зарегистрированы", "Подтверждение регистрации")
	assert.Contains(t, result.Text, "категории", "Список категорий")

}
