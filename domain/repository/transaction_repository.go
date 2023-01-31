package repository

import (
	"context"
)

type ITxRepository interface {
	DoInTx(ctx context.Context, f func(ctx context.Context) (interface{}, error)) (interface{}, error)
}
