package start

import "context"

type TxManager interface {
	Do(ctx context.Context, fn func(db DBTX) error) error
}
