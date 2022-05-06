package app

import (
	"log"
	"production_service/internal/config"
)

func main() {
	log.Print("config initializing")
	cfg := config.GetConfig()

	log.Print("logger initializing")

	app, err := app.NewApp(cfg)
	if err != nil {

	}
}
