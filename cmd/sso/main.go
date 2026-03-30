package main

import (
	"log/slog"
	"os"
	"os/signal"
	"pet-grpc/internal/app"
	"pet-grpc/internal/config"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info(
		"starting app",
		slog.Any("cfg", cfg),
	)

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	// application.GRPCSrv.MustRun()

	go application.GRPCSrv.MustRun() //

	stop := make(chan os.Signal, syscall.SIGINT)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("stopping application", slog.String("signal", sign.String()))
	application.GRPCSrv.Stop()

	log.Info("Application stopped")

	//  серивис либо проект наш(бинарник), либо объект из сервисного слоя приложения(есть траспортные(вроде понятно что), сервисные(бизнес логика) )
	//  TODO: запустить grpc-сервер приложения

}

// создаем логер для разных сред окружения
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log

}
