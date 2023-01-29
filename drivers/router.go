package drivers

import (
	"coin-api/adapters/controller"
	"coin-api/adapters/gateways/rdb"
	"coin-api/database"
	"coin-api/usecase/interactor"
	"coin-api/usecase/presenter"
	"github.com/gin-gonic/gin"
)

const (
	apiVersion  = "/v1"
	userApiRoot = apiVersion + "/user"
	coinApiRoot = apiVersion + "/coin"
)

func InitRouter() *gin.Engine {
	g := gin.Default()

	// DB
	con := database.NewPostgreSQLConnector()

	// User
	uop := presenter.NewUserOutputPort
	uip := interactor.NewUserUseCase
	ur := rdb.NewUserRepository

	// Coin
	cop := presenter.NewCoinOutputPort
	cip := interactor.NewCoinUseCase
	cr := rdb.NewCoinRepository

	ug := g.Group(userApiRoot)
	{
		uc := controllers.NewUserController(uop, uip, ur, con)
		// POST RegisterUserAPI
		ug.POST("", uc.CreateUser())
		// GET GetBalanceByUserIdAPI
		ug.GET("/:userid", uc.GetBalanceById())
	}

	cg := g.Group(coinApiRoot)
	{
		cc := controllers.NewCoinController(cop, cip, cr, ur, con)
		// POST AddUseCoinAPI
		cg.POST("", cc.AddUseCoin())
		// PUT SendCoinAPI
		cg.PUT("/send", cc.SendCoin())
		// GET GetHistoriesById
		cg.GET("/:userid", cc.GetHistoryByUserId())
	}

	return g
}
