package app

import (
	"fmt"
	"github.com/Arman92/go-tdlib"
	log "github.com/sirupsen/logrus"
)

func (a *AppClient) FilterMsg()  {
	go func() {
		// Create an filter function which will be used to filter out unwanted tdlib messages
		eventFilter := func(msg *tdlib.TdMessage) bool {
			updateMsg := (*msg).(*tdlib.UpdateNewMessage)
			// For example, we want incomming messages from user with below id:
			//if updateMsg.Message.Sender {
			//	return true
			//}
			log.Infof("sender=%+v ",updateMsg.Message.Sender)
			return true
		}

		// Here we can add a receiver to retreive any message type we want
		// We like to get UpdateNewMessage events and with a specific FilterFunc
		receiver := a.Cli.AddEventReceiver(&tdlib.UpdateNewMessage{}, eventFilter, 5)
		for newMsg := range receiver.Chan {
			fmt.Println(newMsg)
			updateMsg := (newMsg).(*tdlib.UpdateNewMessage)
			// We assume the message content is simple text: (should be more sophisticated for general use)
			msgText := updateMsg.Message.Content.(*tdlib.MessageText)
			fmt.Println("MsgText:  ", msgText.Text)
			fmt.Print("\n\n")
		}

	}()
}
