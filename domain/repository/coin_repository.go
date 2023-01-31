package repository

import (
	"coin-api/domain/model"
	"context"
)

type ICoinRepository interface {
	SelectHistoriesByUserId(uid uint) ([]model.CoinHistory, error)
	Insert(ctx context.Context, history *model.CoinHistory) (*model.CoinHistory, error)
	BatchInsert(ctx context.Context, histories []*model.CoinHistory) ([]*model.CoinHistory, error)
}
