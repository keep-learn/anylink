package base

// Start 初始化 命令行、配置信息、日志信息
func Start() {

	// 注册cobra 命令
	execute()

	// 将配置文件的值映射到 Cfg 结构体上；便于后期的使用
	initCfg()

	// 初始化项目的日志信息
	initLog()
}

func Test() {
	initLog()
}
