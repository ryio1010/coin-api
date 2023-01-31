package rdb

import (
	"coin-api/common"
	"coin-api/domain/model"
	"coin-api/domain/repository"
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.IUserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (ur *UserRepository) SelectById(uid uint) (*model.User, error) {
	// 取得用モデル定義
	user := model.User{}

	// id検索でのユーザー取得処理
	result := ur.DB.First(&user, "id=?", uid)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// エラーまたはレコードを取得できない場合、ログを出力
		log.Error().Msg(fmt.Sprintf("ユーザー取得処理でエラー発生 ユーザーID : %d", uid))
		return nil, result.Error
	}

	return &user, result.Error
}

func (ur *UserRepository) Insert(user *model.User) (*model.User, error) {
	// ユーザー新規登録時にコイン残高を0で登録
	balance := 0
	user.CoinBalance = &balance

	// ユーザー登録処理
	result := ur.DB.Create(&user)

	if result.Error != nil {
		// エラーの場合、ログを出力
		log.Error().Msg(fmt.Sprintf("ユーザー登録処理でエラー発生 ユーザー : %s", common.CreateJsonString(user)))
		return nil, result.Error
	}

	return user, result.Error
}

func (ur *UserRepository) Update(ctx context.Context, user *model.User) (*model.User, error) {
	// トランザクション取得
	tr, ok := GetTx(ctx)
	if !ok {
		tr = ur.DB
	}

	// ユーザー情報更新処理
	result := tr.Updates(&user)

	if result.Error != nil {
		// エラーの場合、ログを出力
		log.Error().Msg(fmt.Sprintf("ユーザー更新処理でエラー発生 ユーザー : %s", common.CreateJsonString(user)))
		return nil, result.Error
	}

	return user, result.Error
}
