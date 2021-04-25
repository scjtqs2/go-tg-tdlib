package entity

type TextMsg struct {
	UserID int32 `json:"userid"`
	Text string `json:"text"`
	ChatID int64 `json:"chatid"`
	SenderChatID int64 `json:"sender_chatid"`
}
