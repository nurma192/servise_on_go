package main

import (
	"fmt"
	"service_on_go/internal/config"
)

func main() {
	// TODO: init config: cleanEnv
	cfg := config.MustLoad()
	fmt.Println(cfg)

	// TODO: init logger: slog

	// TODO: init storage: sqlLite

	// TODO: init router: chi, "chi render"

	// TODO: run server

}
