package app

import (
	"fmt"
	"github.com/scjtqs/go-tg-tdlib/config"
	"os"
	"os/signal"
	"syscall"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

func Start(conf *config.JsonConfig)  {
	cronConf:=config.LoadCron()
	client :=NewClient(conf)
	// Handle Ctrl+C
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		client.Cli.DestroyInstance()
		os.Exit(1)
	}()
	//定时任务开启
	client.Cron = cron.New()
	for k,v := range cronConf.Cron{
		_, err := client.Cron.AddFunc(v.Cron, func() {
			log.Infof("crontab %d start",(k+1))
			client.SendMessageByName(v.ToUserName,v.TextMsg)
		})
		if err!=nil {
			panic("cron start with error:"+err.Error())
		}
	}

	//
	// rawUpdates gets all updates comming from tdlib
	rawUpdates := client.Cli.GetRawUpdatesChannel(100)
	for update := range rawUpdates {
		// Show all updates
		fmt.Println(update.Data)
		fmt.Print("\n\n")
	}

}
