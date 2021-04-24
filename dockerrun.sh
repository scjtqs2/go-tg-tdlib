#!/bin/bash
## 第一次使用的时候需要执行一下，进行账号初始化

docker run -it --rm --name tg_cli \
-v `pwd`:/home \
-e ProxyStatus="true" \
-e ProxyType="Socks5" \
-e ProxyAddr="127.0.0.1" \
-e ProxyPort="7890" \
-e ProxyUser="" \
-e ProxyPasswd="" \
scjtqs/tg:test

# ProxyType 支持`Socks5`、`HTTP`、`HTTPS`、`MtProto`

# ProxyPasswd 如果是MtProto 这里填secret \

docker run -it --rm --name tg_cli \
-v `pwd`:/home \
-e ProxyStatus="true" \
-e ProxyType="Socks5" \
-e ProxyAddr="pi.scjtqs.com" \
-e ProxyPort="10808" \
-e ProxyUser="" \
-e ProxyPasswd="" \
scjtqs/tg:test