package service

import (
	"bit.monitor.com/global"
	"bit.monitor.com/model"
)

// AddAlarmRecord 新增预警记录
func AddAlarmRecord(alarmRecord *model.AmsAlarmRecord) error {
	db := global.WM_DB.Model(&model.AmsAlarmRecord{})
	return db.Save(alarmRecord).Error
}
