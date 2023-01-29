package controllers

import (
	"coin-api/database"
	"coin-api/domain/repository"
	"coin-api/usecase/model"
	"coin-api/usecase/port"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"strconv"
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
	return func(context *gin.Context) {
		// request情報をformにマッピング
		var form model.CoinAddUseForm
		err := context.ShouldBind(&form)
		if err != nil {
			log.Error().Err(err)
		}

		err = c.newInputPort(context).AddUseCoin(&form)
		if err != nil {
			log.Error().Err(err)
		}
	}
}

func (c *CoinController) SendCoin() gin.HandlerFunc {
	return func(context *gin.Context) {
		var form model.CoinSendForm
		err := context.ShouldBind(&form)
		if err != nil {
			log.Error().Err(err)
		}
		fmt.Println(form)
		err = c.newInputPort(context).SendCoin(&form)
		if err != nil {
			log.Error().Err(err)
		}
	}
}

func (c *CoinController) GetHistoryByUserId() gin.HandlerFunc {
	return func(context *gin.Context) {
		uid := context.Param("userid")
		uidInt, _ := strconv.ParseUint(uid, 10, 64)

		err := c.newInputPort(context).SelectHistoryByUserId(uint(uidInt))
		if err != nil {
			log.Error().Err(err)
		}
	}
}

func (c *CoinController) newInputPort(context *gin.Context) ports.CoinInputPort {
	op := c.OutputFactory(context)
	cr := c.CoinRepositoryFactory(c.ClientFactory.Conn)
	ur := c.UserRepositoryFactory(c.ClientFactory.Conn)
	return c.InputFactory(op, cr, ur)
}
