package app

import (
	"context"
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
	appid, _ := strconv.ParseInt(conf.AppID, 10, 32)
	tdlibParameters := &client.SetTdlibParametersRequest{
		UseTestDc:           false,
		DatabaseDirectory:   conf.DatabaseDirectory,
		FilesDirectory:      conf.FileDirectory,
		UseFileDatabase:     conf.UseFileDatabase,
		UseChatInfoDatabase: conf.UseChatInfoDatabase,
		UseMessageDatabase:  conf.UseMessageDatabase,
		UseSecretChats:      false,
		ApiId:               int32(appid),
		ApiHash:             conf.AppHash,
		SystemLanguageCode:  "en",
		DeviceModel:         "Docker_Server",
		SystemVersion:       "1.8.0",
		ApplicationVersion:  "1.8.0",
		// EnableStorageOptimizer: true,
		// IgnoreFileNames:        conf.IgnoreFileNames,
	}

	// client authorizer
	authorizer := client.ClientAuthorizer(tdlibParameters)
	// go client.CliInteractor(authorizer)
	go func() {
		for {
			select {
			case currentState, ok := <-authorizer.State:
				if !ok {
					// log.Errorf("get state not readdy")
					return
				}
				switch currentState.AuthorizationStateConstructor() {
				case client.ConstructorAuthorizationStateWaitPhoneNumber:
					fmt.Print("Enter phone number: ")
					var phoneNumber string
					fmt.Scanln(&phoneNumber)

					authorizer.PhoneNumber <- phoneNumber

				case client.ConstructorAuthorizationStateWaitCode:
					var code string

					fmt.Print("Enter code: ")
					fmt.Scanln(&code)

					authorizer.Code <- code

				case client.ConstructorAuthorizationStateWaitPassword:
					fmt.Print("Enter password: ")
					var password string
					fmt.Scanln(&password)

					authorizer.Password <- password

				case client.ConstructorAuthorizationStateReady:
					log.Info("Authorization Ready! Let's rock")
					return
				}
			}
		}
	}()

	var options []client.Option
	if conf.Proxy.ProxyStatus {
		switch conf.Proxy.ProxyType {
		case "Socks5":
			port, err := strconv.ParseInt(conf.Proxy.ProxyPort, 10, 32)
			if err != nil {
				panic("error proxy port")
			}
			options = append(options, client.WithProxy(&client.AddProxyRequest{
				Proxy: &client.Proxy{
					Server: conf.Proxy.ProxyAddr,
					Port:   int32(port),
					Type: &client.ProxyTypeSocks5{
						Username: conf.Proxy.ProxyUser,
						Password: conf.Proxy.ProxyPasswd,
					},
				},
				Enable: true,
			}))
		case "HTTP":
			port, err := strconv.ParseInt(conf.Proxy.ProxyPort, 10, 32)
			if err != nil {
				panic("error proxy port")
			}
			options = append(options, client.WithProxy(&client.AddProxyRequest{
				Proxy: &client.Proxy{
					Server: conf.Proxy.ProxyAddr,
					Port:   int32(port),
					Type: &client.ProxyTypeHttp{
						Username: conf.Proxy.ProxyUser,
						Password: conf.Proxy.ProxyPasswd,
					},
				},
				Enable: true,
			}))
		case "HTTPS":
			port, err := strconv.ParseInt(conf.Proxy.ProxyPort, 10, 32)
			if err != nil {
				panic("error proxy port")
			}
			options = append(options, client.WithProxy(&client.AddProxyRequest{
				Proxy: &client.Proxy{
					Server: conf.Proxy.ProxyAddr,
					Port:   int32(port),
					Type: &client.ProxyTypeHttp{
						Username: conf.Proxy.ProxyUser,
						Password: conf.Proxy.ProxyPasswd,
					},
				},
				Enable: true,
			}))
		case "MtProto":
			port, err := strconv.ParseInt(conf.Proxy.ProxyPort, 10, 32)
			if err != nil {
				panic("error proxy port")
			}
			options = append(options, client.WithProxy(&client.AddProxyRequest{
				Proxy: &client.Proxy{
					Server: conf.Proxy.ProxyAddr,
					Port:   int32(port),
					Type: &client.ProxyTypeMtproto{
						Secret: conf.Proxy.ProxyPasswd,
					},
				},
				Enable: true,
			}))
		default:
			log.Fatalf("proxyType error,only  'Socks5'、'HTTP'、'HTTPS'、'MtProto' supportd for proxyType ,proxyConf=%+v", conf.Proxy)
		}
	}

	_, _ = client.SetLogVerbosityLevel(&client.SetLogVerbosityLevelRequest{
		NewVerbosityLevel: 0,
	})

	tdlibClient, err := client.NewClient(authorizer, options...)
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
	chat, err := a.Cli.SearchPublicChat(context.Background(), &client.SearchPublicChatRequest{
		Username: name,
	})
	if err != nil {
		log.Errorf("faild to check username %s,err:=%v \n\n", name, err)
		return err
	}
	chatID := chat.Id
	inputMsgTxt := &client.InputMessageText{
		Text: &client.FormattedText{Text: message},
		// DisableWebPagePreview: true,
		ClearDraft: true,
	}
	_, err = a.Cli.SendMessage(context.Background(), &client.SendMessageRequest{
		ChatId:              chatID,
		InputMessageContent: inputMsgTxt,
	})
	return err
}
