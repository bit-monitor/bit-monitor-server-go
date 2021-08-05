package service

import (
	"bit.monitor.com/global"
	"bit.monitor.com/model"
	"time"
)

// AddSubscriberNotifyRecord 新增
func AddSubscriberNotifyRecord(entity *model.AmsSubscriberNotifyRecord) error {
	var err error
	db := global.WM_DB.Model(&model.AmsSubscriberNotifyRecord{})
	entity.CreateTime = time.Now()
	err = db.Save(entity).Error
	return err
}
