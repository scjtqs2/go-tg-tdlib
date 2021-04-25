package main

import (
	"github.com/Arman92/go-tdlib"
	"github.com/scjtqs/go-tg/app"
	"github.com/scjtqs/go-tg/config"
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"io"
	"os"
	"os/signal"
	"path"
	"time"
)

var configPath = "config.json"

func init() {
	log.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "[%time%] [%lvl%]: %msg% \n",
	})
	w, err := rotatelogs.New(path.Join("logs", "%Y-%m-%d.log"), rotatelogs.WithRotationTime(time.Hour*24))
	if err == nil {
		log.SetOutput(io.MultiWriter(os.Stderr, w))
	}
	if os.Getenv("DEBUG")=="true" {
		log.SetLevel(log.DebugLevel)
		log.Warnf("已开启Debug模式.")
	}
}

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
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}
