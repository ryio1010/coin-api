package rdb

import (
	"coin-api/common"
	"coin-api/domain/model"
	"coin-api/domain/repository"
	"fmt"
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
	// 取得用モデル定義
	var histories []model.CoinHistory

	// idに紐づく全履歴取得
	result := cr.DB.Find(&histories, "userid=?", uid)
	if result.Error != nil {
		// エラーまたはレコードを取得できない場合、ログを出力
		log.Log().Msg(fmt.Sprintf("履歴取得処理でエラー発生 ユーザーID : %d", uid))
	}

	return histories, result.Error
}

func (cr *CoinRepository) Insert(history *model.CoinHistory) (*model.CoinHistory, error) {
	// 履歴登録処理
	result := cr.DB.Create(history)
	if result.Error != nil {
		// エラーの場合、ログを出力
		log.Log().Msg(fmt.Sprintf("履歴登録処理でエラー発生 履歴 : %s", common.CreateJsonString(history)))
	}

	return history, result.Error
}

func (cr *CoinRepository) BatchInsert(histories []*model.CoinHistory) ([]*model.CoinHistory, error) {
	// 履歴登録処理
	results := cr.DB.Create(histories)
	if results.Error != nil {
		// エラーの場合、ログを出力
		log.Log().Msg(fmt.Sprintf("履歴一括登録処理でエラー発生 履歴 : %s", common.CreateJsonString(histories)))
	}

	return histories, results.Error
}
