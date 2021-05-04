package app

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/scjtqs/go-tg/app/web"
	"github.com/scjtqs/go-tg/app/webhook"
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
	//client.Cron = cron.New()
	// Seconds field, optional
	client.Cron = cron.New(cron.WithParser(cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)))
	for k, v := range cronConf.Cron {
		name := v.ToUserName
		crontab := v.Cron
		text := v.TextMsg
		log.Infof("index %d cron %s  to user %s registed", (k + 1), v.Cron, name)
		_, err := client.Cron.AddFunc(crontab, func() {
			log.Infof("crontab to username %s start", name)
			e := client.SendMessageByName(name, text)
			if e != nil {
				log.Errorf("send message to username error err:%v \n\n", name, e)
			}
		})
		if err != nil {
			panic("cron start with error:" + err.Error())
		}
	}

	// filter msg
	if conf.WebHook != nil {
		for k, v := range conf.WebHook {
			client.FilterMsg(k, v)
		}
	}

	client.Cron.Start()
	/**
	rawUpdates gets all updates comming from tdlib
	rawUpdates := client.Cli.GetRawUpdatesChannel(100)
	client.Cli.GetRawUpdatesChannel(100)
	for update := range rawUpdates {
		// Show all updates
		//log.Debug(update.Data)
		//log.Debug("\n")
		if update.Data!=nil{

		}
	}
	*/
	if client.Conf.WebApi.WebApiStatus {
		//开启 http api
		web.HttpServer.
			Run(fmt.Sprintf("%s:%s", conf.WebApi.WebApiHost, conf.WebApi.WebApiPort), conf.WebApi.WebApiToken, client.Cli, conf)
	}
	webhook.Start(conf, client.Cli)
	log.Infof("started ok \n")
}
