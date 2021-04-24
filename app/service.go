package app

import (
	"github.com/robfig/cron/v3"
	"github.com/scjtqs/go-tg-tdlib/config"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func Start(conf *config.JsonConfig) {
	cronConf := config.LoadCron()
	client := NewClient(conf)
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
	for k, v := range cronConf.Cron {
		_, err := client.Cron.AddFunc(v.Cron, func() {
			log.Infof("crontab %d start \n\n", (k + 1))
			e := client.SendMessageByName(v.ToUserName, v.TextMsg)
			if e != nil {
				log.Errorf("send message %d error err:%v \n\n", k, e)
			}
		})
		if err != nil {
			panic("cron start with error:" + err.Error())
		}
	}

	// rawUpdates gets all updates comming from tdlib
	rawUpdates := client.Cli.GetRawUpdatesChannel(100)
	for update := range rawUpdates {
		// Show all updates
		log.Debug(update.Data,"\n\n")
		log.Debug("\n\n")
	}

}
