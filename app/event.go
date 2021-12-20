package app

import (
	"github.com/Arman92/go-tdlib/v2/tdlib"
	"github.com/scjtqs/go-tg/app/entity"
	"github.com/scjtqs/go-tg/app/webhook"
	"github.com/scjtqs/go-tg/config"
	"github.com/scjtqs/go-tg/utils"
	log "github.com/sirupsen/logrus"
	"strconv"
)

func (a *AppClient) FilterMsg(index int, conf *config.WebHook) {
	if !conf.WebHookStatus {
		return
	}
	var chatIds []string
	if conf.WebHookFilter != nil && conf.WebHookFilter.FilterStatus {
		for _, v := range conf.WebHookFilter.FilterNames {
			id, err := a.Cli.SearchPublicChat(v)
			if err != nil {
				log.Errorf("search for chat faild name=%s,err=%v", v, err)
				continue
			}
			chatIds = append(chatIds, strconv.FormatInt(id.ID, 10))
		}
	}

	go func() {
		// Create an filter function which will be used to filter out unwanted tdlib messages
		eventFilter := func(msg *tdlib.TdMessage) bool {
			updateMsg := (*msg).(*tdlib.UpdateNewMessage)
			switch updateMsg.Message.Sender.GetMessageSenderEnum() {
			case tdlib.MessageSenderUserType:
				sender := updateMsg.Message.Sender.(*tdlib.MessageSenderUser)
				//log.Debugf("senderUser=%+v", sender)
				if utils.InArrayString(strconv.FormatInt(int64(sender.UserID), 10), chatIds) {
					return true
				}
			case tdlib.MessageSenderChatType:
				sender := updateMsg.Message.Sender.(*tdlib.MessageSenderChat)
				//log.Debugf("senderChat=%+v", sender)
				if utils.InArrayString(strconv.FormatInt(int64(sender.ChatID), 10), chatIds) {
					return true
				}
			}
			return !conf.WebHookFilter.FilterStatus
		}

		// Here we can add a receiver to retreive any message type we want
		// We like to get UpdateNewMessage events and with a specific FilterFunc
		receiver := a.Cli.AddEventReceiver(&tdlib.UpdateNewMessage{}, eventFilter, 5)
		for newMsg := range receiver.Chan {
			updateMsg := (newMsg).(*tdlib.UpdateNewMessage)
			// We assume the message content is simple text: (should be more sophisticated for general use)
			switch updateMsg.Message.Content.GetMessageContentEnum() {
			case "messageText":
				msgText := updateMsg.Message.Content.(*tdlib.MessageText)
				log.Debugf("msgText: %s", msgText.Text.Text)
				a.makeMsgPush(index, "messageText", updateMsg.Message)
			case "messageAnimation":
				a.makeMsgPush(index, "messageAnimation", updateMsg.Message)
			case "messageAudio":
				audio := updateMsg.Message.Content.(*tdlib.MessageAudio)
				log.Debugf("audio: %+v", audio)
				a.makeMsgPush(index, "messageAudio", updateMsg.Message)
			case "messageDocument":
				doc := updateMsg.Message.Content.(*tdlib.MessageDocument)
				log.Debugf("document:%+v", doc)
				a.makeMsgPush(index, "messageDocument", updateMsg.Message)
			case "messagePhoto":
				photo := updateMsg.Message.Content.(*tdlib.MessagePhoto)
				log.Debugf("photo:%+v", photo)
				a.makeMsgPush(index, "messagePhoto", updateMsg.Message)
			case "messageExpiredPhoto":
				expPhoto := updateMsg.Message.Content.(*tdlib.MessageExpiredPhoto)
				log.Debugf("messageExpiredPhoto:%+v", expPhoto)
				a.makeMsgPush(index, "messageExpiredPhoto", updateMsg.Message)
			case "messageSticker":
				a.makeMsgPush(index, "messageSticker", updateMsg.Message)
			case "messageVideo":
				video := updateMsg.Message.Content.(*tdlib.MessageVideo)
				log.Debugf("video:%+v", video)
				a.makeMsgPush(index, "messageVideo", updateMsg.Message)
			case "messageExpiredVideo":
				a.makeMsgPush(index, "messageExpiredVideo", updateMsg.Message)
			case "messageVideoNote":
				a.makeMsgPush(index, "messageVideoNote", updateMsg.Message)
			case "messageVoiceNote":
				a.makeMsgPush(index, "messageVoiceNote", updateMsg.Message)
			case "messageLocation":
				a.makeMsgPush(index, "messageLocation", updateMsg.Message)
			case "messageVenue":
				a.makeMsgPush(index, "messageVenue", updateMsg.Message)
			case "messageContact":
				a.makeMsgPush(index, "messageContact", updateMsg.Message)
			case "messageDice":
				a.makeMsgPush(index, "messageDice", updateMsg.Message)
			case "messageGame":
				a.makeMsgPush(index, "messageGame", updateMsg.Message)
			case "messagePoll":
				a.makeMsgPush(index, "messagePoll", updateMsg.Message)
			case "messageInvoice":
				a.makeMsgPush(index, "messageInvoice", updateMsg.Message)
			case "messageCall":
				a.makeMsgPush(index, "messageCall", updateMsg.Message)
			case "messageVoiceChatStarted":
				a.makeMsgPush(index, "messageVoiceChatStarted", updateMsg.Message)
			case "messageVoiceChatEnded":
				a.makeMsgPush(index, "messageVoiceChatEnded", updateMsg.Message)
			case "messageInviteVoiceChatParticipants":
				a.makeMsgPush(index, "messageInviteVoiceChatParticipants", updateMsg.Message)
			case "messageBasicGroupChatCreate":
				a.makeMsgPush(index, "messageBasicGroupChatCreate", updateMsg.Message)
			case "messageSupergroupChatCreate":
				a.makeMsgPush(index, "messageSupergroupChatCreate", updateMsg.Message)
			case "messageChatChangeTitle":
				a.makeMsgPush(index, "messageChatChangeTitle", updateMsg.Message)
			case "messageChatChangePhoto":
				a.makeMsgPush(index, "messageChatChangePhoto", updateMsg.Message)
			case "messageChatDeletePhoto":
				a.makeMsgPush(index, "messageChatDeletePhoto", updateMsg.Message)
			case "messageChatAddMembers":
				a.makeMsgPush(index, "messageChatAddMembers", updateMsg.Message)
			case "messageChatJoinByLink":
				a.makeMsgPush(index, "messageChatJoinByLink", updateMsg.Message)
			case "messageChatDeleteMember":
				a.makeMsgPush(index, "messageChatDeleteMember", updateMsg.Message)
			case "messageChatUpgradeTo":
				a.makeMsgPush(index, "messageChatUpgradeTo", updateMsg.Message)
			case "messageChatUpgradeFrom":
				a.makeMsgPush(index, "messageChatUpgradeFrom", updateMsg.Message)
			case "messagePinMessage":
				a.makeMsgPush(index, "messagePinMessage", updateMsg.Message)
			case "messageScreenshotTaken":
				a.makeMsgPush(index, "messageScreenshotTaken", updateMsg.Message)
			case "messageChatSetTTL":
				a.makeMsgPush(index, "messageChatSetTTL", updateMsg.Message)
			case "messageCustomServiceAction":
				a.makeMsgPush(index, "messageCustomServiceAction", updateMsg.Message)
			case "messageGameScore":
				a.makeMsgPush(index, "messageGameScore", updateMsg.Message)
			case "messagePaymentSuccessful":
				a.makeMsgPush(index, "messagePaymentSuccessful", updateMsg.Message)
			case "messagePaymentSuccessfulBot":
				a.makeMsgPush(index, "messagePaymentSuccessfulBot", updateMsg.Message)
			case "messageContactRegistered":
				a.makeMsgPush(index, "messageContactRegistered", updateMsg.Message)
			case "messageWebsiteConnected":
				a.makeMsgPush(index, "messageWebsiteConnected", updateMsg.Message)
			case "messagePassportDataSent":
				a.makeMsgPush(index, "messagePassportDataSent", updateMsg.Message)
			case "messagePassportDataReceived":
				a.makeMsgPush(index, "messagePassportDataReceived", updateMsg.Message)
			case "messageProximityAlertTriggered":
				a.makeMsgPush(index, "messageProximityAlertTriggered", updateMsg.Message)
			case "messageUnsupported":
				a.makeMsgPush(index, "messageUnsupported", updateMsg.Message)
			}
		}
	}()
}

// makeMsgPush 处理消息推送
func (a *AppClient) makeMsgPush(index int, msgType string, message *tdlib.Message) {
	switch message.Sender.GetMessageSenderEnum() {
	case tdlib.MessageSenderUserType:
		sender := message.Sender.(*tdlib.MessageSenderUser)
		webhook.AddMsg(index, &entity.MSG{
			"msgType":    msgType,
			"msg":        message,
			"userID":     sender.UserID,
			"senderType": message.Sender.GetMessageSenderEnum(),
		}, a.Cli)
	case tdlib.MessageSenderChatType:
		sender := message.Sender.(*tdlib.MessageSenderChat)
		webhook.AddMsg(index, &entity.MSG{
			"msgType":      msgType,
			"msg":          message,
			"senderType":   message.Sender.GetMessageSenderEnum(),
			"senderChatID": sender.ChatID,
		}, a.Cli)
	}
}
