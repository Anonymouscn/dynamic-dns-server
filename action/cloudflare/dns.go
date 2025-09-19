package cloudflare

import (
	"net/http"

	"github.com/Anonymouscn/dynamic-dns-server/data/req/cloudflare"
	"github.com/Anonymouscn/dynamic-dns-server/data/resp"
	"github.com/Anonymouscn/go-partner/restful"
)

// GetDNSRecordList 获取 DNS 信息
// more: https://developers.cloudflare.com/api/operations/dns-records-for-a-zone-list-dns-records
func GetDNSRecordList(email, zoneID, authKey string,
	search *cloudflare.DNSRecordSearchReq) (*resp.CloudflareResp[resp.DNSRecord], error) {
	rc := restful.NewRestClient().
		SetURL("https://api.cloudflare.com/client/v4/zones").
		SetPath(restful.Path{zoneID, "dns_records"}).
		SetHeaders(restful.Data{"X-Auth-Email": email, "Authorization": "Bearer " + authKey})
	// 条件构造
	if search != nil {
		if search.Name != "" {
			rc.SetQuery(restful.Data{"name": search.Name})
		}
	}
	rc.Get()
	result := &resp.CloudflareResp[resp.DNSRecord]{}
	if err := rc.Map(result); err != nil {
		return nil, err
	}
	return result, nil
}

// UpdateDNSRecord 更新 DNS 记录
// more: https://developers.cloudflare.com/api/operations/dns-records-for-a-zone-patch-dns-record
func UpdateDNSRecord(email, zoneID, recordID, authKey string,
	updateInfo *cloudflare.DNSRecordUpdateReq) (*resp.CloudflareResp[resp.DNSRecordUpdateInfo], error) {
	rc := restful.NewRestClient().
		SetURL("https://api.cloudflare.com/client/v4/zones").
		SetPath(restful.Path{zoneID, "dns_records", recordID}).
		SetHeaders(restful.Data{"X-Auth-Email": email, "Authorization": "Bearer " + authKey}).
		SetBody(updateInfo).
		Do(http.MethodPatch)
	result := &resp.CloudflareResp[resp.DNSRecordUpdateInfo]{}
	if err := rc.Map(result); err != nil {
		return nil, err
	}
	return result, nil
}
