package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Ahhasha/Tracker-bot/internal/contracts"
	expense "github.com/Ahhasha/Tracker-bot/internal/contracts/expense"
	"github.com/Ahhasha/Tracker-bot/internal/model"
	"github.com/jackc/pgx/v5"
)

type postgresRepo struct{}

func NewPostgresRepo() expense.UserRepository {
	return &postgresRepo{}
}

func (r *postgresRepo) GetIDByTgID(ctx context.Context, db contracts.DBTX, tgUserID int64) (int64, error) {
	const op = "repo.user.GetIDByTgID"

	const q = `SELECT id FROM users WHERE tg_id = $1`
	var internalUserID int64
	err := db.QueryRow(ctx, q, tgUserID).Scan(&internalUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, model.ErrUserNotRegistered
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return internalUserID, nil
}
