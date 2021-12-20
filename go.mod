module github.com/scjtqs/go-tg

go 1.16

require (
	github.com/Arman92/go-tdlib/v2 v2.0.1-0.20210605080123-2454be49c341
	github.com/gin-gonic/gin v1.7.2-0.20211215152723-fb5f04541787
	github.com/guonaihong/gout v0.1.9
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.4 // indirect
	github.com/robfig/cron/v3 v3.0.1
	github.com/sirupsen/logrus v1.8.1
	github.com/t-tomalak/logrus-easy-formatter v0.0.0-20190827215021-c074f06c5816
	github.com/tidwall/gjson v1.7.5
)

//replace (
//	github.com/Arman92/go-tdlib master => /Users/apple/Workspace/git/go-tdlib
//)
//replace github.com/Arman92/go-tdlib v1.0.1-0.20210605080123-2454be49c341 => github.com/Arman92/go-tdlib/v2 v2.0.0-20211210144712-d8b8869d8e49
replace github.com/Arman92/go-tdlib/v2 v2.0.1-0.20210605080123-2454be49c341 => github.com/ffenix113/go-tdlib/v2 v2.0.0-20211204191913-dbb38e1deb80
