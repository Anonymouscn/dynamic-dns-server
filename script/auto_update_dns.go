package script

import (
	"fmt"
	"github.com/Anonymouscn/dynamic-dns-server/action/checkIP"
	"github.com/Anonymouscn/dynamic-dns-server/action/cloudflare"
	"github.com/Anonymouscn/dynamic-dns-server/constant"
	cloudflare2 "github.com/Anonymouscn/dynamic-dns-server/data/req/cloudflare"
	"github.com/Anonymouscn/dynamic-dns-server/provider"
	"log"
)

// AutoUpdateDNS 自动更新 DNS 脚本
func AutoUpdateDNS() {
	conf := provider.ScriptConfig
	proxy := conf.Cloudflare.Proxy
	currentIP, err := checkIP.GetMyIP()
	if err != nil {
		panic("get current ip error")
	}
	log.Println(fmt.Sprintf("current IP: %v[%v]", currentIP, ProxyMsg(proxy)))
	secret := provider.CloudflareSecret
	res, err := cloudflare.GetDNSRecordList(secret.Email, secret.ZoneID, secret.Authorization,
		&cloudflare2.DNSRecordSearchReq{
			Name: secret.TargetDomain,
		},
	)
	if err != nil {
		log.Println(fmt.Sprintf("get dns record list fail: %v", err))
		return
	}
	match := res.Result
	if len(match) == 0 {
		log.Println("no target domain name matches, please check whether the target domain name is incorrect")
	} else {
		target := match[0]
		log.Println(fmt.Sprintf("current dns: %v[%v]", target.Content, ProxyMsg(target.Proxied)))
		if currentIP == target.Content && proxy == target.Proxied {
			log.Println("current IP config has not changed and does not need to be updated.")
		} else {
			domainType := conf.Cloudflare.Type
			if domainType == "" {
				domainType = "A"
			}
			ttl := conf.Cloudflare.TTL
			if ttl < constant.DefaultTTL {
				ttl = constant.DefaultTTL
			}
			res, err := cloudflare.UpdateDNSRecord(secret.Email, secret.ZoneID, target.ID, secret.Authorization,
				&cloudflare2.DNSRecordUpdateReq{
					Comment: fmt.Sprintf("Dynamic DNS client auto update [ip: %v]", currentIP),
					Name:    secret.TargetDomain,
					Proxied: conf.Cloudflare.Proxy,
					TTL:     ttl,
					Content: currentIP,
					Type:    domainType,
				})
			if err != nil {
				log.Println(fmt.Sprintf("DNS record update fail: %v", err))
			}
			if !res.Success {
				log.Println(fmt.Sprintf("DNS record update fail: [messages: %v, errors: %v]", res.Messages, res.Errors))
			}
			log.Println(fmt.Sprintf("DNS record update success: %v[%v] -> %v[%v]", target.Content, ProxyMsg(target.Proxied), currentIP, ProxyMsg(proxy)))
		}
	}
}

// ProxyMsg CDN 代理状态解释信息
func ProxyMsg(proxy bool) string {
	if proxy {
		return "proxied"
	}
	return "only dns"
}
