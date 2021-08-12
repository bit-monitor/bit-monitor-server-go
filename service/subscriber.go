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
		global.WM_LOG.Error("[失败]保存预警订阅通知记录", zap.Any("err", err))
	}
	return err
}

// DeleteAllSubscriberByAlarmId 根据alarmId删除所有关联的subscriber
func DeleteAllSubscriberByAlarmId(tx *gorm.DB, alarmId uint64) (err error) {
	db := tx.Model(&model.AmsSubscriber{})
	err = db.Where("`alarm_id` = ?", alarmId).Delete(model.AmsSubscriber{}).Error
	return err
}

// GetAllSubscriberByAlarmId 根据alarmId获取所有关联的subscriber
func GetAllSubscriberByAlarmId(alarmId uint64) []*model.AmsSubscriber {
	db := global.WM_DB.Model(&model.AmsSubscriber{})
	var records []*model.AmsSubscriber
	db.Where("`alarm_id` = ?", alarmId).Find(&records)
	if records == nil {
		global.WM_LOG.Error("[失败]根据alarmId获取所有关联的subscriber", zap.Any("err", "找不到记录"))
	}
	return records
}

// GetSubscriberById 通过id获取subscriber实体
func GetSubscriberById(id uint64) (error, *model.AmsSubscriber) {
	var err error
	var subscriber model.AmsSubscriber
	db := global.WM_DB.Model(&model.AmsSubscriber{})
	err = db.Where("`id` = ?", id).First(&subscriber).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("subscriber不存在")
		}
		return err, nil
	}
	return nil, &subscriber
}
