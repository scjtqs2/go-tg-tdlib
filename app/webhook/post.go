package webhook

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"github.com/zelenin/go-tdlib/client"
	"github.com/guonaihong/gout"
	"github.com/scjtqs/go-tg/app/entity"
	"github.com/scjtqs/go-tg/config"
	log "github.com/sirupsen/logrus"
	"time"
)

var HttpClient *httpClient

type httpClient struct {
	Bot     *client.Client
	Secret  string
	Addr    string
	Timeout int32
	Conf    *config.JsonConfig
}

func NewHttpClient(config *config.JsonConfig, addr string, secret string, timeout int32, bot *client.Client) *httpClient {
	HttpClient = &httpClient{
		Secret:  secret,
		Addr:    addr,
		Timeout: timeout,
		Conf:    config,
		Bot:     bot,
	}
	log.Infof("HTTP POST上报器已启动: %v", addr)
	return HttpClient
}

func (c *httpClient) PushEvent(m *entity.MSG) {
	var res string
	err := gout.POST(c.Addr).SetJSON(m).BindBody(&res).SetHeader(func() gout.H {
		h := gout.H{
			"X-Self-ID":  c.Conf.Phone,
			"User-Agent": "TGHttp/4.15.0",
		}
		if c.Secret != "" {
			mac := hmac.New(sha1.New, []byte(c.Secret))
			mac.Write([]byte(m.ToJson()))
			h["X-Signature"] = "sha1=" + hex.EncodeToString(mac.Sum(nil))
		}
		return h
	}()).SetTimeout(time.Second * time.Duration(c.Timeout)).Do()
	if err != nil {
		log.Warnf("上报Event数据到 %v 失败: %v", c.Addr, err)
		return
	}
}

func (c *httpClient) AddBot(bot *client.Client) {
	if c.Bot != nil {
		return
	}
	c.Bot = bot
}
