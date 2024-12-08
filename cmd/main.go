package main

import (
	"context"
	"fmt"
	"github.com/Anonymouscn/dynamic-dns-server/bootstrap"
	"github.com/Anonymouscn/dynamic-dns-server/cron"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	// 初始化脚本配置
	if err := bootstrap.Init(); err != nil {
		panic(fmt.Sprintf("init config fail: %v", err))
	}
	// 定时执行 DNS 更新任务
	cron.AutoUpdateDNSCronTask()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-c
		cancel()
	}()
	<-ctx.Done()
}
