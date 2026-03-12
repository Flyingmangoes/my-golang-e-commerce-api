package config

import (
	"errors"
	"os"
)


type Application struct {
	DBConf *DBConfig
	ServConf *ServerConfig
	AppConf *AppConfig
	RateConf *RateLimitingConfig
}

type AppConfig struct {
	CustomAliasLength int
}

type ServerConfig struct {
	Host string
	Port string
}

type DBConfig struct {
	DBAddr string
}

type RateLimitingConfig struct {
	RequestPerMinute int
	Burst int
}

func (a *Application) Validate() error {
	if a.DBConf.DBAddr == ""{
		return errors.New("DB_URL is required")
	}

	if a.ServConf.Host == "" {
		return errors.New("SERV_HOST is required")
	}

	if a.ServConf.Port == "" {
		return errors.New("SERV_PORT is required")
	}

    return nil
}

func NewConfig() *Application {
	return &Application{
		DBConf: &DBConfig{
			DBAddr: os.Getenv("DB_URL"),
		},
		ServConf: &ServerConfig{
			Host: os.Getenv("SERV_HOST"),
			Port: os.Getenv("SERV_PORT"),
		},
		AppConf: &AppConfig{
			CustomAliasLength: 6,
		},
		RateConf: &RateLimitingConfig{
			RequestPerMinute: 30,
			Burst: 3,
		},
	}
}
	