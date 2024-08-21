package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"golang.org/x/exp/slog"
	"os"
	"service_on_go/internal/config"
	"service_on_go/internal/http-server/middleware/logger"
	"service_on_go/internal/lib/logger/handler/slogpretty"
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
	//configPath := "./config/local.yaml"
	cfg := config.MustLoad()
	fmt.Println(cfg)

	// TODO: init logger: slog
	var log *slog.Logger = setupLogger(cfg.Env)

	//log = log.With("Project realized from nurma192", 192) // mozhno dobavit eshe parametry is sebiya
	log.Info("Starter url is shorter!!!", slog.String("env", cfg.Env))
	log.Debug("Debug messages enabled")
	log.Error("error message: massage")

	// TODO: init storage: postgresql, gorm
	storage, err := postgresql.New()

	if err != nil {
		log.Error("Failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	fmt.Println("Storage created: ", storage)

	// TODO: init router: gin
	router := gin.New()

	// middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(logger.New(log))

	// TODO: run server

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = setupPrettySlog()
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

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{Level: slog.LevelDebug},
	}

	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}
