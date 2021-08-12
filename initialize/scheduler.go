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
			// 在循环遍历的情况下，如果在协程里直接引用遍历的变量会出现问题，
			// 因为主流程跑得比协程快，所以等到协程拿到遍历的变量基本已经是后面的数据了，
			// 所以可以借助闭包的原理，在循环体内再定义一个局部变量
			// 参考资料：https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
			scheduler := scheduler
			go func() {
				err := service.StartScheduler(scheduler)
				if err != nil {
					global.WM_LOG.Error("[失败]程序启动时，自动启动所有激活状态的定时任务", zap.Any("RunScheduler", err))
				}
			}()
		}
	}
	global.WM_LOG.Info("[成功]程序启动时，自动启动所有激活状态的定时任务")
}
