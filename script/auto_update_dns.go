package script

import (
	"fmt"
	"log"

	"github.com/Anonymouscn/dynamic-dns-server/action/checkIP"
	"github.com/Anonymouscn/dynamic-dns-server/action/cloudflare"
	"github.com/Anonymouscn/dynamic-dns-server/constant"
	cloudflare2 "github.com/Anonymouscn/dynamic-dns-server/data/req/cloudflare"
	"github.com/Anonymouscn/dynamic-dns-server/provider"
	alidns "github.com/alibabacloud-go/alidns-20150109/v2/client"
	aliapi "github.com/alibabacloud-go/darabonba-openapi/client"
	teautil "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
)

// Cloudflare 域名代理状态常量
const (
	CfProxyState         = "Proxy"
	CfDnsOnlyState       = "DNS Only"
	ProxyNotSupportState = "Proxy Not Support"
)

// 代理商常量定义
const (
	Cloudflare = "cloudflare"
	Aliyun     = "aliyun"
	Tencent    = "tencent"
	Bytedance  = "bytedance"
)

// 全局变量定义
var (
	handlers = HandlerChain{
		&CloudflareHandler{},
		&AliyunHandler{},
		&TencentHandler{},
		&BytedanceHandler{},
	}
	worker Handler
)

// InitHandlers 初始化守护进程处理器
func InitHandlers() {
	for _, handler := range handlers {
		if handler.Name() == provider.ScriptConfig.Type {
			worker = handler.Instance()
			break
		}
	}
	if worker == nil {
		panic("no script handler found")
	}
}

// HandlerChain 守护进程处理链
type HandlerChain []Handler

// Handler 守护进程处理器
type Handler interface {
	Name() string      // 获取处理器名称
	Instance() Handler // 获取处理器实例
	AutoUpdateDNS()    // 自动更新 DNS (启动 DNS 更新守护线程程)
}

// CloudflareHandler Cloudflare 守护进程处理器
type CloudflareHandler struct{}

func (handler *CloudflareHandler) Name() string {
	return Cloudflare
}

func (handler *CloudflareHandler) Instance() Handler {
	return handler
}

// CFProxyMsg CF 代理状态解释信息
func (handler *CloudflareHandler) CFProxyMsg(proxy bool) string {
	if proxy {
		return CfProxyState
	}
	return CfDnsOnlyState
}

func (handler *CloudflareHandler) AutoUpdateDNS() {
	conf := provider.ScriptConfig
	proxy := conf.Cloudflare.Proxy
	currentIP, err := checkIP.GetMyIP()
	if err != nil {
		panic("get current ip error")
	}
	log.Println(fmt.Sprintf("current IP: %v[%v]", currentIP, handler.CFProxyMsg(proxy)))
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
		log.Println(fmt.Sprintf("current dns: %v[%v]", target.Content, handler.CFProxyMsg(target.Proxied)))
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
			log.Println(fmt.Sprintf("DNS record update success: %v[%v] -> %v[%v]",
				target.Content, handler.CFProxyMsg(target.Proxied), currentIP, handler.CFProxyMsg(proxy)),
			)
		}
	}
}

// AliyunHandler Aliyun 守护进程处理器
type AliyunHandler struct {
	*alidns.Client // 阿里云 SDK 注入
}

func (handler *AliyunHandler) Name() string {
	return Aliyun
}

func (handler *AliyunHandler) Instance() Handler {
	return &AliyunHandler{
		handler.createAliyunDnsClient(),
	}
}

// createAliyunDnsClient 创建 aliyun dns 客户端
func (handler *AliyunHandler) createAliyunDnsClient() *alidns.Client {
	conf := provider.AliyunSecret
	client, err := alidns.NewClient(&aliapi.Config{
		AccessKeyId:     &conf.AccessKeyID,
		AccessKeySecret: &conf.AccessKeySecret,
		RegionId:        &conf.RegionID,
	})
	if err != nil {
		panic("create aliyun client error")
	}
	return client
}

// getTargetDNSRecord 获取目标 DNS 记录
func (handler *AliyunHandler) getTargetDNSRecord() (string, string, error) {
	conf := provider.ScriptConfig
	aliConf := provider.AliyunSecret
	resp, err := handler.DescribeDomainRecords(&alidns.DescribeDomainRecordsRequest{
		DomainName: &aliConf.TargetDomain,
		Type:       &conf.Aliyun.Type,
	})
	if err != nil {
		return "", "", err
	}
	if tea.BoolValue(teautil.IsUnset(tea.ToMap(resp))) ||
		tea.BoolValue(teautil.IsUnset(tea.ToMap(resp.Body.DomainRecords.Record[0]))) {
		return "", "", fmt.Errorf("get dns record fail")
	}
	record := resp.Body.DomainRecords.Record[0]
	return *record.Value, *record.RecordId, err
}

// modifyAliyunDnsRecord 修改 DNS 记录
func (handler *AliyunHandler) modifyAliyunDnsRecord(recordID, value string) error {
	_, err := handler.UpdateDomainRecord(&alidns.UpdateDomainRecordRequest{
		RecordId: &recordID,
		Value:    &value,
		RR:       &provider.AliyunSecret.RR,
		Type:     &provider.ScriptConfig.Aliyun.Type,
	})
	if err != nil {
		return err
	}
	return nil
}

// AutoUpdateDNS 自动更新 DNS (DNS 更新守护进程)
func (handler *AliyunHandler) AutoUpdateDNS() {
	// 获取当前 ip
	currentIP, err := checkIP.GetMyIP()
	if err != nil {
		panic("get current ip error")
	}
	log.Println(fmt.Sprintf("current IP: %v[%v]", currentIP, ProxyNotSupportState))
	// 获取当前 DNS 生效的 IP
	record, id, err := handler.getTargetDNSRecord()
	if err != nil {
		log.Println(fmt.Sprintf("get dns record fail: %v", err))
		return
	}
	log.Println(fmt.Sprintf("current dns: %v, record id: %v", record, id))
	// DNS 记录与当前 IP 比对处理
	if record == currentIP {
		log.Println("current IP config has not changed and does not need to be updated.")
	} else {
		if err := handler.modifyAliyunDnsRecord(id, currentIP); err != nil {
			log.Println(fmt.Sprintf("update dns record fail: %v", err))
			return
		}
		log.Println(fmt.Sprintf("DNS record update success: %v[%v] -> %v[%v]",
			record, ProxyNotSupportState, currentIP, ProxyNotSupportState),
		)
	}
}

// TencentHandler 腾讯云守护进程处理器 todo 后续接入
type TencentHandler struct{}

func (handler *TencentHandler) Name() string {
	return Tencent
}

func (handler *TencentHandler) Instance() Handler {
	return handler
}

func (handler *TencentHandler) AutoUpdateDNS() {}

// BytedanceHandler 字节火山引擎守护进程处理器 todo 后续接入
type BytedanceHandler struct{}

func (handler *BytedanceHandler) Name() string {
	return Bytedance
}

func (handler *BytedanceHandler) Instance() Handler {
	return handler
}

func (handler *BytedanceHandler) AutoUpdateDNS() {}

// AutoUpdateDNS 自动更新 DNS (DNS 更新守护进程)
func AutoUpdateDNS() {
	worker.AutoUpdateDNS()
}
