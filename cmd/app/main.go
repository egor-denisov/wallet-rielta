package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/egor-denisov/wallet-rielta/config"
	app "github.com/egor-denisov/wallet-rielta/internal/app"
	sl "github.com/egor-denisov/wallet-rielta/pkg/logger"
)

func main() {
	// Init configuration
	cfg := config.MustLoad()

	// Init logger
	log := sl.SetupLogger(cfg.Log.Level)

	application := app.New(log, cfg)

	// Run servers
	go func() {
		application.HTTPServer.MustRun()
	}()

	go func() {
		application.RMQServer.MustRun()
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-stop:
	case <-application.RMQServer.Notify():
	}

	log.Info("Starting graceful shutdown")

	if err := application.HTTPServer.Shutdown(); err != nil {
		log.Error("HTTPServer.Shutdown error", sl.Err(err))
	}

	if err := application.RMQServer.Shutdown(); err != nil {
		log.Error("RMQServer.Shutdown error", sl.Err(err))
	}

	if err := application.DB.Close(); err != nil {
		log.Error("Close db connection error", sl.Err(err))
	}

	log.Info("Gracefully stopped")
}
