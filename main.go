package main

import (
	"bit.monitor.com/core"
	"bit.monitor.com/global"
	"bit.monitor.com/initialize"
)

func main() {
	// 初始化Viper配置工具
	global.WM_VP = core.Viper()

	// 初始化zap日志
	global.WM_LOG = core.Zap()

	// 初始化gorm连接数据库
	global.WM_DB = initialize.Gorm()
	// 自动迁移数据库表结构
	if global.WM_DB != nil {
		initialize.AutoMigrateMysqlTables(global.WM_DB)
		// 程序结束前，关闭数据库连接
		db, _ := global.WM_DB.DB()
		defer db.Close()
	}

	// 初始化redis服务
	initialize.Redis()

	// 程序启动时，自动启动所有激活状态的定时任务
	go initialize.RunScheduler()

	// 运行http服务
	core.RunWindowsServer()
}
