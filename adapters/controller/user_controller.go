package controllers

import (
	"coin-api/common"
	"coin-api/database"
	"coin-api/domain/repository"
	"coin-api/usecase/model"
	"coin-api/usecase/port"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type UserOutputFactory func(*gin.Context) ports.UserOutputPort
type UserInputFactory func(ports.UserOutputPort, repository.IUserRepository) ports.UserInputPort
type UserRepositoryFactory func(*gorm.DB) repository.IUserRepository

type UserController struct {
	OutputFactory         UserOutputFactory
	InputFactory          UserInputFactory
	UserRepositoryFactory UserRepositoryFactory
	ClientFactory         *database.PostgreSQLConnector
}

func NewUserController(outputFactory UserOutputFactory, inputFactory UserInputFactory, userRepositoryFactory UserRepositoryFactory, clientFactory *database.PostgreSQLConnector) *UserController {
	return &UserController{
		OutputFactory:         outputFactory,
		InputFactory:          inputFactory,
		UserRepositoryFactory: userRepositoryFactory,
		ClientFactory:         clientFactory,
	}
}

func (u *UserController) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// request情報をformにマッピング
		var form model.UserAddForm

		err := c.ShouldBind(&form)
		if err != nil {
			// エラーの場合、ログを出力
			log.Log().Msg(fmt.Sprintf("バインドエラー UserAddForm : %s", common.CreateJsonString(&form)))
			log.Error().Err(err).Send()
		}

		// ユーザー登録処理実行
		err = u.newInputPort(c).RegisterUser(&form)
		if err != nil {
			log.Error().Err(err).Send()
		}
	}
}

func (u *UserController) GetBalanceById() gin.HandlerFunc {
	return func(c *gin.Context) {
		// request情報からユーザーIDを取得
		uid := c.Param("userid")

		// コイン残高取得処理実行
		err := u.newInputPort(c).GetBalanceByUserId(uid)
		if err != nil {
			log.Error().Err(err).Send()
		}
	}
}

func (u *UserController) newInputPort(c *gin.Context) ports.UserInputPort {
	op := u.OutputFactory(c)
	ur := u.UserRepositoryFactory(u.ClientFactory.Conn)
	return u.InputFactory(op, ur)
}
