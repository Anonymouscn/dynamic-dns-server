package provider

import "github.com/Anonymouscn/dynamic-dns-server/data/config"

var (
	ScriptConfig     *config.ScriptConfig    // 脚本配置信息
	GetMyIPApiSecret config.GetMyIPApiSecret // GetMyIP API 私密配置
	CloudflareSecret config.CloudflareSecret // Cloudflare 私密配置信息
	AliyunSecret     config.AliyunSecret     // Aliyun 私密配置信息
)
