package cron

import (
	"github.com/Anonymouscn/dynamic-dns-server/constant"
	"github.com/Anonymouscn/dynamic-dns-server/provider"
	"github.com/Anonymouscn/dynamic-dns-server/script"
	"time"
)

// AutoUpdateDNSCronTask 自动更新 DNS 定时任务
func AutoUpdateDNSCronTask() {
	go func() {
		duration := provider.ScriptConfig.Duration
		if duration < constant.DefaultDuration {
			duration = constant.DefaultDuration
		}
		for {
			script.AutoUpdateDNS()
			t := time.NewTimer(time.Duration(duration) * time.Second)
			<-t.C
		}
	}()
}
