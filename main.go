package main

import (
	"flag"
	"fmt"
	"github.com/Arman92/go-tdlib/client"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/scjtqs/go-tg/app"
	"github.com/scjtqs/go-tg/config"
	"github.com/scjtqs/go-tg/utils"
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strings"
	"time"
)

var c string
var d bool
var h bool
var Version = "unknown"

func init() {
	var debug bool
	flag.StringVar(&c, "c", config.ConfigPath, "configuration filename default is config.json")
	flag.BoolVar(&d, "d", false, "running as a daemon")
	flag.BoolVar(&debug, "D", false, "debug mode")
	flag.BoolVar(&h, "h", false, "this help")
	flag.Parse()

	logFormatter := (&easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "[%time%] [%lvl%]: %msg% \n",
	})
	w, err := rotatelogs.New(path.Join("logs", "%Y-%m-%d.log"), rotatelogs.WithRotationTime(time.Hour*24))
	if err != nil {
		log.Errorf("rotatelogs init err: %v", err)
		panic(err)
	}
	if os.Getenv("DEBUG") == "true" || debug {
		log.SetLevel(log.DebugLevel)
		log.Warnf("已开启Debug模式.")
	}
	log.AddHook(config.NewLocalHook(w, logFormatter, config.GetLogLevel("info")...))
}

func main() {
	if h {
		help()
	}
	if d {
		Daemon()
	}
	client.SetLogVerbosityLevel(1)
	client.SetFilePath("./errors.txt")
	var conf *config.JsonConfig
	switch os.Getenv("IS_DOCKER") {
	case "true":
		conf = config.DefaultConfig()
		conf.Save(config.ConfigPath)
	default:
		if !utils.PathExists(config.ConfigPath) {
			conf = config.DefaultConfig()
			_ = conf.Save(config.ConfigPath)
		}
		conf = config.Load(config.ConfigPath)
	}
	app.Start(conf)
	log.Infof("started ok verison=%s", Version)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	os.Exit(1)
}

// Daemon go-tg server 的 daemon的实现函数
func Daemon() {
	args := os.Args[1:]

	execArgs := make([]string, 0)

	l := len(args)
	for i := 0; i < l; i++ {
		if strings.Index(args[i], "-d") == 0 {
			continue
		}

		execArgs = append(execArgs, args[i])
	}

	proc := exec.Command(os.Args[0], execArgs...)
	err := proc.Start()

	if err != nil {
		panic(err)
	}

	log.Info("[PID] ", proc.Process.Pid)
	// pid写入到pid文件中，方便后续stop的时候kill
	pidErr := savePid("go-tg.pid", fmt.Sprintf("%d", proc.Process.Pid))
	if pidErr != nil {
		log.Errorf("save pid file error: %v", pidErr)
	}

	os.Exit(0)
}

// savePid 保存pid到文件中，便于后续restart/stop的时候kill pid用。
func savePid(path string, data string) error {
	return utils.WriteAllText(path, data)
}

// help cli命令行-h的帮助提示
func help() {
	fmt.Printf(`go-tg service
version: %s

Usage:

server [OPTIONS]

Options:
`, Version)
	flag.PrintDefaults()
	os.Exit(0)
}
