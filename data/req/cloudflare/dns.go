package cloudflare

// DNSRecordSearchReq DNS 记录搜索
type DNSRecordSearchReq struct {
	Name string `json:"name"`
}

// DNSSettings DNS 设置附带信息
type DNSSettings struct {
}

// DNSRecordUpdateReq DNS 记录更新
type DNSRecordUpdateReq struct {
	Comment  string      `json:"comment"`
	Name     string      `json:"name"`
	Proxied  bool        `json:"proxied"`
	Settings DNSSettings `json:"settings"`
	Tags     []any       `json:"tags"`
	TTL      int64       `json:"ttl"`
	Content  string      `json:"content"`
	Type     string      `json:"type"`
}
