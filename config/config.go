package config

import (
	"encoding/json"
	"github.com/scjtqs/go-tg/utils"
	log "github.com/sirupsen/logrus"
	"os"
)

var cronPath = "cron.json"

type JsonConfig struct {
	Phone               string     `json:"phone_number"`
	Password            string     `json:"password"`
	AppID               string     `json:"appid"`
	AppHash             string     `json:"app_hash"`
	UseMessageDatabase  bool       `json:"use_message_database"`
	UseFileDatabase     bool       `json:"use_file_database"`
	UseChatInfoDatabase bool       `json:"use_chat_info_database"`
	UseTestDataCenter   bool       `json:"use_test_data_center"`
	DatabaseDirectory   string     `json:"database_directory"`
	FileDirectory       string     `json:"file_directory"`
	IgnoreFileNames     bool       `json:"ignore_file_name"`
	Proxy               *Proxy     `json:"proxy"`
	WebHook             []*WebHook `json:"webhook"`
	WebApi              *WebApi    `json:"webapi"`
}

type Proxy struct {
	ProxyStatus bool   `json:"status"` //开关
	ProxyType   string `json:"type"`   // "HTTP"，HTTPS，Socks5、MtProto
	ProxyAddr   string `json:"addr"`   // 10.0.0.1
	ProxyPort   string `json:"port"`   // 1234
	ProxyUser   string `json:"user"`   // user
	ProxyPasswd string `json:"passwd"` //passwd
}

// webhook
type WebHook struct {
	WebHookStatus bool          `json:"status"`        //开关
	WebHookUrl    string        `json:"http_post_url"` //post推送地址
	WebHookSecret string        `json:"secret"`        //推送签名校验用的secret
	WebHookFilter *WebHookFilter `json:"filter"`        //过滤仅这些用户推送
}

type WebHookFilter struct {
	FilterStatus bool     `json:"status"` //开关
	FilterNames  []string `json:"names"`  //过滤的用户名们
}

// api
type WebApi struct {
	WebApiStatus bool   `json:"status"`    //开关
	WebApiHost   string `json:"bind_addr"` //绑定地址
	WebApiPort   string `json:"port"`      //api端口地址
	WebApiToken  string `json:"token"`     //简易的api鉴权参数
}

func Load(p string) *JsonConfig {
	if !utils.PathExists(p) {
		log.Warnf("尝试加载配置文件 %v 失败: 文件不存在", p)
		return nil
	}
	c := JsonConfig{}
	err := json.Unmarshal([]byte(utils.ReadAllText(p)), &c)
	if err != nil {
		log.Warnf("尝试加载配置文件 %v 时出现错误: %v", p, err)
		return nil
	}
	return &c
}

func (c *JsonConfig) Save(p string) error {
	data, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}
	utils.WriteAllText(p, string(data))
	return nil
}

func DefaultConfig() *JsonConfig {
	conf := &JsonConfig{
		Phone:               "",
		Password:            "",
		AppID:               "187786",
		AppHash:             "e782045df67ba48e441ccb105da8fc85",
		UseMessageDatabase:  true,
		UseFileDatabase:     true,
		UseChatInfoDatabase: true,
		UseTestDataCenter:   false,
		DatabaseDirectory:   "./tdlib-db",
		FileDirectory:       "./tdlib-files",
		IgnoreFileNames:     false,
		Proxy: &Proxy{
			ProxyStatus: false,
			ProxyType:   "Socks5",
			ProxyAddr:   "vpn.abc.com",
			ProxyPort:   "1234",
			ProxyUser:   "",
			ProxyPasswd: "", //mtp 的secret也是这个字段
		},
		WebHook: []*WebHook{
			&WebHook{
				WebHookStatus: false,
				WebHookFilter: &WebHookFilter{FilterStatus: false},
				WebHookUrl:    "http://192.168.50.85:1234",
				WebHookSecret: "abcde",
			},
			&WebHook{
				WebHookStatus: false,
				WebHookFilter: &WebHookFilter{FilterStatus: false},
				WebHookUrl:    "http://192.168.50.85:1234",
				WebHookSecret: "abcde",
			},
		},
		WebApi: &WebApi{
			WebApiStatus: false,
			WebApiHost:   "",
			WebApiPort:   "9001",
			WebApiToken:  "abcde",
		},
	}
	if os.Getenv("Phone") != "" {
		conf.Phone = os.Getenv("Phone")
	}
	if os.Getenv("Password") != "" {
		conf.Password = os.Getenv("Password")
	}
	if os.Getenv("AppID") != "" {
		conf.AppID = os.Getenv("AppID")
	}
	if os.Getenv("AppHash") != "" {
		conf.AppHash = os.Getenv("AppHash")
	}
	if os.Getenv("UseMessageDatabase") != "" {
		switch os.Getenv("UseMessageDatabase") {
		case "true":
			conf.UseMessageDatabase = true
		case "false":
			conf.UseMessageDatabase = false
		}
	}
	if os.Getenv("UseFileDatabase") != "" {
		switch os.Getenv("UseFileDatabase") {
		case "true":
			conf.UseFileDatabase = true
		case "false":
			conf.UseFileDatabase = false
		}
	}
	if os.Getenv("UseChatInfoDatabase") != "" {
		switch os.Getenv("UseChatInfoDatabase") {
		case "true":
			conf.UseChatInfoDatabase = true
		case "false":
			conf.UseChatInfoDatabase = false
		}
	}
	if os.Getenv("UseTestDataCenter") != "" {
		switch os.Getenv("UseTestDataCenter") {
		case "true":
			conf.UseTestDataCenter = true
		case "false":
			conf.UseTestDataCenter = false
		}
	}
	if os.Getenv("DatabaseDirectory") != "" {
		conf.DatabaseDirectory = os.Getenv("DatabaseDirectory")
	}
	if os.Getenv("FileDirectory") != "" {
		conf.FileDirectory = os.Getenv("FileDirectory")
	}
	if os.Getenv("IgnoreFileNames") != "" {
		switch os.Getenv("IgnoreFileNames") {
		case "true":
			conf.IgnoreFileNames = true
		case "false":
			conf.IgnoreFileNames = false
		}
	}
	if os.Getenv("ProxyStatus") != "" {
		switch os.Getenv("ProxyStatus") {
		case "true":
			conf.Proxy.ProxyStatus = true
		case "false":
			conf.Proxy.ProxyStatus = false
		}
	}
	if os.Getenv("ProxyType") != "" {
		conf.Proxy.ProxyType = os.Getenv("ProxyType")
	}
	if os.Getenv("ProxyAddr") != "" {
		conf.Proxy.ProxyAddr = os.Getenv("ProxyAddr")
	}
	if os.Getenv("ProxyPort") != "" {
		conf.Proxy.ProxyPort = os.Getenv("ProxyPort")
	}
	if os.Getenv("ProxyUser") != "" {
		conf.Proxy.ProxyUser = os.Getenv("ProxyUser")
	}
	if os.Getenv("ProxyPasswd") != "" {
		conf.Proxy.ProxyPasswd = os.Getenv("ProxyPasswd")
	}

	// 将json数据写入到环境变量
	if os.Getenv("WebHook") != "" {
		var webhook []*WebHook
		err := json.Unmarshal([]byte(os.Getenv("WebHook")), &webhook)
		if err == nil {
			conf.WebHook = webhook
		}else{
			log.Errorf("parase webhook faild, WebHook=%s, err=%v",os.Getenv("WebHook"),err)
		}
	}

	if os.Getenv("WebApiStatus") != "" {
		switch os.Getenv("WebApiStatus") {
		case "true":
			conf.WebApi.WebApiStatus = true
		case "false":
			conf.WebApi.WebApiStatus = false
		}
	}
	if os.Getenv("WebApiHost") != "" {
		conf.WebApi.WebApiHost = os.Getenv("WebApiHost")
	}
	if os.Getenv("WebApiPort") != "" {
		conf.WebApi.WebApiPort = os.Getenv("WebApiPort")
	}
	if os.Getenv("WebApiToken") != "" {
		conf.WebApi.WebApiToken = os.Getenv("WebApiToken")
	}
	return conf
}

type CronJobConfig struct {
	Cron []*CronMessage `json:"cron_config"`
}

type CronMessage struct {
	Cron       string `json:"cron"`
	ToUserName string `json:"to_user_name"`
	TextMsg    string `json:"text_msg"`
}

func LoadCron() *CronJobConfig {
	p := cronPath
	if !utils.PathExists(p) {
		log.Warnf("尝试加载配置文件 %v 失败: 文件不存在", p)
		InitDefaultCronJobConf()
	}
	c := CronJobConfig{}
	err := json.Unmarshal([]byte(utils.ReadAllText(p)), &c)
	if err != nil {
		log.Errorf("尝试加载配置文件 %v 时出现错误: %v", p, err)
		panic("plese check your cron.json")
	}
	return &c
}
func (c CronJobConfig) Save(p string) error {
	data, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}
	utils.WriteAllText(p, string(data))
	return nil
}
func InitDefaultCronJobConf() {
	conf := &CronJobConfig{
		[]*CronMessage{
			{
				Cron:       "* * * * *",
				ToUserName: "@LvanLamCommitCodeBot",
				TextMsg:    "/start",
			},
			{
				Cron:       "* * * * *",
				ToUserName: "@TuringLabbot",
				TextMsg:    "/start",
			},
		},
	}
	conf.Save(cronPath)
}
