package app

import (
	"github.com/Arman92/go-tdlib"
	"github.com/scjtqs/go-tg/config"
	"github.com/scjtqs/go-tg/entity"
	"github.com/scjtqs/go-tg/utils"
	"github.com/scjtqs/go-tg/webhook"
	log "github.com/sirupsen/logrus"
	"strconv"
)

func (a *AppClient) FilterMsg(index int,conf *config.WebHook) {
	if !conf.WebHookStatus {
		return
	}
	var chatIds []string
	if conf.WebHookFilter.FilterStatus {
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
				switch updateMsg.Message.Sender.GetMessageSenderEnum() {
				case tdlib.MessageSenderUserType:
					sender := updateMsg.Message.Sender.(*tdlib.MessageSenderUser)
					webhook.AddMsg(index,entity.TextMsg{UserID: sender.UserID, ChatID: updateMsg.Message.ChatID, Text: msgText.Text.Text})
				case tdlib.MessageSenderChatType:
					sender := updateMsg.Message.Sender.(*tdlib.MessageSenderChat)
					webhook.AddMsg(index,entity.TextMsg{ChatID: updateMsg.Message.ChatID, Text: msgText.Text.Text,SenderChatID: sender.ChatID})
				}
			case "messageAnimation":
			case "messageAudio":
				audio := updateMsg.Message.Content.(*tdlib.MessageAudio)
				log.Debugf("audio: %+v", audio)
			case "messageDocument":
				doc := updateMsg.Message.Content.(*tdlib.MessageDocument)
				log.Debugf("document:%+v", doc)
			case "messagePhoto":
				photo := updateMsg.Message.Content.(*tdlib.MessagePhoto)
				log.Debugf("photo:%+v", photo)
			case "messageExpiredPhoto":
				expPhoto := updateMsg.Message.Content.(*tdlib.MessageExpiredPhoto)
				log.Debugf("messageExpiredPhoto:%+v", expPhoto)
			case "messageSticker":
			case "messageVideo":
				video := updateMsg.Message.Content.(*tdlib.MessageVideo)
				log.Debugf("video:%+v", video)
			case "messageExpiredVideo":
			case "messageVideoNote":
			case "messageVoiceNote":
			case "messageLocation":
			case "messageVenue":
			case "messageContact":
			case "messageDice":
			case "messageGame":
			case "messagePoll":
			case "messageInvoice":
			case "messageCall":
			case "messageVoiceChatStarted":
			case "messageVoiceChatEnded":
			case "messageInviteVoiceChatParticipants":
			case "messageBasicGroupChatCreate":
			case "messageSupergroupChatCreate":
			case "messageChatChangeTitle":
			case "messageChatChangePhoto":
			case "messageChatDeletePhoto":
			case "messageChatAddMembers":
			case "messageChatJoinByLink":
			case "messageChatDeleteMember":
			case "messageChatUpgradeTo":
			case "messageChatUpgradeFrom":
			case "messagePinMessage":
			case "messageScreenshotTaken":
			case "messageChatSetTTL":
			case "messageCustomServiceAction":
			case "messageGameScore":
			case "messagePaymentSuccessful":
			case "messagePaymentSuccessfulBot":
			case "messageContactRegistered":
			case "messageWebsiteConnected":
			case "messagePassportDataSent":
			case "messagePassportDataReceived":
			case "messageProximityAlertTriggered":
			case "messageUnsupported":

			}
		}
	}()
}
