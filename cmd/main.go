package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ktrntrsv/userBalanceService/internal/adapters/api"
	"github.com/ktrntrsv/userBalanceService/internal/config"
	"github.com/ktrntrsv/userBalanceService/pkg/httpserver"
	"github.com/ktrntrsv/userBalanceService/pkg/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig("./config.yml")
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	l := logger.New(cfg.Logger.Level)

	// HTTP Server
	handler := gin.New()
	gin.SetMode(cfg.Server.Mode)
	api.NewRouter(handler, l)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.Server.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

}
