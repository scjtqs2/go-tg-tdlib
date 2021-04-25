package webhook

import (
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
			log.Infof("register hook index:%d,uri:%s,secret:%s", k, v.WebHookUrl, v.WebHookSecret)
			go func() {
				for {
					msg := <-MsgCh[index]
					log.Infof("index:%d ,msg:%+v",index, msg)
					//TODO 做推送+熔断
				}
			}()
		}

	}

}

func AddMsg(k int, msg interface{}) {
	log.Debugf("msg index:%d added,%+v",k, msg)
	MsgCh[k] <- msg
}
