# docker-compose的使用方式

### 1、创建docker-compose.yaml模板
```yaml
version: "2.2"
services:
  tg_bot:
    image: scjtqs/go-tg:1.8.0
    restart: always
    stdin_open: true
    tty: true
    ports:
      - 9099:9001
    volumes:
      - ./data:/home
    environment:
      ProxyStatus: 'true'
      ProxyType: 'Socks5'
      ProxyAddr: '192.168.50.85'
      ProxyPort: '10808'
      ProxyUser: ''
      ProxyPasswd: ''
      WebHook: '[{"status":false,"http_post_url":"http://192.168.50.85:9991","secret":"abcd","filter":{"status":false}}]'
      WebApiStatus: 'true'
      WebApiPort: '9001'
      WebApiToken: 'abcdefg'
    mem_limit: 250M
    cpus: 1
```

### 2、执行初始化（账号登录）
> docker-compose run --rm tg_bot
>
> 输入 手机号 +86186xxxxxxxx
> 
> 输入验证code ，一般会推送到已登录的tg客户端上
> 
> 登录成功. ctrl+c 取消
> 

### 3、后台运行
> docker-compose up -d

