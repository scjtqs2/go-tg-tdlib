package web

import (
	"github.com/Arman92/go-tdlib"
	"github.com/gin-gonic/gin"
	"github.com/scjtqs/go-tg/config"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
)

var HttpServer = &httpServer{}

type httpServer struct {
	engine *gin.Engine
	bot    *tdlib.Client
	conf   *config.JsonConfig
}

func (s *httpServer) Run(addr, authToken string, bot *tdlib.Client,conf *config.JsonConfig) {
	gin.SetMode(gin.ReleaseMode)
	s.engine = gin.New()
	s.bot = bot
	s.conf = conf
	s.engine.Use(func(c *gin.Context) {
		if c.Request.Method != "GET" && c.Request.Method != "POST" {
			log.Warnf("已拒绝客户端 %v 的请求: 方法错误", c.Request.RemoteAddr)
			c.Status(404)
			return
		}
		if c.Request.Method == "POST" && c.Request.Header.Get("Content-Type") == "application/json" {
			d, err := c.GetRawData()
			if err != nil {
				log.Warnf("获取请求 %v 的Body时出现错误: %v", c.Request.RequestURI, err)
				c.Status(400)
				return
			}
			if !gjson.ValidBytes(d) {
				log.Warnf("已拒绝客户端 %v 的请求: 非法Json", c.Request.RemoteAddr)
				c.Status(400)
				return
			}
			c.Set("json_body", gjson.ParseBytes(d))
		}
		c.Next()
	})

	if authToken != "" {
		s.engine.Use(func(c *gin.Context) {
			if auth := c.Request.Header.Get("Authorization"); auth != "" {
				if strings.SplitN(auth, " ", 2)[1] != authToken {
					c.AbortWithStatus(401)
					return
				}
			} else if c.Query("access_token") != authToken {
				c.AbortWithStatus(401)
				return
			} else {
				c.Next()
			}
		})
	}

	// 通过 chatid发送消息
	s.engine.Any("/send_msg", s.SendMessage)
	// 通过用户名获取chatid
	s.engine.Any("/get_chat_info", s.GetChatInfo)
	// 查询当前登录用户信息
	s.engine.Any("/getme",s.GetUserInfo)
	// 搜索chat
	s.engine.Any("/search_chat_infos",s.SearchChatInfos)
	// 通过用户id查询用户信息
	s.engine.Any("/get_userinfo_by_userid",s.GetUserByUserId)
	// 通过chatid和messageid拉取消息信息
	s.engine.Any("/get_message",s.GetMessage)
	// 获取聊天列表 chatlist
	s.engine.Any("/get_chat_list",s.getChatList)
	// GetMessagesByChatID
	s.engine.Any("get_messages",s.GetMessagesByChatID)

	go func() {
		log.Infof("go-tg HTTP 服务器已启动: %v", addr)
		log.Fatal(s.engine.Run(addr))
	}()
}


func getParam(c *gin.Context, k string) string {
	p, _ := getParamWithType(c, k)
	return p
}

func getParamWithType(c *gin.Context, k string) (string, gjson.Type) {
	if q := c.Query(k); q != "" {
		return q, gjson.Null
	}
	if c.Request.Method == "POST" {
		if h := c.Request.Header.Get("Content-Type"); h != "" {
			if h == "application/x-www-form-urlencoded" {
				if p, ok := c.GetPostForm(k); ok {
					return p, gjson.Null
				}
			}
			if h == "application/json" {
				if obj, ok := c.Get("json_body"); ok {
					res := obj.(gjson.Result).Get(k)
					if res.Exists() {
						switch res.Type {
						case gjson.JSON:
							return res.Raw, gjson.JSON
						case gjson.String:
							return res.Str, gjson.String
						case gjson.Number:
							return strconv.FormatInt(res.Int(), 10), gjson.Number // 似乎没有需要接受 float 类型的api
						case gjson.True:
							return "true", gjson.True
						case gjson.False:
							return "false", gjson.False
						}
					}
				}
			}
		}
	}
	return "", gjson.Null
}
