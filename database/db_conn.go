package database

import (
	"coin-api/config"
	"coin-api/domain/model"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgreSQLConnector struct {
	Conn *gorm.DB
}

func NewPostgreSQLConnector() *PostgreSQLConnector {
	conf := config.LoadConfig()
	dsn := postgresConnInfo(*conf.PostgreSQLInfo)

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// gormのmigrate
	err = conn.AutoMigrate(&model.User{}, &model.CoinHistory{})

	return &PostgreSQLConnector{
		Conn: conn,
	}
}

func postgresConnInfo(postgresInfo config.PostgreSQLInfo) string {
	dataSourceName := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		postgresInfo.Host,
		postgresInfo.User,
		postgresInfo.Password,
		postgresInfo.DbName,
	)

	return dataSourceName
}
