package main

import (
	app "github.com/Imomali1/metrics/internal/app/agent"
)

func main() {
	cfg := app.LoadConfig()

	app.Run(cfg)
}
