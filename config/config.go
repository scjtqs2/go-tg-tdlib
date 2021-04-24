package config

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"os"
)

var cronPath="cron.json"

type JsonConfig struct {
	Phone               string `json:"phone_number"`
	Password            string `json:"password"`
	AppID               string `json:"appid"`
	AppHash             string `json:"app_hash"`
	UseMessageDatabase  bool   `json:"use_message_database"`
	UseFileDatabase     bool   `json:"use_file_database"`
	UseChatInfoDatabase bool   `json:"use_chat_info_database"`
	UseTestDataCenter   bool   `json:"use_test_data_center"`
	DatabaseDirectory   string `json:"database_directory"`
	FileDirectory       string `json:"file_directory"`
	IgnoreFileNames     bool   `json:"ignore_file_name"`
	Proxy               Proxy  `json:"proxy"`
}

type Proxy struct {
	ProxyStatus bool   `json:"status"` //开关
	ProxyType   string `json:"type"`   // "HTTP"，HTTPS，Socks5、MtProto
	ProxyAddr   string `json:"addr"`   // 10.0.0.1
	ProxyPort   string `json:"port"`   // 1234
	ProxyUser   string `json:"user"`   // user
	ProxyPasswd string `json:"passwd"` //passwd
}

func Load(p string) *JsonConfig {
	if !PathExists(p) {
		log.Warnf("尝试加载配置文件 %v 失败: 文件不存在", p)
		return nil
	}
	c := JsonConfig{}
	err := json.Unmarshal([]byte(ReadAllText(p)), &c)
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
	WriteAllText(p, string(data))
	return nil
}

func DefaultConfig() *JsonConfig {
	conf := &JsonConfig{
		Phone:               "+8618611500511",
		Password:            "12345",
		AppID:               "187786",
		AppHash:             "e782045df67ba48e441ccb105da8fc85",
		UseMessageDatabase:  true,
		UseFileDatabase:     true,
		UseChatInfoDatabase: true,
		UseTestDataCenter:   false,
		DatabaseDirectory:   "./tdlib-db",
		FileDirectory:       "./tdlib-files",
		IgnoreFileNames:     false,
		Proxy: Proxy{
			ProxyStatus: true,
			ProxyType:   "Socks5",
			ProxyAddr:   "127.0.0.1",
			ProxyPort:   "1234",
			ProxyUser:   "",
			ProxyPasswd: "",//mtp 的secret也是这个字段
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
	return conf
}

type CronJobConfig struct {
	Cron []CronMessage `json:"cron_config"`
}

type CronMessage struct {
	Cron string `json:"cron"`
	ToUserName string `json:"to_user_name"`
	TextMsg string `json:"text_msg"`
}

func LoadCron() *CronJobConfig {
	p:=cronPath
	if !PathExists(p) {
		log.Warnf("尝试加载配置文件 %v 失败: 文件不存在", p)
		InitDefaultCronJobConf()
	}
	c := CronJobConfig{}
	err := json.Unmarshal([]byte(ReadAllText(p)), &c)
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
	WriteAllText(p, string(data))
	return nil
}
func InitDefaultCronJobConf()  {
	conf:=&CronJobConfig{
		[]CronMessage{
			{
				Cron: "* * * * *",
				ToUserName: "https://t.me/LvanLamCommitCodeBot",
				TextMsg: "/start",
			},
			{
				Cron: "* * * * *",
				ToUserName: "https://t.me/TuringLabbot",
				TextMsg: "/start",
			},
		},
	}
	conf.Save(cronPath)
}