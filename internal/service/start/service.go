package start

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Ahhasha/Tracker-bot/internal/contracts"
	repoStart "github.com/Ahhasha/Tracker-bot/internal/repository/start"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
	repo   contracts.RegistrationRepo
}

func NewService(pool *pgxpool.Pool, logger *slog.Logger, repo contracts.RegistrationRepo) *Service {
	return &Service{
		pool:   pool,
		logger: logger,
		repo:   repo,
	}
}

func (s *Service) RegIfNotExist(ctx context.Context, tgID int64, username string) error {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	repoWithTx := repoStart.NewRepo(tx)

	exists, err := repoWithTx.UserExists(ctx, tgID)
	if err != nil {
		return fmt.Errorf("check user exists: %w", err)
	}

	if exists {
		s.logger.Info("user already exists", slog.Int64("tg_id", tgID))
		return nil
	}

	userID, err := repoWithTx.Create(ctx, tgID, username)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	if err := repoWithTx.CreateDefaultCategory(ctx, userID); err != nil {
		return fmt.Errorf("create default categories: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	s.logger.Info("user registered", slog.Int64("tg_id", tgID), slog.Int64("user_id", userID))
	return nil
}
