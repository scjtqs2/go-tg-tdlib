package web

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/Arman92/go-tdlib"
	"github.com/gin-gonic/gin"
	"github.com/scjtqs/go-tg/entity"
	"github.com/scjtqs/go-tg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"path"
	"strconv"
	"strings"
)

// SendMessage 发送信息
func (s *httpServer) SendMessage(c *gin.Context) {
	chatID, _ := strconv.ParseInt(getParam(c, "chat_id"), 10, 64)
	message, t := getParamWithType(c, "message")
	log.Infof("sendMessage to chat_id=%v,message=%v", chatID, message)
	if t == gjson.JSON {
		inputMsg, err := s.makeMsg(message)
		if err != nil {
			log.Error(err)
			c.JSON(200, entity.Failed(404, err.Error()))
			return
		}
		msg, err := s.bot.SendMessage(chatID, int64(0), int64(0), nil, nil, inputMsg)
		if err != nil {
			log.Error(err)
			//消息发送失败
			c.JSON(200, entity.Failed(400, err.Error()))
			return
		}
		c.JSON(200, entity.OK(msg))
		return
	}
	log.Error("invalid json")
	c.JSON(404, entity.Failed(404, "invalid json"))
}

// GetChatInfo 通过 名称获取 chat信息
func (s *httpServer) GetChatInfo(c *gin.Context) {
	name := getParam(c, "name")
	if name == "" {
		log.Error("invalid name")
		c.JSON(400, entity.Failed(400, "invalid name"))
		return
	}
	chat, err := s.bot.SearchPublicChat(name)
	if err != nil {
		log.Error(err)
		c.JSON(200, entity.Failed(400, err.Error()))
		return
	}
	c.JSON(200, entity.OK(chat))
}

// GetUserInfo 获取当前用户信息
func (s *httpServer) GetUserInfo(c *gin.Context) {
	info, err := s.bot.GetMe()
	if err != nil {
		log.Error(err)
		c.JSON(400, entity.Failed(400, err.Error()))
		return
	}
	c.JSON(200, entity.OK(info))
}

// makeMsg 消息体构造
func (s *httpServer) makeMsg(message string) (tdlib.InputMessageContent, error) {
	var inputMsg tdlib.InputMessageContent
	msg := gjson.Parse(message)
	switch msg.Get("msgtype").String() {
	case entity.TEXT:
		inputMsg = tdlib.NewInputMessageText(tdlib.NewFormattedText(msg.Get("content").String(), nil), true, true)
	case entity.PHOTO:
		f := msg.Get("file").String()
		var filePath string
		if strings.HasPrefix(f, "http") || strings.HasPrefix(f, "https") {
			cache := msg.Get("cache").String()
			if cache == "" || !msg.Get("cache").Exists() {
				cache = "1"
			}
			hash := md5.Sum([]byte(f))
			cacheFile := path.Join("/tmp", hex.EncodeToString(hash[:])+".cache")
			if !utils.PathExists(cacheFile) || cache == "0" {
				b, err := utils.GetBytes(f)
				if err != nil {
					return nil, err
				}
				_ = ioutil.WriteFile(cacheFile, b, 0644)
			}
			filePath = cacheFile
		} else {
			filePath = f
		}
		log.Infof("send photo  file=%s,path=%s", f, filePath)
		inputMsg = tdlib.NewInputMessagePhoto(tdlib.NewInputFileLocal(filePath), nil, nil, 400, 400,
			tdlib.NewFormattedText(msg.Get("content").String(), nil), 0)
	default:
		return nil, errors.New("not support msg")
	}
	return inputMsg, nil
}

// SearchChatInfos 通过 名称获取搜索 chat信息
// query
func (s *httpServer) SearchChatInfos(c *gin.Context) {
	query := getParam(c, "query")
	if query == "" {
		c.JSON(400, entity.Failed(400, "invalid query"))
		return
	}
	chat, err := s.bot.SearchChatsOnServer(query, 50)
	if err != nil {
		log.Error(err)
		c.JSON(200, entity.Failed(400, err.Error()))
		return
	}
	c.JSON(200, entity.OK(chat))
}

// SearchChatInfos 通过 名称获取搜索 chat信息
// userID
func (s *httpServer) GetUserByUserId(c *gin.Context) {
	userID := getParam(c, "userID")
	if userID == "" {
		c.JSON(400, entity.Failed(400, "invalid userID"))
		return
	}
	uid, err := strconv.ParseInt(userID, 10, 32)
	if err != nil {
		log.Error(err)
		c.JSON(400, entity.Failed(400, "invalid userID"))
		return
	}
	chat, err := s.bot.GetUser(int32(uid))
	if err != nil {
		log.Error(err)
		c.JSON(200, entity.Failed(400, err.Error()))
		return
	}
	c.JSON(200, entity.OK(chat))
}

// GetMessage 获取消息
// chatID
// messageID
func (s *httpServer) GetMessage(c *gin.Context) {
	chatID := getParam(c, "chatID")
	if chatID == "" {
		c.JSON(400, entity.Failed(400, "invalid chatID"))
		return
	}
	cid, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		log.Error(err)
		c.JSON(400, entity.Failed(400, "invalid chatID"))
		return
	}
	messageID := getParam(c, "messageID")
	if messageID == "" {
		c.JSON(400, entity.Failed(400, "invalid messageID"))
		return
	}
	mid, err := strconv.ParseInt(messageID, 10, 64)
	if err != nil {
		log.Error(err)
		c.JSON(400, entity.Failed(400, "invalid messageID"))
		return
	}
	chat, err := s.bot.GetMessage(cid, mid)
	if err != nil {
		log.Error(err)
		c.JSON(200, entity.Failed(400, err.Error()))
		return
	}
	c.JSON(200, entity.OK(chat))
}
