package app

import (
	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"production_service/internal/config"
	"production_service/pkg/logging"
)

type App struct {
	config *config.Config
	logger *logging.Logger
}

func NewApp(cfg *config.Config, logger *logging.Logger) (App, error) {
	logger.Println("router init")
	router := httprouter.New()
	router.Handler(http.MethodGet, "/swagger", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently))
	router.Handler(http.MethodGet, "/swagger/*any", httpSwagger.WrapHandler)

	return App{
		config: cfg,
		logger: logger,
	}, nil
}
