# G25WebDAV



这是一个简易的小应用，可实现按照配置文件开启WEBDAV。

```cmd
> go version
go version go1.23.9 windows/amd64
```

## 编译


```cmd
> go build -ldflags="-s -w" -o G25WEBDAV.exe main.go
```

## 初始化

```cmd
> G25WEBDAV.exe
```

首次运行后生成默认配置文件`config.ini`

```
[server]
addr = 127.0.0.1:80
username = pro
password = pro
dir = ./data
```

## 运行

```cmd
> G25WEBDAV.exe
```

