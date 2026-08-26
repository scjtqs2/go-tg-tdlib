package app

import (
	"context"
	"github.com/scjtqs/go-tg/config"
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
	"strconv"
)

func (a *AppClient) FilterMsg(index int, conf *config.WebHook) {
	if !conf.WebHookStatus {
		return
	}
	var chatIds []string
	if conf.WebHookFilter != nil && conf.WebHookFilter.FilterStatus {
		for _, v := range conf.WebHookFilter.FilterNames {
			id, err := a.Cli.SearchPublicChat(context.Background(), &client.SearchPublicChatRequest{Username: v})
			if err != nil {
				log.Errorf("search for chat faild name=%s,err=%v", v, err)
				continue
			}
			chatIds = append(chatIds, strconv.FormatInt(id.Id, 10))
		}
	}
}

