package config

const (
	dbUser     = "admin"
	dbPassword = "admin"
	dbName     = "coin_db"
	dbHost     = "coin_db"
	dbPort     = "5433"
)

type AppConfig struct {
	PostgreSQLInfo *PostgreSQLInfo
}
type PostgreSQLInfo struct {
	User     string
	Password string
	DbName   string
	Host     string
	Port     string
}

func LoadConfig() *AppConfig {
	dbInfo := &PostgreSQLInfo{
		User:     dbUser,
		Password: dbPassword,
		DbName:   dbName,
		Host:     dbHost,
		Port:     dbPort,
	}

	conf := AppConfig{
		PostgreSQLInfo: dbInfo,
	}

	return &conf
}
