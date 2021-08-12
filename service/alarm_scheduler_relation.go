package service

import (
	"bit.monitor.com/global"
	"bit.monitor.com/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

func AddAlarmSchedulerRelation(tx *gorm.DB, alarmId uint64, schedulerId uint64) error {
	db := tx.Model(&model.AmsAlarmSchedulerRelation{})
	entity := model.AmsAlarmSchedulerRelation{
		AlarmId:     alarmId,
		SchedulerId: schedulerId,
		CreateTime:  time.Now(),
	}
	err := db.Save(&entity).Error
	if err != nil {
		global.WM_LOG.Error("新增预警计划关联关系失败", zap.Any("err", err))
		return err
	}
	return nil
}

// GetAllAlarmSchedulerRelationByAlarmId 根据预警id获取所有关联实体
func GetAllAlarmSchedulerRelationByAlarmId(alarmId uint64) (error, []*model.AmsAlarmSchedulerRelation) {
	var err error
	var list []*model.AmsAlarmSchedulerRelation
	db := global.WM_DB.Model(&model.AmsAlarmSchedulerRelation{})
	err = db.Where("`alarm_id` = ?", alarmId).Find(&list).Error
	if err != nil {
		return err, nil
	}
	return nil, list
}

// DeleteAlarmSchedulerRelationById 根据id删除AlarmSchedulerRelation实体
func DeleteAlarmSchedulerRelationById(id uint64) error {
	err := global.WM_DB.Where("`id` = ?", id).Delete(&model.AmsAlarmSchedulerRelation{}).Error
	return err
}
