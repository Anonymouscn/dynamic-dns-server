package resp

import "github.com/Anonymouscn/dynamic-dns-server/data/req/cloudflare"

type CloudflareResp[T any] struct {
	Result     []T              `json:"result"`
	Errors     []CloudflareInfo `json:"errors"`
	Messages   []CloudflareInfo `json:"messages"`
	Success    bool             `json:"success"`
	ResultInfo ResultInfo       `json:"result_info"`
}

type CloudflareInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type DNSRecord struct {
	ID         string      `json:"id"`
	ZoneID     string      `json:"zone_id"`
	ZoneName   string      `json:"zone_name"`
	Name       string      `json:"name"`
	Type       string      `json:"type"`
	Content    string      `json:"content"`
	Proxiable  bool        `json:"proxiable"`
	Proxied    bool        `json:"proxied"`
	TTL        int         `json:"ttl"`
	Settings   DNSSettings `json:"settings"`
	Meta       DNSMeta     `json:"meta"`
	Comment    string      `json:"comment"`
	Tags       []any       `json:"tags"`
	CreatedOn  string      `json:"createdOn"`
	ModifiedOn string      `json:"modified_on"`
}

type DNSMeta struct {
	AutoAdded           bool `json:"auto_added"`
	ManagedByApps       bool `json:"managed_by_apps"`
	ManagedByArgoTunnel bool `json:"managed_by_argo_tunnel"`
}

type DNSSettings struct {
}

type ResultInfo struct {
	Page       int64 `json:"page"`
	PerPage    int64 `json:"per_page"`
	Count      int64 `json:"count"`
	TotalCount int64 `json:"total_count"`
	TotalPage  int64 `json:"total_page"`
}

type DNSRecordUpdateInfo struct {
	cloudflare.DNSRecordUpdateReq
	CommentModifiedOn string   `json:"comment_modified_on"`
	CreatedOn         string   `json:"created_on"`
	ID                string   `json:"id"`
	Meta              *DNSMeta `json:"meta"`
	ModifiedOn        string   `json:"modified_on"`
	Proxiable         bool     `json:"proxiable"`
	TagsModifiedOn    string   `json:"tags_modified_on"`
}
