#使用说明

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

