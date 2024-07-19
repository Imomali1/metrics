package main

import (
	"os"

	app "github.com/Imomali1/metrics/internal/app/server"
	"github.com/Imomali1/metrics/internal/pkg/logger"
	"fmt"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func printServerInfo() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}

func main() {
	printServerInfo()

	cfg := app.LoadConfig()

	log := logger.NewLogger(os.Stdout, cfg.LogLevel, cfg.ServiceName)

	app.Run(cfg, log)
}
