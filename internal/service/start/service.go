package start

import (
	"context"
	"fmt"

	startrepo "github.com/Ahhasha/Tracker-bot/internal/repository/start"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	pool *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

func (s *Service) RegIfNotExist(ctx context.Context, tgID int64, username string) error {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback(ctx)

	repo := startrepo.NewStart(tx)
	if _, err := repo.FindOrCreateStart(ctx, tgID, username); err != nil {
		return fmt.Errorf("find or create user and categories: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}
