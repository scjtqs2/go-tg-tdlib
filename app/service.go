package app

import (
	"github.com/robfig/cron/v3"
	"github.com/scjtqs/go-tg/config"
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
		client.Cron.Stop()
		os.Exit(1)
	}()
	//定时任务开启
	client.Cron = cron.New()
	for k, v := range cronConf.Cron {
		name:=v.ToUserName
		crontab:=v.Cron
		text:=v.TextMsg
		log.Infof("cron %d to user %s registed",(k+1),name)
		_, err := client.Cron.AddFunc(crontab, func() {
			log.Infof("crontab to username %s start \n\n",name )
			e := client.SendMessageByName(name, text)
			if e != nil {
				log.Errorf("send message to username error err:%v \n\n", name, e)
			}
		})
		if err != nil {
			panic("cron start with error:" + err.Error())
		}
	}
	//测试
	//e := client.SendMessageByName("LvanLamCommitCodeBot", "testing")
	//if e != nil {
	//	log.Errorf("send message %d error err:%v \n\n", 0, e)
	//}
	client.Cron.Start()
	// rawUpdates gets all updates comming from tdlib
	rawUpdates := client.Cli.GetRawUpdatesChannel(100)
	for update := range rawUpdates {
		// Show all updates
		log.Debug(update.Data,"\n\n")
		log.Debug("\n\n")
	}
	log.Infof("started ok \n")
}
