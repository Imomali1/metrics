package main

import (
	app "github.com/Imomali1/metrics/internal/app/server"
	"github.com/Imomali1/metrics/internal/pkg/logger"
	"os"
)

func main() {
	cfg := app.LoadConfig()

	log := logger.NewLogger(os.Stdout, cfg.LogLevel, cfg.ServiceName)

	app.Run(cfg, log)
}
