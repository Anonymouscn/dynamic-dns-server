package bootstrap

import (
	"fmt"
	"github.com/Anonymouscn/dynamic-dns-server/constant"
	"github.com/Anonymouscn/dynamic-dns-server/provider"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

// Init 脚本初始化
func Init() error {
	InitScriptConfig()
	return nil
}

// InitScriptConfig 初始化脚本配置
func InitScriptConfig() {
	if _, err := os.Stat(constant.PathToConfig); err != nil {
		if os.IsNotExist(err) {
			panic("script config file not found")
		} else {
			panic(fmt.Sprintf("unknown error: %v", err))
		}
	}
	config, err := os.ReadFile(constant.PathToConfig)
	if err != nil {
		panic(fmt.Sprintf("read config file error: %v", err))
	}
	if err = yaml.Unmarshal(config, &provider.ScriptConfig); err != nil {
		panic(fmt.Sprintf("map config error: %v", err))
	}
	conf := provider.ScriptConfig
	// 加载 GetMyIP API 私密配置
	LoadGetMyIPApi()
	// 根据 DNS 注册商使用不同脚本 (方便后续扩展)
	switch conf.Type {
	case "cloudflare":
		log.Println("cloudflare dynamic DNS is active")
		InitCloudflareSecret()
	}
}

// InitCloudflareSecret 初始化 Cloudflare secret 配置
func InitCloudflareSecret() {
	conf := provider.ScriptConfig
	path := conf.Cloudflare.Secret
	if path != "" {
		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				panic("cloudflare secret file not found")
			} else {
				panic(fmt.Sprintf("unknown error: %v", err))
			}
		}
		secret, err := os.ReadFile(path)
		if err != nil {
			panic(fmt.Sprintf("read secret file error: %v", err))
		}
		if err = yaml.Unmarshal(secret, &provider.CloudflareSecret); err != nil {
			panic(fmt.Sprintf("map cloudflare secret error: %v", err))
		}
	}
}

// LoadGetMyIPApi 加载 get my ip API 接口
func LoadGetMyIPApi() {
	conf := provider.ScriptConfig
	path := conf.GetMyIpApi
	if path == "" {
		panic("path to 'get my ip' API secret file can not be empty")
	}
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			panic("path to 'get my ip' API secret file not found")
		} else {
			panic(fmt.Sprintf("unknown error: %v", err))
		}
	}
	secret, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("read secret file error: %v", err))
	}
	if err = yaml.Unmarshal(secret, &provider.GetMyIPApiSecret); err != nil {
		panic(fmt.Sprintf("map cloudflare secret error: %v", err))
	}
}
