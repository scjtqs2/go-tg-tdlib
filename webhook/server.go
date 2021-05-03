package webhook

import (
	"github.com/Arman92/go-tdlib"
	"github.com/scjtqs/go-tg/config"
	"github.com/scjtqs/go-tg/entity"
	log "github.com/sirupsen/logrus"
)

var MsgCh []chan *entity.MSG

func Start(conf *config.JsonConfig, bot *tdlib.Client) {
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

func AddMsg(k int, msg *entity.MSG, client *tdlib.Client) {
	log.Debugf("msg index:%d added,%+v", k, msg)
	//HttpClient.AddBot(client)
	MsgCh[k] <- msg
}
