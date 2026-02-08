package user

import (
	"context"
	"fmt"

	"github.com/Ahhasha/Tracker-bot/internal/contracts"
	"github.com/Ahhasha/Tracker-bot/internal/contracts/add"
)

type postgresRepo struct{}

func NewPostgresRepo() add.UserRepository {
	return &postgresRepo{}
}

func (r *postgresRepo) GetIDByTgID(ctx context.Context, db contracts.DBTX, tgUserID int64) (int64, error) {
	const op = "repo.user.GetIDByTgID"

	const q = `SELECT id FROM users WHERE tg_id = $1`
	var internalUserID int64
	err := db.QueryRow(ctx, q, tgUserID).Scan(&internalUserID)
	if err != nil {
		return 0, fmt.Errorf("%s: scan: %w", op, err)
	}
	return internalUserID, nil
}
