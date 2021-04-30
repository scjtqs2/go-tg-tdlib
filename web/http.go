package web

import (
	"github.com/Arman92/go-tdlib"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/scjtqs/go-tg/entity"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
)

var HttpServer = &httpServer{}

type httpServer struct {
	engine *gin.Engine
	bot    *tdlib.Client
}

func (s *httpServer) Run(addr, authToken string, bot *tdlib.Client) {
	gin.SetMode(gin.ReleaseMode)
	s.engine = gin.New()
	s.bot = bot
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

	//s.engine.Any("/get_login_info", s.GetLoginInfo)
	//s.engine.Any("/get_login_info_async", s.GetLoginInfo)
	//
	//s.engine.Any("/get_friend_list", s.GetFriendList)
	//s.engine.Any("/get_friend_list_async", s.GetFriendList)
	//
	//s.engine.Any("/get_group_list", s.GetGroupList)
	//s.engine.Any("/get_group_list_async", s.GetGroupList)
	//
	//s.engine.Any("/get_group_info", s.GetGroupInfo)
	//s.engine.Any("/get_group_info_async", s.GetGroupInfo)
	//
	//s.engine.Any("/get_group_member_list", s.GetGroupMemberList)
	//s.engine.Any("/get_group_member_list_async", s.GetGroupMemberList)
	//
	//s.engine.Any("/get_group_member_info", s.GetGroupMemberInfo)
	//s.engine.Any("/get_group_member_info_async", s.GetGroupMemberInfo)
	//
	// 通过 chatid发送消息
	s.engine.Any("/send_msg", s.SendMessage)
	// 通过用户名获取chatid
	s.engine.Any("/get_chat_info", s.GetChatInfo)
	//
	//s.engine.Any("/send_private_msg", s.SendPrivateMessage)
	//s.engine.Any("/send_private_msg_async", s.SendPrivateMessage)
	//
	//s.engine.Any("/send_group_msg", s.SendGroupMessage)
	//s.engine.Any("/send_group_msg_async", s.SendGroupMessage)
	//
	//s.engine.Any("/get_status", s.GetStatus)
	//s.engine.Any("/get_status_async", s.GetStatus)

	go func() {
		log.Infof("go-tg HTTP 服务器已启动: %v", addr)
		log.Fatal(s.engine.Run(addr))
	}()
}

// SendMessage 发送信息
func (s *httpServer) SendMessage(c *gin.Context) {
	chatID, _ := strconv.ParseInt(getParam(c, "chat_id"), 10, 64)
	message, t := getParamWithType(c, "message")
	if t == gjson.JSON {
		inputMsg, err := s.makeMsg(message)
		if err != nil {
			c.JSON(200, entity.Failed(404))
			return
		}
		msg, err := s.bot.SendMessage(chatID, int64(0), int64(0), nil, nil, inputMsg)
		if err != nil {
			//消息发送失败
			c.JSON(200, entity.Failed(400))
			return
		}
		c.JSON(200, entity.OK(msg))
		return
	}
	c.JSON(404, entity.Failed(404))
}

// GetChatInfo 通过 名称获取 chat信息
func (s *httpServer) GetChatInfo(c *gin.Context) {
	name := getParam(c, "name")
	if name == "" {
		c.JSON(400, entity.Failed(400))
		return
	}
	chat, err := s.bot.SearchPublicChat(name)
	if err != nil {
		c.JSON(200, entity.Failed(400))
		return
	}
	c.JSON(200, entity.OK(chat))
}

func (s *httpServer) makeMsg(message string) (tdlib.InputMessageContent, error) {
	var inputMsg tdlib.InputMessageContent
	msg := gjson.Parse(message)
	switch msg.Get("msgtype").String() {
	case entity.TEXT:
		inputMsg = tdlib.NewInputMessageText(tdlib.NewFormattedText(msg.Get("content").String(), nil), true, true)
	case entity.PHOTO:
		//inputMsg = tdlib.NewInputMessagePhoto(tdlib.NewInputFileLocal("./bunny.jpg"), nil, nil, 400, 400,
		//	tdlib.NewFormattedText("A photo sent from go-tdlib!", nil), 0)
		return nil, errors.New("not support msg")
	default:
		return nil, errors.New("not support msg")
	}
	return inputMsg, nil
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
