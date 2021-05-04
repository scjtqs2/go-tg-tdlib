# 使用说明
> 此项目主要用于tg自己账号登录 定时发送消息给某人
> 
> 当前容器支持 amd64、i386、arm64、armhf

# 特别说明！！！ tg的验证码不是发的手机短信，是发送到了tg客户端上！！！！国内手机号需要加上+86！！

## 首次使用需要进行登录
> 需要使用docker run进行登录处理
```shell
docker run -it --rm --name tg_cli \
-v `pwd`:/home \
-e AppID="1807909" \
-e AppHash="4b1594bcfab16b370686b14d85c60559" \
-e ProxyStatus="true" \
-e ProxyType="Socks5" \
-e ProxyAddr="192.168.0.111" \
-e ProxyPort="7890" \
-e ProxyUser="" \
-e ProxyPasswd="" \
scjtqs/go-tg:latest
```
> AppID 可不填
> 
> AppHash 可不填
> 
> ProxyStatus 是否开启 网络代理，默认为 "false"
> 
> ProxyType 支持`Socks5`、`HTTP`、`HTTPS`、`MtProto`
> 
> ProxyPasswd 如果是MtProto 这里填secret 
> 
> ProxyAddr 这个别填 localhost和 127.0.0.1。否则会识别成docker容器内部，而非宿主机
> 
> 国内阿里云镜像 registry.cn-hangzhou.aliyuncs.com/scjtqs/go-tg:latest
> 
> 其他 环境变量，请查看 Dockerfile.go-td 里面的 ENV部分
> 
> WebApiStatus 是否开启api，默认为 "false"。要开启请填 "true"
> 
> WebApiPort api监听端口，默认"9001"
> 
> WebApiToken api鉴权token，默认 "abcde"

## 定时任务配置
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
> cron 是定时配置，参考linux的crontab，5个参数为分钟即，6个参数为秒级
> 
> to_user_name 是对方的username
> 
> text_msg 是要发送的消息类容
> 
> 请自行修改
> 
> cron薅羊毛参考 `1,2,3,4,5 0 1,8,16,24 * ?`

## docker后台运行
```shell
docker run  -D --restart=always --name tg_bot \
-v /yourdir:/home \
-e ProxyStatus="true" \
-e ProxyType="Socks5" \
-e ProxyAddr="192.168.0.111" \
-e ProxyPort="7890" \
-e ProxyUser="" \
-e ProxyPasswd="" \
scjtqs/go-tg:latest
```

## api 功能
[点击查看](doc/api.md)

## 使用docker-compose 方式运行
[点击查看](doc/docker-compose.md)

## for develop
[如果我想自己修改代码，要如何编译当前代码](doc/dev.md)
