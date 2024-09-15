package main

import (
	"fmt"
	"github.com/Imomali1/metrics/internal/pkg/logger"
	"os"

	app "github.com/Imomali1/metrics/internal/app/agent"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func printAgentInfo() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}

func main() {
	printAgentInfo()

	cfg := app.LoadConfig()

	log := logger.NewLogger(os.Stdout, cfg.LogLevel, cfg.ServiceName)

	if err := app.Run(cfg, log); err != nil {
		log.Fatal().Err(err).Send()
	}
}
