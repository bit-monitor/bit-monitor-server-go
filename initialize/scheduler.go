package initialize

import (
	"bit.monitor.com/global"
	"bit.monitor.com/service"
	"go.uber.org/zap"
)

// RunScheduler 程序启动时，自动启动所有激活状态的定时任务
func RunScheduler() {

	// 获取所有运行中的定时任务
	state := int8(1)
	schedulerList := service.GetSchedulerListByState(state)

	// 启动任务
	if len(schedulerList) > 0 {
		for _, scheduler := range schedulerList {
			err := service.StartScheduler(scheduler)
			if err != nil {
				global.WM_LOG.Error("[失败]程序启动时，自动启动所有激活状态的定时任务", zap.Any("RunScheduler", err))
			}
		}
	}
	global.WM_LOG.Info("[成功]程序启动时，自动启动所有激活状态的定时任务")
}
