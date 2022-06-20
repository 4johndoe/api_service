package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	_ "production_service/docs"
	"production_service/internal/config"
	"production_service/pkg/logging"
	"production_service/pkg/metric"
	"time"
)

type App struct {
	cfg        *config.Config
	logger     *logging.Logger
	router     *httprouter.Router
	httpServer *http.Server
}

func NewApp(cfg *config.Config, logger *logging.Logger) (App, error) {
	logger.Println("router init")
	router := httprouter.New()
	router.Handler(http.MethodGet, "/swagger", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently))
	router.Handler(http.MethodGet, "/swagger/*any", httpSwagger.WrapHandler)

	logger.Println("heartbeat metric init")
	metricHandler := metric.Handler{}
	metricHandler.Register(router)

	//pgConfig := postgresql.NewPgConfig(
	//	cfg.PostgreSQL.Username, cfg.PostgreSQL.Password,
	//	cfg.PostgreSQL.Host, cfg.PostgreSQL.Port, cfg.PostgreSQL.Database,
	//)
	//pgClient, err := postgresql.NewClient(context.Background(), 5, time.Second*5, pgConfig)
	//if err != nil {
	//	logger.Fatal(err)
	//}
	//
	//productStorage := storage.NewProductStorage(pgClient, logger)
	//all, err := productStorage.All(context.Background())
	//if err != nil {
	//	logger.Fatal(err)
	//}
	//logger.Fatal(all)

	return App{
		cfg:    cfg,
		logger: logger,
		router: router,
	}, nil
}

func (a *App) Run() {
	a.startHTTP()
}

func (a *App) startHTTP() {
	a.logger.Info("start http")

	var listener net.Listener

	if a.cfg.Listen.Type == config.LISTEN_TYPE_SOCK {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			a.logger.Fatal(err)
		}
		socketPath := path.Join(appDir, a.cfg.Listen.SocketFile)
		a.logger.Infof("socket path: %s", socketPath)

		a.logger.Info("create and listen unix socket")
		listener, err = net.Listen("unix", socketPath)
		if err != nil {
			a.logger.Fatal(err)
		}
	} else {
		a.logger.Info("bind application to host %s and port: %s", a.cfg.Listen.BindIp, a.cfg.Listen.Port)
		var err error
		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", a.cfg.Listen.BindIp, a.cfg.Listen.Port))
		if err != nil {
			a.logger.Fatal(err)
		}
	}

	c := cors.New(cors.Options{
		AllowedMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodOptions},
		AllowedOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
		AllowCredentials:   true,
		AllowedHeaders:     []string{ /*"Authorization", */ "Location", "Charset", "Access-Control-Allow-Origin", "Content-Type", "content-type", "Origin", "Accept", "Content-Length", "Accept-Encoding", "X-CSRF-Token"},
		OptionsPassthrough: true,
		ExposedHeaders:     []string{ /*"Access-Token", "Refresh-Token", */ "Location", "Authorization", "Content-Disposition"},
		// Enable debugging for testing, consider disabling in production
		Debug: true,
	})

	handler := c.Handler(a.router)

	a.httpServer = &http.Server{
		Handler:      handler,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	a.logger.Println("application completely init and started")

	if err := a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			a.logger.Warnf("server shutdown")
		default:
			a.logger.Fatal(err)
		}
	}
	err := a.httpServer.Shutdown(context.Background())
	if err != nil {
		a.logger.Fatal(err)
	}
}
