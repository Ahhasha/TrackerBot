package database

import (
	"context"
	"fmt"

	"github.com/Ahhasha/Tracker-bot/internal/contracts"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TxManager struct {
	pool *pgxpool.Pool
}

func NewTxManager(pool *pgxpool.Pool) *TxManager {
	return &TxManager{pool: pool}
}

func (m *TxManager) Do(ctx context.Context, fn func(db contracts.DBTX) error) error {
	const op = "database.TxManager.Do"

	tx, err := m.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("%s: begin tx: %w", op, err)
	}

	defer tx.Rollback(ctx)

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: commit tx: %w", op, err)
	}
	return nil
}
