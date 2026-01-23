package contracts

import "context"

type RegistrationRepo interface {
	UpsertUser(ctx context.Context, tgID int64, username string) (userID int64, created bool, err error)
	CreateDefaultCategory(ctx context.Context, userID int64) error
}
