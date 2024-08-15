package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
	"service_on_go/internal/config"
	"service_on_go/internal/lib/logger/sl"
	"service_on_go/internal/storage/postgresql"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// TODO: init config: cleanEnv
	cfg := config.MustLoad()
	fmt.Println(cfg)

	// TODO: init logger: slog
	log := setupLogger(cfg.Env)

	//log = log.With("Project realized from nurma192", 192) // mozhno dobavit eshe parametry is sebiya
	log.Info("Starter url is shorter!!!", slog.String("env", cfg.Env))
	log.Debug("Debug messages enabled")

	// TODO: init storage: postgresql, gorm
	storage, err := postgresql.New()

	if err != nil {
		log.Error("Failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	fmt.Println(storage)

	err = storage.SaveUrl("https://google.com", "google3")
	if err != nil {
		log.Error("failed to save url", sl.Err(err))
		os.Exit(1)
	}

	log.Info("url saved to db")

	err = storage.SaveUrl("https://google.com", "google2")
	if err != nil {
		log.Error("failed to save url", sl.Err(err))
		os.Exit(1)
	}

	// TODO: init router: chi, "chi render"

	// TODO: run server

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log

}
