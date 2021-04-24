package main

import (
	"github.com/Arman92/go-tdlib"
	"github.com/scjtqs/go-tg-tdlib/app"
	"github.com/scjtqs/go-tg-tdlib/config"
	"os"
)

var configPath = "config.json"

func main() {
	tdlib.SetLogVerbosityLevel(1)
	tdlib.SetFilePath("./errors.txt")
	var conf *config.JsonConfig
	switch os.Getenv("IS_DOCKER") {
	case "true":
		conf = config.DefaultConfig()
	default:
		if !config.PathExists(configPath) {
			conf = config.DefaultConfig()
			_ = conf.Save(configPath)
		}
		conf = config.Load(configPath)
	}

	app.Start(conf)

}
