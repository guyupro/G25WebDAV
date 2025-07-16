package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/webdav"
	"gopkg.in/ini.v1"
	"log"
	"net/http"
	"os"
)

const configFile = "config.ini"

type Config struct {
	Addr     string
	Username string
	Password string
	Dir      string
}

func defaultConfig() *Config {
	return &Config{
		Addr:     "127.0.0.1:80",
		Username: "pro",
		Password: "pro",
		Dir:      "./data",
	}
}

func createDefaultConfigFile() error {
	cfg := defaultConfig()
	os.MkdirAll(cfg.Dir, 0755)
	f, err := os.Create(configFile)
	if err != nil {
		return err
	}
	defer f.Close()
	fmt.Fprintf(f, "[server]\naddr = %s\nusername = %s\npassword = %s\ndir = %s\n", cfg.Addr, cfg.Username, cfg.Password, cfg.Dir)
	return nil
}

func loadConfig() (*Config, error) {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if err := createDefaultConfigFile(); err != nil {
			return nil, err
		}
		fmt.Println("首次启动，已生成默认配置文件 config.ini，请修改后启动。 ")
		os.Exit(0)
	}
	iniCfg, err := ini.Load(configFile)
	if err != nil {
		return nil, err
	}
	sec := iniCfg.Section("server")
	return &Config{
		Addr:     sec.Key("addr").MustString("127.0.0.1:80"),
		Username: sec.Key("username").MustString("pro"),
		Password: sec.Key("password").MustString("pro"),
		Dir:      sec.Key("dir").MustString("./data"),
	}, nil
}

func basicAuth(username, password string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok || u != username || p != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="webdav"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized\n"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	verbose := flag.Bool("v", false, "详细模式，输出每个HTTP请求日志")
	flag.Parse()

	cfg, err := loadConfig()
	if err != nil {
		log.Fatal("配置文件读取失败:", err)
	}
	os.MkdirAll(cfg.Dir, 0755)

	h := &webdav.Handler{
		Prefix:     "/",
		FileSystem: webdav.Dir(cfg.Dir),
		LockSystem: webdav.NewMemLS(),
	}

	handler := basicAuth(cfg.Username, cfg.Password, h)
	if *verbose {
		handler = logRequest(handler)
	}

	fmt.Printf("WebDAV服务启动: http://%s 目录: %s 用户: %s\n", cfg.Addr, cfg.Dir, cfg.Username)
	if err := http.ListenAndServe(cfg.Addr, handler); err != nil {
		log.Fatal(err)
	}
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s from %s", r.Method, r.URL.Path, r.Proto, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
