package main

import (
	"log"
	"production_service/internal/app"
	"production_service/internal/config"
	"production_service/pkg/logging"
)

func main() {
	log.Print("config initializing")
	cfg := config.GetConfig()

	log.Print("logger initializing")
	logging.Init(cfg.AppConfig.LogLevel) // 40 49
	logger := logging.GetLogger()

	a, err := app.NewApp(cfg, &logger)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("running app")
	a.Run()
}
