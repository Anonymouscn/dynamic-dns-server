package config

// ScriptConfig 脚本配置
type ScriptConfig struct {
	Name       string               `yaml:"name"`
	Type       string               `yaml:"type"`
	Duration   int64                `yaml:"duration"`
	GetMyIpApi string               `yaml:"get_my_ip_api"`
	Cloudflare CloudflareSecretPath `yaml:"cloudflare"`
}

// CloudflareSecretPath cloudflare 私密文件路径
type CloudflareSecretPath struct {
	Secret string `yaml:"secret"`
	Proxy  bool   `yaml:"proxy"`
	TTL    int64  `yaml:"ttl"`
	Type   string `yaml:"type"`
}

// CloudflareSecret cloudflare 私密配置
type CloudflareSecret struct {
	Email         string `yaml:"email"`
	AccountID     string `yaml:"account_id"`
	TargetDomain  string `yaml:"target_domain"`
	ZoneDomain    string `yaml:"zone_domain"`
	Authorization string `yaml:"authorization"`
	ZoneID        string `yaml:"zone_id"`
}

// GetMyIPApiSecret GetMyIP API 私密配置
type GetMyIPApiSecret struct {
	API string `yaml:"api"`
}
