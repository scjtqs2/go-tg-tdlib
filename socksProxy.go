package main

import (
	"fmt"
	"github.com/Arman92/go-tdlib"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	tdlib.SetLogVerbosityLevel(1)
	tdlib.SetFilePath("./errors.txt")

	// Create new instance of client
	client := tdlib.NewClient(tdlib.Config{
		APIID:               "187786",
		APIHash:             "e782045df67ba48e441ccb105da8fc85",
		SystemLanguageCode:  "en",
		DeviceModel:         "Server",
		SystemVersion:       "1.0.0",
		ApplicationVersion:  "1.0.0",
		UseMessageDatabase:  true,
		UseFileDatabase:     true,
		UseChatInfoDatabase: true,
		UseTestDataCenter:   false,
		DatabaseDirectory:   "./tdlib-db",
		FileDirectory:       "./tdlib-files",
		IgnoreFileNames:     false,
	})

	// You can set user-name and password to empty of don't need it
	// Socks5
	//client.AddProxy("pi.scjtqs.com", 1234, true, tdlib.NewProxyTypeSocks5("user-name", "password"))
	client.AddProxy("pi.scjtqs.com", 10808, true, tdlib.NewProxyTypeSocks5("", ""))
	// HTTP - HTTPS proxy
	//client.AddProxy("127.0.0.1", 1234, true, tdlib.NewProxyTypeHttp("user-name", "password", false))
	// MtProto Proxy
	//client.AddProxy("127.0.0.1", 1234, true, tdlib.NewProxyTypeMtproto("MTPROTO-SECRET"))

	// Handle Ctrl+C
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		client.DestroyInstance()
		os.Exit(1)
	}()

	//// Wait while we get AuthorizationReady!
	//// Note: See authorization example for complete auhtorization sequence example
	//currentState, _ := client.Authorize()
	//for ; currentState.GetAuthorizationStateEnum() != tdlib.AuthorizationStateReadyType; currentState, _ = client.Authorize() {
	//	time.Sleep(300 * time.Millisecond)
	//}
	//
	//// Send "/start" text every 5 seconds to Forsquare bot chat
	//go func() {
	//	// Should get chatID somehow, check out "getChats" example
	//	chatID := int64(198529620) // Foursquare bot chat id
	//
	//	inputMsgTxt := tdlib.NewInputMessageText(tdlib.NewFormattedText("/start", nil), true, true)
	//	client.SendMessage(chatID, 0, 0, nil, nil, inputMsgTxt)
	//
	//	time.Sleep(5 * time.Second)
	//}()

	for {
		currentState, _ := client.Authorize()
		if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPhoneNumberType {
			fmt.Print("Enter phone: ")
			var number string
			fmt.Scanln(&number)
			_, err := client.SendPhoneNumber(number)
			if err != nil {
				fmt.Printf("Error sending phone number: %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitCodeType {
			fmt.Print("Enter code: ")
			var code string
			fmt.Scanln(&code)
			_, err := client.SendAuthCode(code)
			if err != nil {
				fmt.Printf("Error sending auth code : %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPasswordType {
			fmt.Print("Enter Password: ")
			var password string
			fmt.Scanln(&password)
			_, err := client.SendAuthPassword(password)
			if err != nil {
				fmt.Printf("Error sending auth password: %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateReadyType {
			fmt.Println("Authorization Ready! Let's rock")
			break
		}
	}

	// rawUpdates gets all updates comming from tdlib
	rawUpdates := client.GetRawUpdatesChannel(100)
	for update := range rawUpdates {
		// Show all updates
		fmt.Println(update.Data)
		fmt.Print("\n\n")
	}

}
