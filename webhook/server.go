package webhook

import (
	"github.com/Arman92/go-tdlib"
	"github.com/scjtqs/go-tg/config"
	log "github.com/sirupsen/logrus"
)

var MsgCh []chan interface{}

func Start(conf *config.JsonConfig) {
	//wg := sync.WaitGroup{}
	//wg.Add(1)
	MsgCh = make([]chan interface{}, len(conf.WebHook))
	for k, v := range conf.WebHook {
		index := k
		MsgCh[k] = make(chan interface{}, 10)
		if v.WebHookStatus {
			httpClient := NewHttpClient()
			httpClient.Init(conf, v.WebHookUrl, v.WebHookSecret, 5)
			log.Infof("register hook index:%d,uri:%s,secret:%s", k, v.WebHookUrl, v.WebHookSecret)
			go func() {
				for {
					msg := <-MsgCh[index]
					log.Infof("index:%d ,msg:%+v", index, msg)
					//TODO 做推送+熔断
				}
			}()
		}
	}

}

func AddMsg(k int, msg map[string]interface{}, client *tdlib.Client) {
	log.Debugf("msg index:%d added,%+v", k, msg)
	HttpClient.AddBot(client)
	MsgCh[k] <- msg
}
