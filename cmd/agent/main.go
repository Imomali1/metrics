package main

import (
	_ "net/http/pprof"

	app "github.com/Imomali1/metrics/internal/app/agent"
)

func main() {
	app.Run()
}
