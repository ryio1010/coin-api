package config

const (
	port       = "8081"
	dbUser     = "admin"
	dbPassword = "admin"
	dbName     = "coin_db"
	dbHost     = "localhost"
	dbPort     = "5433"
)

type AppConfig struct {
	HTTPInfo       *HttpInfo
	PostgreSQLInfo *PostgreSQLInfo
}

type HttpInfo struct {
	Addr string
}

type PostgreSQLInfo struct {
	User     string
	Password string
	DbName   string
	Host     string
	Port     string
}

func LoadConfig() *AppConfig {
	addr := ":" + port

	httpInfo := &HttpInfo{
		Addr: addr,
	}

	dbInfo := &PostgreSQLInfo{
		User:     dbUser,
		Password: dbPassword,
		DbName:   dbName,
		Host:     dbHost,
		Port:     dbPort,
	}

	conf := AppConfig{
		HTTPInfo:       httpInfo,
		PostgreSQLInfo: dbInfo,
	}

	return &conf
}
