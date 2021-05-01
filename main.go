package main

import (
	"github.com/Arman92/go-tdlib"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/scjtqs/go-tg/app"
	"github.com/scjtqs/go-tg/config"
	"github.com/scjtqs/go-tg/utils"
	"github.com/scjtqs/go-tg/webhook"
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"io"
	"os"
	"os/signal"
	"path"
	"time"
)

var configPath = "config.json"

func init() {
	logFormatter:=(&easy.Formatter{
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
	log.AddHook(config.NewLocalHook(w, logFormatter, config.GetLogLevel("warn")...))
}

func main() {
	tdlib.SetLogVerbosityLevel(1)
	tdlib.SetFilePath("./errors.txt")
	var conf *config.JsonConfig
	switch os.Getenv("IS_DOCKER") {
	case "true":
		conf = config.DefaultConfig()
		conf.Save(configPath)
	default:
		if !utils.PathExists(configPath) {
			conf = config.DefaultConfig()
			_ = conf.Save(configPath)
		}
		conf = config.Load(configPath)
	}
	webhook.Start(conf)
	app.Start(conf)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}
