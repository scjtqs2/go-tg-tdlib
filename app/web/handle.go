package web

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/scjtqs/go-tg/app/entity"
	"github.com/scjtqs/go-tg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/zelenin/go-tdlib/client"
	"io/ioutil"
	"path"
	"strconv"
	"strings"
)

var haveFullChatList bool
var allChats []*client.Chat

// SendMessage 发送信息
func (s *httpServer) SendMessage(c *gin.Context) {
	chatID, _ := strconv.ParseInt(getParam(c, "chat_id"), 10, 64)
	message, t := getParamWithType(c, "message")

	log.Debugf("sendMessage to chat_id=%v,message=%v", chatID, message)
	if t == gjson.JSON {
		messTraptIdstr, t := getParamWithType(c, "msg_trapt_id")
		var messTraptId int64
		var replay_id int64
		var err error
		if t != gjson.Number {
			messTraptId, err = strconv.ParseInt(messTraptIdstr, 10, 64)
			if err != nil {
				messTraptId = int64(0)
			}
		} else {
			messTraptId, _ = strconv.ParseInt(messTraptIdstr, 10, 64)
		}
		replay_id_str, t := getParamWithType(c, "msg_replay_id")
		if t != gjson.Number {
			replay_id, err = strconv.ParseInt(replay_id_str, 10, 64)
			if err != nil {
				replay_id = int64(0)
			}
		} else {
			replay_id, _ = strconv.ParseInt(replay_id_str, 10, 64)
		}
		inputMsg, err := s.makeMsg(message)
		if err != nil {
			log.Error(err)
			c.JSON(200, entity.Failed(404, err.Error()))
			return
		}
		msg, err := s.bot.SendMessage(&client.SendMessageRequest{
			ChatId:              chatID,
			MessageThreadId:     messTraptId,
			ReplyTo:             &client.InputMessageReplyToMessage{MessageId: replay_id},
			InputMessageContent: inputMsg,
		})
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
	chatID := getParam(c, "chat_id")
	if chatID != "" {
		chatid, _ := strconv.ParseInt(chatID, 10, 64)
		chat, err := s.bot.GetChat(&client.GetChatRequest{
			ChatId: chatid,
		})
		if err != nil {
			c.JSON(400, entity.Failed(400, err.Error()))
			return
		}
		c.JSON(200, entity.OK(chat))
		return
	}
	name := getParam(c, "name")
	log.Debugf("name=%s", name)
	if name == "" {
		log.Error("invalid name")
		c.JSON(400, entity.Failed(400, "invalid name"))
		return
	}
	chat, err := s.bot.SearchPublicChat(&client.SearchPublicChatRequest{
		Username: name,
	})
	if err != nil {
		log.Error(err)
		c.JSON(200, entity.Failed(400, err.Error()))
		return
	}
	log.Debug("get chat info=", chat)
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
	log.Debug("userinfo=", info)
	c.JSON(200, entity.OK(info))
}

// makeMsg 消息体构造
func (s *httpServer) makeMsg(message string) (client.InputMessageContent, error) {
	var inputMsg client.InputMessageContent
	msg := gjson.Parse(message)
	switch msg.Get("msgtype").String() {
	case entity.TEXT:
		inputMsg = &client.InputMessageText{
			Text: &client.FormattedText{
				Text: msg.Get("content").String(),
			},
			//DisableWebPagePreview: true,
			ClearDraft: true,
		}

	case entity.PHOTO:
		f := msg.Get("file").String()
		var filePath string
		if strings.HasPrefix(f, "http") || strings.HasPrefix(f, "https") {
			cache := msg.Get("cache").String()
			if cache == "" || !msg.Get("cache").Exists() {
				cache = "1"
			}
			hash := md5.Sum([]byte(f))
			cacheFile := path.Join("/tmp", hex.EncodeToString(hash[:])+".gif")
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
		var addedStickerFileIds []int32
		stickerFileId := msg.Get("sfid").Int()
		if stickerFileId == 0 {
			addedStickerFileIds = nil
		} else {
			addedStickerFileIds = []int32{int32(stickerFileId)}
		}
		log.Debugf("send photo  file=%s,path=%s", f, filePath)
		inputMsg = &client.InputMessagePhoto{
			Photo: &client.InputFileLocal{
				Path: filePath,
			},
			AddedStickerFileIds: addedStickerFileIds,
			Width:               400,
			Height:              400,
			Caption: &client.FormattedText{
				Text: msg.Get("content").String(),
			},
		}
	default:
		return nil, errors.New("not support msg")
	}
	return inputMsg, nil
}

// SearchChatInfos 通过 名称获取搜索 chat信息
// query
func (s *httpServer) SearchChatInfos(c *gin.Context) {
	query := getParam(c, "query")
	log.Debug("query=", query)
	if query == "" {
		c.JSON(400, entity.Failed(400, "invalid query"))
		return
	}
	chat, err := s.bot.SearchChatsOnServer(&client.SearchChatsOnServerRequest{Query: query, Limit: 50})
	if err != nil {
		log.Error(err)
		c.JSON(200, entity.Failed(400, err.Error()))
		return
	}
	log.Debug("SearchChatInfos chat=", chat)
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
	uid, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		log.Error(err)
		c.JSON(400, entity.Failed(400, "invalid userID"))
		return
	}
	user, err := s.bot.GetUser(&client.GetUserRequest{UserId: uid})
	if err != nil {
		log.Error(err)
		c.JSON(200, entity.Failed(400, err.Error()))
		return
	}
	log.Debug("GetUserByUserId,user=", user)
	c.JSON(200, entity.OK(user))
}

// GetMessage 获取消息
// chatID
// messageID
func (s *httpServer) GetMessage(c *gin.Context) {
	chatID := getParam(c, "chat_id")
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
	messageID := getParam(c, "message_id")
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
	msg, err := s.bot.GetMessage(&client.GetMessageRequest{
		ChatId:    cid,
		MessageId: mid,
	})
	if err != nil {
		log.Error(err)
		c.JSON(200, entity.Failed(400, err.Error()))
		return
	}
	log.Debug("GetMessage,message=", msg)
	c.JSON(200, entity.OK(msg))
}

// getChatList 获取聊天列表
// limit
func (s *httpServer) getChatList(c *gin.Context) {
	limit := getParam(c, "limit")
	if limit == "" {
		limit = "1000"
	}
	lid, err := strconv.ParseInt(limit, 10, 32)
	if err != nil {
		log.Error(err)
		c.JSON(400, entity.Failed(400, "invalid limit"))
		return
	}
	var chatList client.ChatList

	// get chats (ids) from tdlib
	chats, err := s.bot.GetChats(&client.GetChatsRequest{
		ChatList: chatList,
		Limit:    int32(lid),
	})
	if err != nil {
		log.Error(err)
		c.JSON(400, entity.Failed(400, err.Error()))
		return
	}
	c.JSON(200, entity.OK(chats))
}

// see https://stackoverflow.com/questions/37782348/how-to-use-getchats-in-tdlib
func getChatList(cli *client.Client, limit int) error {

	if !haveFullChatList && limit > len(allChats) {
		var chatList client.ChatList
		// get chats (ids) from tdlib
		chats, err := cli.GetChats(&client.GetChatsRequest{
			ChatList: chatList,
			Limit:    int32(limit - len(allChats)),
		})
		if err != nil {
			return err
		}
		if len(chats.ChatIds) == 0 {
			haveFullChatList = true
			return nil
		}

		for _, chatID := range chats.ChatIds {
			// get chat info from tdlib
			chat, err := cli.GetChat(&client.GetChatRequest{ChatId: chatID})
			if err == nil {
				allChats = append(allChats, chat)
			} else {
				return err
			}
		}
		return getChatList(cli, limit)
	}
	return nil
}

func (s *httpServer) GetMessagesByChatID(c *gin.Context) {
	chatID := getParam(c, "chat_id")
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
	messages, err := s.bot.GetMessages(&client.GetMessagesRequest{ChatId: cid})
	if err != nil {
		log.Error(err)
		c.JSON(400, entity.Failed(400, err.Error()))
		return
	}
	c.JSON(200, entity.OK(messages))
}
