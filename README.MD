#使用说明
>
> 
##首次使用需要进行登录
> 需要使用docker run进行登录处理
```shell
docker run -it --rm --name tg_cli \
-v `pwd`:/home \
-e ProxyStatus="true" \
-e ProxyType="Socks5" \
-e ProxyAddr="127.0.0.1" \
-e ProxyPort="7890" \
-e ProxyUser="" \
-e ProxyPasswd="" \
scjtqs/tg:test
```
>  ProxyType 支持`Socks5`、`HTTP`、`HTTPS`、`MtProto`
> 
>  ProxyPasswd 如果是MtProto 这里填secret \

##定时任务配置
> 首次运行后，当前目录会生成一个cron.json目录
> 
```json
{
	"cron_config": [
		{
			"cron": "* * * * *",
			"to_user_name": "@LvanLamCommitCodeBot",
			"text_msg": "/start"
		},
		{
			"cron": "* * * * *",
			"to_user_name": "@TuringLabbot",
			"text_msg": "/start"
		}
	]
}
```
> cron 是定时配置，参考linux的crontab
> 
> to_user_name 是对方的username
> 
> text_msg 是要发送的消息类容
> 
> 请自行修改