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

//http://localhost/?state=state-token&code=4/0AfJohXn5vkyBpavHlKu4m228rpMZthKjEwUBh4i3l3eENg1kQCgv6BoWu_xDgt4x7v7Ktw&scope=https://www.googleapis.com/auth/spreadsheets.readonly
//http://localhost/?state=state-token&code=4/0AfJohXkUXJUETeKq0KrfhCdTXYr0G_aZR2yhKt6IUJ7BJdd3Suk-xn-CtOEe8UVZ1vpjxg&scope=https://www.googleapis.com/auth/spreadsheets.readonly
