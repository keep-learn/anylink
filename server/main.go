// AnyLink 是一个企业级远程办公vpn软件，可以支持多人同时在线使用。

//go:build !windows
// +build !windows

package main

import (
	"embed"
	"os"
	"os/signal"
	"syscall"

	"github.com/bjdgyc/anylink/admin"
	"github.com/bjdgyc/anylink/base"
	"github.com/bjdgyc/anylink/handler"
)

//go:embed ui
var uiData embed.FS

// 程序版本
var CommitId string

func main() {

	// 设置程序版本号
	base.CommitId = CommitId

	// 设置前端页面的路径信息
	admin.UiData = uiData

	// Start 初始化 命令行、配置信息、日志信息
	base.Start()

	// http 相关的操作（包括后端的路由、vpn相关的操作）
	handler.Start()

	// 监听：关闭、重启 相关的信号
	signalWatch()
}

// 监听：关闭、重启 相关的信号
func signalWatch() {
	base.Info("Server pid: ", os.Getpid())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGALRM)
	for {
		sig := <-sigs
		base.Info("Get signal:", sig)
		switch sig {
		case syscall.SIGUSR2:
			// reload
			base.Info("Reload")
		default:
			// stop
			base.Info("Stop")
			handler.Stop()
			return
		}
	}
}
