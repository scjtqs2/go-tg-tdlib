# How to build this project

## How to build tdlib
> before go build ,you need to have tdlib build.
> 
> [how to build tdlib](https://tdlib.github.io/td/build.html?language=Go)

## On macos
```shell
## brew 有一个预编译的tdlib包，不用自己编译，直接下载就好了。
brew install tdlib
git clone https://github.com/scjtqs/go-tg.git go-tg
cd go-tg
go build -o go-tg
./go-tg -h
```
## Other os 
> 老老实实自己编译tdlib的依赖包，再go build本项目吧。
> 
> [编译方法](https://tdlib.github.io/td/build.html?language=Go)