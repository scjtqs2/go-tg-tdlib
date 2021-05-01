package app

import (
	"fmt"
	"github.com/Arman92/go-tdlib"
	"github.com/robfig/cron/v3"
	"github.com/scjtqs/go-tg/config"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type AppClient struct {
	Cli  *tdlib.Client
	Conf *config.JsonConfig
	Cron *cron.Cron
}

func NewClient(conf *config.JsonConfig) *AppClient {
	// Create new instance of client
	client := tdlib.NewClient(tdlib.Config{
		APIID:               conf.AppID,
		APIHash:             conf.AppHash,
		SystemLanguageCode:  "en",
		DeviceModel:         "Server",
		SystemVersion:       "1.0.0",
		ApplicationVersion:  "1.0.0",
		UseMessageDatabase:  conf.UseMessageDatabase,
		UseFileDatabase:     conf.UseFileDatabase,
		UseChatInfoDatabase: conf.UseChatInfoDatabase,
		UseTestDataCenter:   conf.UseTestDataCenter,
		DatabaseDirectory:   conf.DatabaseDirectory,
		FileDirectory:       conf.FileDirectory,
		IgnoreFileNames:     conf.IgnoreFileNames,
	})

	// You can set user-name and password to empty of don't need it
	// Socks5
	//client.AddProxy("pi.scjtqs.com", 1234, true, tdlib.NewProxyTypeSocks5("user-name", "password"))
	//client.AddProxy("pi.scjtqs.com", 10808, true, tdlib.NewProxyTypeSocks5("", ""))
	// HTTP - HTTPS proxy
	//client.AddProxy("127.0.0.1", 1234, true, tdlib.NewProxyTypeHttp("user-name", "password", false))
	// MtProto Proxy
	//client.AddProxy("127.0.0.1", 1234, true, tdlib.NewProxyTypeMtproto("MTPROTO-SECRET"))
	if conf.Proxy.ProxyStatus {
		switch conf.Proxy.ProxyType {
		case "Socks5":
			port, err := strconv.ParseInt(conf.Proxy.ProxyPort, 10, 32)
			if err != nil {
				panic("error proxy port")
			}
			client.AddProxy(conf.Proxy.ProxyAddr, int32(port), true, tdlib.NewProxyTypeSocks5(conf.Proxy.ProxyUser, conf.Proxy.ProxyPasswd))
		case "HTTP":
			port, err := strconv.ParseInt(conf.Proxy.ProxyPort, 10, 32)
			if err != nil {
				panic("error proxy port")
			}
			client.AddProxy(conf.Proxy.ProxyAddr, int32(port), true, tdlib.NewProxyTypeHttp(conf.Proxy.ProxyUser, conf.Proxy.ProxyPasswd, false))
		case "HTTPS":
			port, err := strconv.ParseInt(conf.Proxy.ProxyPort, 10, 32)
			if err != nil {
				panic("error proxy port")
			}
			client.AddProxy(conf.Proxy.ProxyAddr, int32(port), true, tdlib.NewProxyTypeHttp(conf.Proxy.ProxyUser, conf.Proxy.ProxyPasswd, false))
		case "MtProto":
			port, err := strconv.ParseInt(conf.Proxy.ProxyPort, 10, 32)
			if err != nil {
				panic("error proxy port")
			}
			client.AddProxy(conf.Proxy.ProxyAddr, int32(port), true, tdlib.NewProxyTypeMtproto(conf.Proxy.ProxyPasswd))
		}
	}

	//// Wait while we get AuthorizationReady!
	//// Note: See authorization example for complete auhtorization sequence example
	for {
		currentState, _ := client.Authorize()
		if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPhoneNumberType {
			log.Info("Enter phone: \n")
			var number string
			fmt.Scanln(&number)
			_, err := client.SendPhoneNumber(number)
			if err != nil {
				log.Infof("Error sending phone number: %v \n", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitCodeType {
			log.Info("Enter code: \n")
			var code string
			fmt.Scanln(&code)
			_, err := client.SendAuthCode(code)
			if err != nil {
				log.Infof("Error sending auth code : %v \n", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPasswordType {
			log.Info("Enter Password: \n")
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
	return &AppClient{
		Cli:  client,
		Conf: conf,
	}
}

func (a *AppClient) SendMessageByName(name string, message string) error {
	chat, err := a.Cli.SearchPublicChat(name)
	if err != nil {
		log.Errorf("faild to check username %s,err:=%v \n\n", name, err)
		return err
	}
	chatID := chat.ID
	inputMsgTxt := tdlib.NewInputMessageText(tdlib.NewFormattedText(message, nil), true, true)
	_,err=a.Cli.SendMessage(chatID, int64(0), int64(0), nil, nil, inputMsgTxt)
	return err
}
