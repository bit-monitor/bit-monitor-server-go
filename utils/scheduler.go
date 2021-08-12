package utils

import (
	"bit.monitor.com/global"
	"errors"
	"fmt"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"sync"
)

// 由于map在并发的情况下是读写不安全的，因此需要改用sync.map的方式来解决
var scheduleTasks = sync.Map{}

// Scheduler
// 此处的BeanName、MethodName，是为了保持与Java后台写法一致所以保留的
// 在Java后台里由于使用了Spring框架的Bean进行实例化注入，但在Go里没必要这么做
type Scheduler struct {
	BeanName       string
	MethodName     string
	Params         string
	SchedulerId    uint64
	CronExpression string
}

// Start
// 开始定时任务
func (s Scheduler) Start(f func(s string) error) error {
	var err error

	c := cron.New()
	// 这里AddFunc不使用CronExpression，是因为Java后台里的实现是用了Spring Framework里的scheduling，
	// 在Java的方案里cron表达式是6位，包括秒的，但在这里使用的cron库里的cron表达式是5位，不支持秒，
	// 因此这里不按照Java后台的方式来，改用cron库提供的方法
	_, err = c.AddFunc("@every 10s", func() {
		err = f(s.Params)
	})
	if err != nil {
		return err
	}

	// 开始定时任务后，存入任务队列中
	c.Start()
	scheduleTasks.Store(s.SchedulerId, c)
	global.WM_LOG.Info("[信息]定时任务已存入任务队列", zap.Any("info", fmt.Sprintf("bean：%v，方法：%v，参数：%v", s.BeanName, s.MethodName, s.Params)))

	return err
}

// Stop
// 停止定时任务
func (s Scheduler) Stop() error {
	var err error
	v, ok := scheduleTasks.Load(s.SchedulerId)
	if !ok {
		err = errors.New("无法停止定时任务，因定时任务队列中找不到已启动的任务")
	} else {
		c, _ := v.(*cron.Cron)
		err = c.Stop().Err()
		scheduleTasks.Delete(s.SchedulerId)
	}
	return err
}

// StopAndDeleteBySchedulerId 根据定时任务id结束并删除任务
func StopAndDeleteBySchedulerId(schedulerId uint64) error {
	var err error
	v, ok := scheduleTasks.Load(schedulerId)
	if !ok {
		err = errors.New("无法停止定时任务，因定时任务队列中找不到已启动的任务")
	} else {
		c, _ := v.(*cron.Cron)
		err = c.Stop().Err()
		scheduleTasks.Delete(schedulerId)
	}
	return err
}
