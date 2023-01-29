package rdb

import (
	"coin-api/domain/model"
	"coin-api/domain/repository"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type CoinRepository struct {
	DB *gorm.DB
}

func NewCoinRepository(db *gorm.DB) repository.ICoinRepository {
	return &CoinRepository{
		DB: db,
	}
}

func (cr *CoinRepository) SelectHistoriesByUserId(uid uint) ([]model.CoinHistory, error) {
	var histories []model.CoinHistory
	result := cr.DB.Find(&histories, "userid=?", uid)
	if result.Error != nil {
		log.Fatal().Err(result.Error)
		panic(result.Error)
	}
	return histories, result.Error
}

func (cr *CoinRepository) Insert(history *model.CoinHistory) (*model.CoinHistory, error) {
	result := cr.DB.Create(history)
	if result.Error != nil {
		log.Fatal().Err(result.Error)
		panic(result.Error)
	}
	return history, result.Error
}
