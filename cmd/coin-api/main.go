package main

import (
	"coin-api/drivers"
	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	engine := drivers.InitRouter()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	err := engine.Run(":8081")
	if err != nil {
		panic("ERROR!!!")
	}
}
