package webhook

import (
	"github.com/zelenin/go-tdlib/client"
	"github.com/scjtqs/go-tg/app/entity"
	"github.com/scjtqs/go-tg/config"
	log "github.com/sirupsen/logrus"
)

var MsgCh []chan *entity.MSG

func Start(conf *config.JsonConfig, bot *client.Client) {
	//wg := sync.WaitGroup{}
	//wg.Add(1)
	MsgCh = make([]chan *entity.MSG, len(conf.WebHook))
	for k, v := range conf.WebHook {
		index := k
		MsgCh[k] = make(chan *entity.MSG, 10)
		if v.WebHookStatus {
			pushClient := NewHttpClient(conf, v.WebHookUrl, v.WebHookSecret, 5, bot)

			log.Infof("register hook index:%d,uri:%s,secret:%s", k, v.WebHookUrl, v.WebHookSecret)
			go func() {
				for {
					msg := <-MsgCh[index]
					log.Debugf("index:%d ,msg:%+v", index, msg)
					//TODO 做推送+熔断
					pushClient.PushEvent(msg)
				}
			}()
		}
	}

}

func AddMsg(k int, msg *entity.MSG, client *client.Client) {
	log.Debugf("msg index:%d added,%+v", k, msg)
	//HttpClient.AddBot(client)
	MsgCh[k] <- msg
}
