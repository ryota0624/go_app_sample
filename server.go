package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"net/http"
)

func initConfig() {
	configPath := viper.Get("config_path")
	if configPath == nil {
		configPath = "dev.toml"
	}
	viper.AutomaticEnv()
	viper.SetConfigFile(fmt.Sprintf("./configs/%s", configPath))

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func init() {
	initConfig()
}

func main() {
	e := echo.New()
	e.GET("/config", func(context echo.Context) error {
		return context.JSONPretty(http.StatusOK, viper.AllSettings(), "\t")
	})
	e.GET("/config/:keyName", func(context echo.Context) error {
		keyName := context.Param("keyName")
		return context.String(http.StatusOK, viper.GetString(keyName))
	})


	e.Static("/static", "public")
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", viper.Get("port"))))
}