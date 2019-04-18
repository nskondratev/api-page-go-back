package conf

import (
	"flag"
	"github.com/joho/godotenv"
	"os"
)

type AppConfig struct {
	DBConnectionString string
	Addr               string
	BaseUrl            string
}

func GetAppConfig() (*AppConfig, error) {
	const (
		defaultConnectionString = ""
		defaultAddr             = ""
		defaultBaseUrl          = ""
	)
	conf := &AppConfig{}

	if err := godotenv.Load(); err != nil {
		return conf, err
	}

	flag.StringVar(&conf.DBConnectionString, "db", defaultConnectionString, "Database connection string")
	if len(conf.DBConnectionString) < 1 && len(os.Getenv("DB_CONNECTION_STRING")) > 0 {
		conf.DBConnectionString = os.Getenv("DB_CONNECTION_STRING")
	}
	flag.StringVar(&conf.Addr, "addr", defaultAddr, "Application address to listen")
	if len(conf.Addr) < 1 && len(os.Getenv("ADDR")) > 0 {
		conf.Addr = os.Getenv("ADDR")
	}
	flag.StringVar(&conf.BaseUrl, "base-url", defaultBaseUrl, "Application base url")
	if len(conf.BaseUrl) < 1 && len(os.Getenv("BASE_URL")) > 0 {
		conf.BaseUrl = os.Getenv("BASE_URL")
	}
	flag.Parse()
	return conf, nil
}
