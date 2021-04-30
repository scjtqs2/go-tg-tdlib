package entity

import (
	"encoding/json"
	"reflect"
)

type MSG map[string]interface{}

func (m MSG) ToJson() string {
	b, _ := json.Marshal(m)
	return string(b)
}

func OK(data interface{}) MSG {
	return MSG{"data": data, "retcode": 0, "status": "ok"}
}

func Failed(code int) MSG {
	return MSG{"data": nil, "retcode": code, "status": "failed"}
}


// MakeMsg struct转成MSG 用于发送
func MakeMsg(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

const (
	TEXT                           = "messageText"
	ANIMATION                      = "messageAnimation"
	AUDIO                          = "messageAudio"
	DOCUMENT                       = "messageDocument"
	PHOTO                          = "messagePhoto"
	EXPIRED_PHOTO                  = "messageExpiredPhoto"
	STICKER                        = "messageSticker"
	VIDEO                          = "messageVideo"
	EXPIRED_VIDEO                  = "messageExpiredVideo"
	VIDEO_NOTE                     = "messageVideoNote"
	VOICE_NOTE                     = "messageVoiceNote"
	LOCATION                       = "messageLocation"
	VENUE                          = "messageVenue"
	CONTACT                        = "messageContact"
	DICE                           = "messageDice"
	GAME                           = "messageGame"
	POOL                           = "messagePoll"
	INVOICE                        = "messageInvoice"
	CALL                           = "messageCall"
	VOICE_CHAT_STARTED             = "messageVoiceChatStarted"
	VOICE_CHAT_ENDED               = "messageVoiceChatEnded"
	INVITE_VOICE_CHAT_PARTICIPANTS = "messageInviteVoiceChatParticipants"
	BASIC_GROUP_CHAT_CREATE        = "messageBasicGroupChatCreate"
	SUPERGROUP_CHAT_CREATE         = "messageSupergroupChatCreate"
	CHAT_CHANGE_TITLE              = "messageChatChangeTitle"
	CHAT_CHANGE_PHOTO              = "messageChatChangePhoto"
	CHAT_DELETE_PHOTO              = "messageChatDeletePhoto"
	CHAT_ADD_MEMBERS               = "messageChatAddMembers"
	CHAT_JOIN_BY_LINK              = "messageChatJoinByLink"
	CHAT_DELETE_MEMBER             = "messageChatDeleteMember"
	CHAT_UPGRADE_TO                = "messageChatUpgradeTo"
	CHAT_UPGRADE_FORM              = "messageChatUpgradeFrom"
	PIN_MESSAGE                    = "messagePinMessage"
	SCREENSHOT_TAKEN               = "messageScreenshotTaken"
	CHAT_SET_TTL                   = "messageChatSetTTL"
	CUSTOM_SERVICE_ACTION          = "messageCustomServiceAction"
	GAME_SCORE                     = "messageGameScore"
	PAYMENT_SUCCESSFUL             = "messagePaymentSuccessful"
	PAYMENT_SUCCESSFUL_BOT         = "messagePaymentSuccessfulBot"
	CONTACT_REGISTERED             = "messageContactRegistered"
	WEBSITED_CONNECTED             = "messageWebsiteConnected"
	PASSPORT_DATA_SENT             = "messagePassportDataSent"
	PASSPORT_DATA_RECEIVED         = "messagePassportDataReceived"
	PROXIMIT_ALERT_TIGGERED        = "messageProximityAlertTriggered"
	UNSUPPORTED                    = "messageUnsupported"
)
