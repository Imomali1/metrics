package main

import (
	"os"

	app "github.com/Imomali1/metrics/internal/app/server"
	"github.com/Imomali1/metrics/internal/pkg/logger"
)

func main() {
	cfg := app.LoadConfig()

	log := logger.NewLogger(os.Stdout, cfg.LogLevel, cfg.ServiceName)

	app.Run(cfg, log)
}
