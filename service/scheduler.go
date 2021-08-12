package service

import (
	"bit.monitor.com/global"
	"bit.monitor.com/model"
	"bit.monitor.com/utils"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

const (
	BeanName       = "AlarmScheduler"
	MethodName     = "startScheduler"
	CronExpression = "*/10 * * * * ?"
)

// AddScheduler 保存定时任务
func AddScheduler(tx *gorm.DB, params string) (error, *model.TmsScheduler) {
	db := tx.Model(&model.TmsScheduler{})
	schedulerEntity := model.TmsScheduler{
		BeanName:       BeanName,
		MethodName:     MethodName,
		Params:         params,
		CronExpression: CronExpression,
		CreateTime:     time.Now(),
		UpdateTime:     time.Now(),
		State:          1,
	}
	err := db.Save(&schedulerEntity).Error
	if err != nil {
		return err, nil
	}
	return nil, &schedulerEntity
}

// DeleteSchedulerById 根据id删除Scheduler实体
func DeleteSchedulerById(id uint64) error {
	err := global.WM_DB.Where("`id` = ?", id).Delete(&model.TmsScheduler{}).Error
	return err
}

// StartScheduler 启动定时任务
func StartScheduler(scheduler *model.TmsScheduler) error {
	s := utils.Scheduler{
		BeanName:       scheduler.BeanName,
		MethodName:     scheduler.MethodName,
		Params:         scheduler.Params,
		SchedulerId:    scheduler.Id,
		CronExpression: scheduler.CronExpression,
	}
	err := s.Start(func(params string) error {
		var errTemp error
		state := 0
		errorMsg := ""
		timeBefore := time.Now()

		global.WM_LOG.Info("定时任务开始执行", zap.Any("info", fmt.Sprintf("bean：%v，方法：%v，参数：%v", s.BeanName, s.MethodName, s.Params)))

		// 要执行的任务
		errTemp = startAlarmSchedule(params)

		if errTemp != nil {
			errorMsg = errTemp.Error()
			global.WM_LOG.Info("定时任务执行异常", zap.Any("info", fmt.Sprintf("bean：%v，方法：%v，参数：%v，异常：%v", s.BeanName, s.MethodName, s.Params, errTemp)))
		} else {
			state = 1
		}
		timeAfter := time.Now()
		timeCost := timeAfter.Sub(timeBefore).Milliseconds()
		global.WM_LOG.Info("定时任务执行结束", zap.Any("info", fmt.Sprintf("bean：%v，方法：%v，参数：%v，耗时：%v毫秒", s.BeanName, s.MethodName, s.Params, timeCost)))

		// 保存定时任务执行记录
		// 这里由于cron库的定时任务是通过goroutine启动的异步函数，因此下面保存无法用tx的方式，要么只能使用gorm的手动控制事务的方式
		// gorm手动事务的文档地址：https://gorm.io/docs/transactions.html
		db := global.WM_DB.Model(&model.TmsSchedulerRecord{})
		record := model.TmsSchedulerRecord{
			SchedulerId: s.SchedulerId,
			State:       int8(state),
			CreateTime:  time.Now(),
			TimeCost:    int8(timeCost),
			ErrorMsg:    errorMsg,
		}
		errTemp = db.Save(&record).Error

		return errTemp
	})
	return err
}

// GetSchedulerListByState 根据state获取Scheduler实体列表
func GetSchedulerListByState(state int8) []*model.TmsScheduler {
	list := make([]*model.TmsScheduler, 0)
	db := global.WM_DB.Model(&model.TmsScheduler{})
	if state < 0 || state > 1 {
		return list
	}
	db.Where("`state` = ?", state).Find(&list)
	return list
}
