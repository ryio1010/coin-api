package main

import (
	"coin-api/drivers"
	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

func main() {
	// log設定
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	// Gin設定
	engine := drivers.InitRouter()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	// サーバー起動
	err := engine.Run(":8081")
	if err != nil {
		panic("サーバーの起動に失敗しました。")
	}
}
