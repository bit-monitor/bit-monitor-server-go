package service

import (
	"bit.monitor.com/global"
	"bit.monitor.com/model"
	"bit.monitor.com/model/response"
	"bit.monitor.com/model/validation"
	"math"
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

// GetSubscriberNotifyRecord 查询报警记录
func GetSubscriberNotifyRecord(r validation.GetSubscriberNotifyRecord) (err error, data interface{}) {
	limit := r.PageSize
	offset := limit * (r.PageNum - 1)
	db := global.WM_DB.Model(&model.AmsSubscriberNotifyRecord{})
	var totalNum int64
	var entityList []model.AmsSubscriberNotifyRecord
	records := make([]response.GetSubscriberNotifyRecord, 0)

	// 报警记录id
	if r.AlarmRecordId != nil {
		db = db.Where("`alarm_record_id` = ?", r.AlarmRecordId)
	}
	// 报警订阅方id
	if r.SubscriberId != nil {
		db = db.Where("`subscriber_id` = ?", r.SubscriberId)
	}
	// 通知状态，0-失败，1-成功
	if r.State != nil {
		db = db.Where("`state` = ?", r.State)
	}
	// 报警内容，格式为JSON字符串
	if r.Content != "" {
		db = db.Where("`content` like ?", "%"+r.Content+"%")
	}
	// 开始时间、结束时间
	if r.StartTime != "" && r.EndTime != "" {
		db = db.Where("`create_time` BETWEEN ? AND ?", r.StartTime, r.EndTime)
	} else if r.StartTime != "" {
		db = db.Where("`create_time` >= ?", r.StartTime)
	} else if r.EndTime != "" {
		db = db.Where("`create_time` <= ?", r.EndTime)
	}

	err = db.Count(&totalNum).Error
	err = db.Limit(limit).Offset(offset).Find(&entityList).Error
	if err != nil {
		return
	}

	// 设置category
	for _, entity := range entityList {
		err, subscriber := GetSubscriberById(entity.SubscriberId)
		if err != nil {
			break
		}
		record := response.GetSubscriberNotifyRecord{
			Id:         entity.Id,
			State:      entity.State,
			Content:    entity.Content,
			Category:   subscriber.Category,
			CreateTime: entity.CreateTime,
		}
		records = append(records, record)
	}

	data = map[string]interface{}{
		"totalNum":  totalNum,
		"totalPage": math.Ceil(float64(totalNum) / float64(r.PageSize)),
		"pageNum":   r.PageNum,
		"pageSize":  r.PageSize,
		"records":   records,
	}
	return err, data
}

// GetWithRelatedInfoSubscriberNotifyRecord 查询报警记录-带关联信息
func GetWithRelatedInfoSubscriberNotifyRecord(r validation.GetWithRelatedInfoSubscriberNotifyRecord) (err error, data interface{}) {
	limit := r.PageSize
	offset := limit * (r.PageNum - 1)
	db := global.WM_DB.Model(&model.AmsSubscriberNotifyRecord{})
	var totalNum int64
	var entityList []model.AmsSubscriberNotifyRecord
	records := make([]response.GetWithRelatedInfoSubscriberNotifyRecord, 0)

	// 报警记录id
	if r.AlarmRecordId != nil {
		db = db.Where("`alarm_record_id` = ?", r.AlarmRecordId)
	}
	// 报警订阅方id
	if r.SubscriberId != nil {
		db = db.Where("`subscriber_id` = ?", r.SubscriberId)
	}
	// 通知状态，0-失败，1-成功
	if r.State != nil {
		db = db.Where("`state` = ?", r.State)
	}
	// 报警内容，格式为JSON字符串
	if r.Content != "" {
		db = db.Where("`content` like ?", "%"+r.Content+"%")
	}
	// 开始时间、结束时间
	if r.StartTime != "" && r.EndTime != "" {
		db = db.Where("`create_time` BETWEEN ? AND ?", r.StartTime, r.EndTime)
	} else if r.StartTime != "" {
		db = db.Where("`create_time` >= ?", r.StartTime)
	} else if r.EndTime != "" {
		db = db.Where("`create_time` <= ?", r.EndTime)
	}

	err = db.Count(&totalNum).Error
	err = db.Limit(limit).Offset(offset).Find(&entityList).Error
	if err != nil {
		return
	}

	// 设置category、alarmName
	for _, entity := range entityList {
		err, subscriber := GetSubscriberById(entity.SubscriberId)
		if err != nil {
			break
		}
		err, alarm := GetAlarmById(subscriber.AlarmId)
		record := response.GetWithRelatedInfoSubscriberNotifyRecord{
			Id:         entity.Id,
			State:      entity.State,
			Content:    entity.Content,
			Category:   subscriber.Category,
			AlarmName:  alarm.Name,
			CreateTime: entity.CreateTime,
		}
		records = append(records, record)
	}

	data = map[string]interface{}{
		"totalNum":  totalNum,
		"totalPage": math.Ceil(float64(totalNum) / float64(r.PageSize)),
		"pageNum":   r.PageNum,
		"pageSize":  r.PageSize,
		"records":   records,
	}
	return err, data
}
