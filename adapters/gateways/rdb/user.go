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
	user := model.User{}

	result := ur.DB.First(&user, "id=?", uid)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// エラーまたはレコードを取得できない場合、ログを出力
		log.Log().Msg(fmt.Sprintf("ユーザー取得処理でエラー発生 ユーザーID : %d", uid))
		log.Error().Err(result.Error)
	}

	return &user, result.Error
}

func (ur *UserRepository) Insert(user model.User) (*model.User, error) {
	// ユーザー新規登録時にコイン残高を0で登録
	balance := 0
	user.CoinBalance = &balance

	result := ur.DB.Create(&user)

	if result.Error != nil {
		log.Error().Err(result.Error)
	}

	return &user, result.Error
}

func (ur *UserRepository) Update(user *model.User) (*model.User, error) {
	result := ur.DB.Updates(&user)
	if result.Error != nil {
		log.Error().Err(result.Error)
	}
	return user, result.Error
}
