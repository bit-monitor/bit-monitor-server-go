package service

import (
	"bit.monitor.com/global"
	"bit.monitor.com/model"
	"bit.monitor.com/model/validation"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AddSubscriber 保存预警订阅通知记录
func AddSubscriber(tx *gorm.DB, s validation.Subscriber, alarmId uint64) (err error) {
	db := tx.Model(&model.AmsSubscriber{})

	// 保存实体
	entity := model.AmsSubscriber{
		AlarmId:    alarmId,
		Subscriber: s.Subscriber,
		IsActive:   *s.IsActive,
		Category:   *s.Category,
	}
	err = db.Create(&entity).Error
	if err != nil {
		global.WM_LOG.Error("保存预警订阅通知记录失败", zap.Any("err", err))
	}
	return err
}

// DeleteAllByAlarmId 根据alarmId删除所有关联的subscriber
func DeleteAllByAlarmId(tx *gorm.DB, alarmId uint64) (err error) {
	db := tx.Model(&model.AmsSubscriber{})
	err = db.Where("`alarm_id` = ?", alarmId).Delete(model.AmsSubscriber{}).Error
	return err
}

// GetAllByAlarmId 根据alarmId获取所有关联的subscriber
func GetAllByAlarmId(alarmId uint64) []*model.AmsSubscriber {
	db := global.WM_DB.Model(&model.AmsSubscriber{})
	var records []*model.AmsSubscriber
	db.Where("`alarm_id` = ?", alarmId).Find(&records)
	if records == nil {
		global.WM_LOG.Error("根据alarmId获取所有关联的subscriber失败", zap.Any("err", "找不到记录"))
	}
	return records
}

// GetCategoryBySubscriberId 根据id获取category
func GetCategoryBySubscriberId(subscriberId uint64) (err error, category int8) {
	db := global.WM_DB.Model(&model.AmsSubscriber{})
	var e model.AmsSubscriber
	err = db.Where("`id` = ?", subscriberId).First(&e).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("subscriber不存在")
		}
		global.WM_LOG.Error("subscriber不存在", zap.Any("err", err))
		return
	}
	category = e.Category
	return
}
