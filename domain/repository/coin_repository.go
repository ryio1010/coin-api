package repository

import (
	"coin-api/domain/model"
)

type ICoinRepository interface {
	SelectHistoriesByUserId(uid uint) ([]model.CoinHistory, error)
	Insert(history *model.CoinHistory) (*model.CoinHistory, error)
}
