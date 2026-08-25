# How to build this project

## Prerequisites
- Go 1.23 或更高版本
- TDLib 库及开发头文件（对应 commit: `49b3bcbb6bfebf2ed44dd9f25102d2e1a94a58c4`）

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

若 TDLib 安装在自定义路径，可以使用环境变量指定头文件和库路径：
```shell
CGO_CFLAGS=-I/path/to/tdlib/include CGO_LDFLAGS="-Wl,-rpath,/path/to/tdlib/lib -L/path/to/tdlib/lib -ltdjson" go build -o go-tg
```

## Other os 
> 老老实实自己编译tdlib的依赖包，再go build本项目吧。
> 
> [编译方法](https://tdlib.github.io/td/build.html?language=Go)