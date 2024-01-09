package app

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/scjtqs/go-tg/config"
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
	"strconv"
)

type AppClient struct {
	Cli  *client.Client
	Conf *config.JsonConfig
	Cron *cron.Cron
}

// NewClient 初始化 bot方法
func NewClient(conf *config.JsonConfig) *AppClient {
	// client authorizer
	authorizer := client.ClientAuthorizer()
	// go client.CliInteractor(authorizer)
	go func() {
		for {
			select {
			case currentState, ok := <-authorizer.State:
				if !ok {
					// log.Errorf("get state not readdy")
					return
				}
				switch currentState.AuthorizationStateType() {
				case client.TypeAuthorizationStateWaitPhoneNumber:
					fmt.Print("Enter phone number: ")
					var phoneNumber string
					fmt.Scanln(&phoneNumber)

					authorizer.PhoneNumber <- phoneNumber

				case client.TypeAuthorizationStateWaitCode:
					var code string

					fmt.Print("Enter code: ")
					fmt.Scanln(&code)

					authorizer.Code <- code

				case client.TypeAuthorizationStateWaitPassword:
					fmt.Print("Enter password: ")
					var password string
					fmt.Scanln(&password)

					authorizer.Password <- password

				case client.TypeAuthorizationStateReady:
					log.Info("Authorization Ready! Let's rock")
					return
				}
			}
		}
	}()
	appid, _ := strconv.ParseInt(conf.AppID, 10, 32)
	authorizer.TdlibParameters <- &client.SetTdlibParametersRequest{
		UseTestDc:              false,
		DatabaseDirectory:      conf.DatabaseDirectory,
		FilesDirectory:         conf.FileDirectory,
		UseFileDatabase:        conf.UseFileDatabase,
		UseChatInfoDatabase:    conf.UseChatInfoDatabase,
		UseMessageDatabase:     conf.UseMessageDatabase,
		UseSecretChats:         false,
		ApiId:                  int32(appid),
		ApiHash:                conf.AppHash,
		SystemLanguageCode:     "en",
		DeviceModel:            "Docker_Server",
		SystemVersion:          "1.8.0",
		ApplicationVersion:     "1.8.0",
		EnableStorageOptimizer: true,
		IgnoreFileNames:        conf.IgnoreFileNames,
	}

	// You can set user-name and password to empty of don't need it
	// Socks5
	// client.AddProxy("pi.scjtqs.com", 1234, true, tdlib.NewProxyTypeSocks5("user-name", "password"))
	// client.AddProxy("pi.scjtqs.com", 10808, true, tdlib.NewProxyTypeSocks5("", ""))
	// HTTP - HTTPS proxy
	// client.AddProxy("127.0.0.1", 1234, true, tdlib.NewProxyTypeHttp("user-name", "password", false))
	// MtProto Proxy
	// client.AddProxy("127.0.0.1", 1234, true, tdlib.NewProxyTypeMtproto("MTPROTO-SECRET"))
	var proxy client.Option
	if conf.Proxy.ProxyStatus {
		switch conf.Proxy.ProxyType {
		case "Socks5":
			port, err := strconv.ParseInt(conf.Proxy.ProxyPort, 10, 32)
			if err != nil {
				panic("error proxy port")
			}
			proxy = client.WithProxy(&client.AddProxyRequest{
				Server: conf.Proxy.ProxyAddr,
				Port:   int32(port),
				Enable: true,
				Type: &client.ProxyTypeSocks5{
					Username: conf.Proxy.ProxyUser,
					Password: conf.Proxy.ProxyPasswd,
				},
			})
		case "HTTP":
			port, err := strconv.ParseInt(conf.Proxy.ProxyPort, 10, 32)
			if err != nil {
				panic("error proxy port")
			}
			proxy = client.WithProxy(&client.AddProxyRequest{
				Server: conf.Proxy.ProxyAddr,
				Port:   int32(port),
				Enable: true,
				Type: &client.ProxyTypeHttp{
					Username: conf.Proxy.ProxyUser,
					Password: conf.Proxy.ProxyPasswd,
				},
			})
		case "HTTPS":
			port, err := strconv.ParseInt(conf.Proxy.ProxyPort, 10, 32)
			if err != nil {
				panic("error proxy port")
			}
			proxy = client.WithProxy(&client.AddProxyRequest{
				Server: conf.Proxy.ProxyAddr,
				Port:   int32(port),
				Enable: true,
				Type: &client.ProxyTypeHttp{
					Username: conf.Proxy.ProxyUser,
					Password: conf.Proxy.ProxyPasswd,
				},
			})
		case "MtProto":
			port, err := strconv.ParseInt(conf.Proxy.ProxyPort, 10, 32)
			if err != nil {
				panic("error proxy port")
			}
			proxy = client.WithProxy(&client.AddProxyRequest{
				Server: conf.Proxy.ProxyAddr,
				Port:   int32(port),
				Enable: true,
				Type: &client.ProxyTypeMtproto{
					Secret: conf.Proxy.ProxyPasswd,
				},
			})
		default:
			log.Fatalf("proxyType error,only  'Socks5'、'HTTP'、'HTTPS'、'MtProto' supportd for proxyType ,proxyConf=%+v", conf.Proxy)
		}
	}

	logVerbosity := client.WithLogVerbosity(&client.SetLogVerbosityLevelRequest{
		NewVerbosityLevel: 0,
	})
	var (
		tdlibClient *client.Client
		err         error
	)
	switch proxy {
	case nil:
		tdlibClient, err = client.NewClient(authorizer, logVerbosity)
	default:
		tdlibClient, err = client.NewClient(authorizer, logVerbosity, proxy)
	}
	if err != nil {
		log.Fatalf("NewClient error: %s", err)
	}

	return &AppClient{
		Cli:  tdlibClient,
		Conf: conf,
	}
}

// SendMessageByName 给cron的定时发送使用 仅支持文本消息
func (a *AppClient) SendMessageByName(name string, message string) error {
	chat, err := a.Cli.SearchPublicChat(&client.SearchPublicChatRequest{
		Username: name,
	})
	if err != nil {
		log.Errorf("faild to check username %s,err:=%v \n\n", name, err)
		return err
	}
	chatID := chat.Id
	inputMsgTxt := &client.InputMessageText{
		Text: &client.FormattedText{Text: message},
		//DisableWebPagePreview: true,
		ClearDraft: true,
	}
	_, err = a.Cli.SendMessage(&client.SendMessageRequest{
		ChatId:              chatID,
		InputMessageContent: inputMsgTxt,
	})
	return err
}
