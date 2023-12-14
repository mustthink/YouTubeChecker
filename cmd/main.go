package main

import (
	"flag"

	"github.com/mustthink/YouTubeChecker/config"
	"github.com/mustthink/YouTubeChecker/internal"
)

func main() {
	configPath := flag.String("config", config.DefaultConfig, "path to config")
	isDebug := flag.Bool("d", true, "is debug active")
	flag.Parse()

	app := internal.NewApplication(*configPath, *isDebug)
	app.Run()
}
