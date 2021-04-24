#!/bin/bash
## 第一次使用的时候需要执行一下，进行账号初始化

docker run -it --rm -name tg_cli \
-v `pwd`:/home \
-e ProxyStatus="true" \
#支持`Socks5`、`HTTP`、`HTTPS`、`MtProto`
-e ProxyType="Socks5" \
-e ProxyAddr="127.0.0.1" \
-e ProxyPort="7890" \
-e ProxyUser="" \
#如果是MtProto 这里填secret \
-e ProxyPasswd=""