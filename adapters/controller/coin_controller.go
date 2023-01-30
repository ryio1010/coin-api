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

type CoinOutputFactory func(*gin.Context) ports.CoinOutputPort
type CoinInputFactory func(ports.CoinOutputPort, repository.ICoinRepository, repository.IUserRepository) ports.CoinInputPort
type CoinRepositoryFactory func(*gorm.DB) repository.ICoinRepository

type CoinController struct {
	OutputFactory         CoinOutputFactory
	InputFactory          CoinInputFactory
	CoinRepositoryFactory CoinRepositoryFactory
	UserRepositoryFactory UserRepositoryFactory
	ClientFactory         *database.PostgreSQLConnector
}

func NewCoinController(outputFactory CoinOutputFactory, inputFactory CoinInputFactory, coinRepositoryFactory CoinRepositoryFactory, userRepositoryFactory UserRepositoryFactory, clientFactory *database.PostgreSQLConnector) *CoinController {
	return &CoinController{
		OutputFactory:         outputFactory,
		InputFactory:          inputFactory,
		CoinRepositoryFactory: coinRepositoryFactory,
		UserRepositoryFactory: userRepositoryFactory,
		ClientFactory:         clientFactory,
	}
}

func (c *CoinController) AddUseCoin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// request情報をformにマッピング
		var form model.CoinAddUseForm
		err := ctx.ShouldBind(&form)

		if err != nil {
			// エラーの場合、ログを出力
			log.Log().Msg(fmt.Sprintf("バインドエラー CoinAddUseForm : %s", common.CreateJsonString(&form)))
			log.Error().Err(err).Send()
		}

		// コイン追加消費処理
		err = c.newInputPort(ctx).AddUseCoin(&form)
		if err != nil {
			log.Error().Stack().Err(err).Send()
		}
	}
}

func (c *CoinController) SendCoin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// request情報をformにマッピング
		var form model.CoinSendForm
		err := ctx.ShouldBind(&form)

		if err != nil {
			log.Log().Msg(fmt.Sprintf("バインドエラー CoinSendForm : %s", common.CreateJsonString(&form)))
			log.Error().Err(err).Send()
		}

		// コイン送金処理
		err = c.newInputPort(ctx).SendCoin(&form)

		if err != nil {
			log.Error().Stack().Err(err).Send()
		}
	}
}

func (c *CoinController) GetHistoryByUserId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// request情報からユーザーIDを取得
		uid := ctx.Param("userid")

		// コイン履歴取得処理
		err := c.newInputPort(ctx).SelectHistoryByUserId(uid)

		if err != nil {
			log.Error().Stack().Err(err).Send()
		}
	}
}

func (c *CoinController) newInputPort(ctx *gin.Context) ports.CoinInputPort {
	op := c.OutputFactory(ctx)
	cr := c.CoinRepositoryFactory(c.ClientFactory.Conn)
	ur := c.UserRepositoryFactory(c.ClientFactory.Conn)
	return c.InputFactory(op, cr, ur)
}
