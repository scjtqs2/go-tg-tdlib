package app

import (
	"fmt"
	"github.com/Arman92/go-tdlib/v2/client"
	"github.com/Arman92/go-tdlib/v2/tdlib"
	"github.com/robfig/cron/v3"
	"github.com/scjtqs/go-tg/config"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type AppClient struct {
	Cli  *client.Client
	Conf *config.JsonConfig
	Cron *cron.Cron
}

// NewClient 初始化 bot方法
func NewClient(conf *config.JsonConfig) *AppClient {
	// Create new instance of client
	client := client.NewClient(client.Config{
		APIID:                  conf.AppID,
		APIHash:                conf.AppHash,
		SystemLanguageCode:     "en",
		DeviceModel:            "Docker_Server",
		SystemVersion:          "1.0.0",
		ApplicationVersion:     "1.0.0",
		UseMessageDatabase:     conf.UseMessageDatabase,
		UseFileDatabase:        conf.UseFileDatabase,
		UseChatInfoDatabase:    conf.UseChatInfoDatabase,
		UseTestDataCenter:      conf.UseTestDataCenter,
		DatabaseDirectory:      conf.DatabaseDirectory,
		FileDirectory:          conf.FileDirectory,
		IgnoreFileNames:        conf.IgnoreFileNames,
		UseSecretChats:         true,
		EnableStorageOptimizer: true,
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
		default:
			log.Fatalf("proxyType error,only  'Socks5'、'HTTP'、'HTTPS'、'MtProto' supportd for proxyType ,proxyConf=%+v", conf.Proxy)
		}
	}

	//// Wait while we get AuthorizationReady!
	//// Note: See authorization example for complete auhtorization sequence example
	for {
		currentState, _ := client.Authorize()
		if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPhoneNumberType {
			fmt.Print("Enter phone: ")
			var number string
			fmt.Scanln(&number)
			conf.Phone = number
			_, err := client.SendPhoneNumber(number)
			if err != nil {
				log.Errorf("Error sending phone number: %v", err)
			}
			conf.Save(config.ConfigPath)
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitCodeType {
			fmt.Print("Enter code: ")
			var code string
			fmt.Scanln(&code)
			_, err := client.SendAuthCode(code)
			if err != nil {
				log.Errorf("Error sending auth code : %v ", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPasswordType {
			fmt.Print("Enter Password: ")
			var password string
			fmt.Scanln(&password)
			conf.Password = password
			_, err := client.SendAuthPassword(password)
			if err != nil {
				log.Errorf("Error sending auth password: %v", err)
			}
			conf.Save(config.ConfigPath)
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateReadyType {
			log.Info("Authorization Ready! Let's rock")
			break
		}
	}
	return &AppClient{
		Cli:  client,
		Conf: conf,
	}
}

// SendMessageByName 给cron的定时发送使用 仅支持文本消息
func (a *AppClient) SendMessageByName(name string, message string) error {
	chat, err := a.Cli.SearchPublicChat(name)
	if err != nil {
		log.Errorf("faild to check username %s,err:=%v \n\n", name, err)
		return err
	}
	chatID := chat.ID
	inputMsgTxt := tdlib.NewInputMessageText(tdlib.NewFormattedText(message, nil), true, true)
	_, err = a.Cli.SendMessage(chatID, int64(0), int64(0), nil, nil, inputMsgTxt)
	return err
}
