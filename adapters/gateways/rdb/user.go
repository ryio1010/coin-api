package rdb

import (
	"coin-api/domain/model"
	"coin-api/domain/repository"
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
		log.Log().Msg(fmt.Sprintf("ユーザー取得処理でエラー発生 ユーザーID : %d", uid))
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
		log.Log().Msg(fmt.Sprintf("ユーザー取得処理でエラー発生 ユーザー : %s", model.CreateJsonStringFromUserModel(user)))
	}

	return user, result.Error
}

func (ur *UserRepository) Update(user *model.User) (*model.User, error) {
	// ユーザー情報更新処理
	result := ur.DB.Updates(&user)

	if result.Error != nil {
		// エラーの場合、ログを出力
		log.Log().Msg(fmt.Sprintf("ユーザー更新処理でエラー発生 ユーザー : %s", model.CreateJsonStringFromUserModel(user)))
	}

	return user, result.Error
}
