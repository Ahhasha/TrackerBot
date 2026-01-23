package contracts

import "context"

type RegistrationRepo interface {
	UserExists(ctx context.Context, tgID int64) (bool, error)
	Create(ctx context.Context, tgID int64, username string) (int64, error)
	CreateDefaultCategory(ctx context.Context, userID int64) error
}
