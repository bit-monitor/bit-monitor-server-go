package service

import (
	"bit.monitor.com/global"
	"bit.monitor.com/model"
	"bit.monitor.com/model/response"
	"bit.monitor.com/model/validation"
	"bit.monitor.com/utils"
	"bit.monitor.com/utils/dingtalk"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"math"
	"strconv"
	"strings"
	"time"
)

func AddAlarm(r validation.AddAlarm, userId uint64) (err error, data interface{}) {
	var entity model.AmsAlarm

	// 使用事务的方式提交，避免中途异常后导致错误改动了关联关系的问题
	err = global.WM_DB.Transaction(func(tx *gorm.DB) error {
		db := tx.Model(&model.AmsAlarm{})

		// 保存实体
		entity = model.AmsAlarm{
			Name:              r.Name,
			Level:             *r.Level,
			Category:          *r.Category,
			Rule:              r.Rule,
			StartTime:         r.StartTime,
			EndTime:           r.EndTime,
			SilentPeriod:      *r.SilentPeriod,
			IsActive:          *r.IsActive,
			ProjectIdentifier: r.ProjectIdentifier,
			IsDeleted:         0,
			CreateTime:        time.Now(),
			UpdateTime:        time.Now(),
			CreateBy:          userId,
		}
		err = db.Create(&entity).Error
		if err != nil {
			global.WM_LOG.Error("[失败]新增预警", zap.Any("err", err))
			return err
		}

		// 保存subscriberList
		alarmId := entity.Id
		if len(r.SubscriberList) > 0 {
			for _, subscribe := range r.SubscriberList {
				err = AddSubscriber(tx, subscribe, alarmId)
				if err != nil {
					return err
				}
			}
		}

		// 启动预警定时任务
		if entity.IsActive == 1 {
			err = startAlarmScheduler(tx, &entity)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		global.WM_LOG.Error("[失败]新增预警-事务回滚", zap.Any("err", err))
		return err, nil
	}

	return nil, entity
}

func UpdateAlarm(r validation.UpdateAlarm) (err error, data interface{}) {
	var entity model.AmsAlarm

	// 使用事务的方式提交，避免中途异常后导致错误改动了关联关系的问题
	err = global.WM_DB.Transaction(func(tx *gorm.DB) error {
		db := tx.Model(&model.AmsAlarm{})
		err := db.Where("`id` = ?", r.Id).First(&entity).Error

		// 预警不存在
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err = errors.New("该预警不存在")
			}
			return err
		}

		// 编辑内容
		// 预警名称
		if r.Name != "" {
			entity.Name = r.Name
		}
		// 项目标识
		if r.ProjectIdentifier != "" {
			entity.ProjectIdentifier = r.ProjectIdentifier
		}
		// 报警等级
		if r.Level != nil {
			entity.Level = *r.Level
		}
		// 过滤条件
		if r.Category != nil {
			entity.Category = *r.Category
		}
		// 预警规则
		if r.Rule != "" {
			entity.Rule = r.Rule
		}
		// 报警时段-开始时间
		if r.StartTime != "" {
			entity.StartTime = r.StartTime
		}
		// 报警时段-结束时间
		if r.EndTime != "" {
			entity.EndTime = r.EndTime
		}
		// 静默期
		if r.SilentPeriod != nil {
			entity.SilentPeriod = *r.SilentPeriod
		}
		// 是否启用
		if r.IsActive != nil {
			// 若状态改为启动，则先停止已有的定时任务，再重新启动对应的定时任务
			stopAlarmScheduler(&entity)

			if *r.IsActive == 1 {
				err := startAlarmScheduler(tx, &entity)
				if err != nil {
					return err
				}
			}
			entity.IsActive = *r.IsActive
		}
		// 是否已被删除
		if r.IsDeleted != nil {
			entity.IsDeleted = *r.IsDeleted
		}
		if len(r.SubscriberList) > 0 {
			alarmId := entity.Id
			// 先删除已有的关联关系
			err = DeleteAllSubscriberByAlarmId(tx, alarmId)
			if err != nil {
				return err
			}

			// 再创建新的关联关系
			if len(r.SubscriberList) > 0 {
				for _, subscribe := range r.SubscriberList {
					err = AddSubscriber(tx, subscribe, alarmId)
					if err != nil {
						return err
					}
				}
			}
		}
		entity.UpdateTime = time.Now()
		err = db.Save(&entity).Error
		return err
	})

	if err != nil {
		global.WM_LOG.Error("[失败]编辑预警-事务回滚", zap.Any("err", err))
		return err, nil
	}

	return nil, entity
}

func GetAlarm(r validation.GetAlarm) (err error, data interface{}) {
	limit := r.PageSize
	offset := limit * (r.PageNum - 1)
	db := global.WM_DB.Model(&model.AmsAlarm{})
	var totalNum int64
	var alarmList []model.AmsAlarm
	records := make([]response.GetAlarm, 0)

	// 预警名称
	if r.Name != "" {
		db = db.Where("`name` = ?", r.Name)
	}
	// 项目标识
	if r.ProjectIdentifier != "" {
		db = db.Where("`project_identifier` = ?", r.ProjectIdentifier)
	}
	// 报警等级
	if r.Level != nil {
		db = db.Where("`level` = ?", r.Level)
	}
	// 过滤条件
	if r.Category != nil {
		db = db.Where("`category` = ?", r.Category)
	}
	// 预警规则
	if r.Rule != "" {
		db = db.Where("`rule` like ?", "%"+r.Rule+"%")
	}
	// 报警时段-开始时间
	if r.StartTime != "" {
		db = db.Where("`start_time` = ?", r.StartTime)
	}
	// 报警时段-结束时间
	if r.EndTime != "" {
		db = db.Where("`end_time` = ?", r.EndTime)
	}
	// 静默期
	if r.SilentPeriod != nil {
		db = db.Where("`silent_period` = ?", r.SilentPeriod)
	}
	// 是否启用
	if r.IsActive != nil {
		db = db.Where("`is_active` = ?", r.IsActive)
	}
	// 创建人ID
	if r.CreateBy != nil {
		db = db.Where("`create_by` = ?", r.CreateBy)
	}
	// 是否已被删除
	if r.IsDeleted != nil {
		db = db.Where("`is_deleted` = ?", r.IsDeleted)
	}

	err = db.Count(&totalNum).Error
	err = db.Limit(limit).Offset(offset).Find(&alarmList).Error
	if err != nil {
		return
	}

	// 设置subscriberList
	for _, alarm := range alarmList {
		record := response.GetAlarm{
			Id:                alarm.Id,
			Name:              alarm.Name,
			ProjectIdentifier: alarm.ProjectIdentifier,
			Level:             alarm.Level,
			Category:          alarm.Category,
			Rule:              alarm.Rule,
			StartTime:         alarm.StartTime,
			EndTime:           alarm.EndTime,
			SilentPeriod:      alarm.SilentPeriod,
			IsActive:          alarm.IsActive,
			CreateBy:          alarm.CreateBy,
			IsDeleted:         alarm.IsDeleted,
			SubscriberList:    "",
		}
		err = setSubscriberList(&record)
		if err != nil {
			break
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

func DeleteAlarm(alarmId uint64) (err error, data interface{}) {
	var entity model.AmsAlarm

	// 使用事务的方式提交，避免中途异常后导致错误改动了关联关系的问题
	err = global.WM_DB.Transaction(func(tx *gorm.DB) error {

		db := tx.Model(&model.AmsAlarm{})
		err = db.Where("`id` = ?", alarmId).First(&entity).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("预警不存在")
		}

		// 删除关联的订阅者
		err = DeleteAllSubscriberByAlarmId(tx, alarmId)
		if err != nil {
			return err
		}

		// 停止预警定时任务
		stopAlarmScheduler(&entity)

		// 删除实体记录
		err = db.Delete(&entity).Error
		return err
	})

	if err != nil {
		global.WM_LOG.Error("[失败]删除预警-事务回滚", zap.Any("err", err))
		return err, nil
	}

	return nil, true
}

func GetAlarmRecord(r validation.GetAlarmRecord) (err error, data interface{}) {
	limit := r.PageSize
	offset := limit * (r.PageNum - 1)
	db := global.WM_DB.Model(&model.AmsAlarmRecord{})
	var totalNum int64
	var alarmRecordList []model.AmsAlarmRecord
	records := make([]response.GetAlarmRecord, 0)

	// 预警id
	if r.AlarmId != nil {
		db = db.Where("`alarm_id` = ?", r.AlarmId)
	}
	// 报警内容，格式为JSON字符串
	if r.AlarmData != "" {
		db = db.Where("`alarm_data` like ?", "%"+r.AlarmData+"%")
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
	err = db.Limit(limit).Offset(offset).Order("create_time desc").Find(&alarmRecordList).Error
	if err != nil {
		return
	}

	// 设置alarmName
	for _, alarmRecord := range alarmRecordList {
		err, alarm := GetAlarmById(alarmRecord.AlarmId)
		if err != nil {
			break
		}
		record := response.GetAlarmRecord{
			Id:         alarmRecord.Id,
			AlarmId:    alarmRecord.AlarmId,
			AlarmData:  alarmRecord.AlarmData,
			AlarmName:  alarm.Name,
			CreateTime: alarmRecord.CreateTime,
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

// GetProjectNameByAlarmId 根据alarmId获取关联的项目名称
func GetProjectNameByAlarmId(alarmId uint64) (error, string) {
	db := global.WM_DB.Model(&model.AmsAlarm{})
	var err error
	var alarm model.AmsAlarm
	err = db.Where("`id` = ?", alarmId).First(&alarm).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("找不到预警信息")
		}
		return err, ""
	}
	err, projectName := GetProjectNameByProjectIdentifier(alarm.ProjectIdentifier)
	if err != nil {
		return err, ""
	}
	return nil, projectName
}

// 设置subscriberList
func setSubscriberList(a *response.GetAlarm) error {
	var err error
	alarmId := a.Id
	subscriberEntityList := GetAllSubscriberByAlarmId(alarmId)
	subscriberList, err := json.Marshal(subscriberEntityList)
	a.SubscriberList = string(subscriberList)
	return err
}

// GetAlarmById 根据id获取alarm实体
func GetAlarmById(id uint64) (error, *model.AmsAlarm) {
	var err error
	var alarm model.AmsAlarm
	db := global.WM_DB.Model(&model.AmsAlarm{})
	err = db.Where("`id` = ?", id).First(&alarm).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("预警不存在")
		return err, nil
	}
	return nil, &alarm
}

// 启动预警定时任务
func startAlarmScheduler(tx *gorm.DB, entity *model.AmsAlarm) error {
	params, err := utils.StructToJson(entity)
	if err != nil {
		global.WM_LOG.Error("[失败]启动预警定时任务", zap.Any("err", err))
		return err
	}

	// 保存定时任务
	err, scheduler := AddScheduler(tx, params)
	if err != nil {
		return err
	}

	// 保存预警-定时任务关联表
	alarmId := entity.Id
	schedulerId := scheduler.Id
	err = AddAlarmSchedulerRelation(tx, alarmId, schedulerId)
	if err != nil {
		return err
	}

	// 创建定时任务并启动
	err = StartScheduler(scheduler)

	return err
}

// 停止预警定时任务
func stopAlarmScheduler(alarm *model.AmsAlarm) {

	// 若预警未启动，则跳过
	if alarm.IsActive == 0 {
		return
	}

	// 若为启动中的预警
	err, relationList := GetAllAlarmSchedulerRelationByAlarmId(alarm.Id)
	if err != nil {
		global.WM_LOG.Error("[失败]停止预警定时任务", zap.Any("err", err))
		return
	}
	if len(relationList) == 0 {
		global.WM_LOG.Info("[失败]停止预警定时任务", zap.String("relationList", "没有找到关联的信息"))
		return
	}

	for _, asr := range relationList {
		schedulerId := asr.SchedulerId

		// 终止其关联的执行中的定时任务
		err = utils.StopAndDeleteBySchedulerId(schedulerId)
		if err != nil {
			global.WM_LOG.Error("[失败]停止预警定时任务", zap.Any("StopAndDeleteBySchedulerId", err))
			continue
		}

		// 删除定时任务
		err = DeleteSchedulerById(schedulerId)
		if err != nil {
			global.WM_LOG.Error("[失败]停止预警定时任务", zap.Any("DeleteSchedulerById", err))
			return
		}

		// 删除预警-定时任务关联表
		err = DeleteAlarmSchedulerRelationById(asr.Id)
		if err != nil {
			global.WM_LOG.Error("[失败]停止预警定时任务", zap.Any("DeleteAlarmSchedulerRelationById", err))
			return
		}
	}
}

// startAlarmSchedule 预警定时任务执行的内容
func startAlarmSchedule(params string) error {
	var entity model.AmsAlarm
	err := json.Unmarshal([]byte(params), &entity)
	if err != nil {
		global.WM_LOG.Error("[失败]定时任务执行-解析params参数异常", zap.Any("err", err))
		return err
	}

	// 获取预警规则
	var rule validation.SchedulerRule
	err = json.Unmarshal([]byte(entity.Rule), &rule)
	if err != nil {
		return err
	}

	// 过滤条件
	category := entity.Category
	tableNameMap := map[int8]string{
		1: "lms_js_error_log",
		2: "lms_http_error_log",
		3: "lms_resource_load_error_log",
		4: "lms_custom_error_log",
	}

	resultList := make([]*validation.AlarmScheduleResult, 0)

	// 根据category过滤条件，查询对应的日志表
	if category == 0 {
		tempList := make([]*validation.AlarmScheduleResult, 0)
		for _, tableName := range tableNameMap {
			for _, ruleItem := range rule.Rules {
				tempList = append(tempList, CheckIsExceedAlarmThreshold(tableName, &ruleItem))
			}
		}
		// 聚合分析
		setResultListByTempList(&tempList, &resultList)
	} else {
		tableName := tableNameMap[category]
		for _, ruleItem := range rule.Rules {
			resultList = append(resultList, CheckIsExceedAlarmThreshold(tableName, &ruleItem))
		}
	}

	// 根据条件规则判断是否触发报警
	op := rule.Op
	isAnd := true
	isOr := false
	for _, resultItem := range resultList {
		isExceedAlarmThreshold := resultItem.IsExceedAlarmThreshold
		isAnd = isAnd && isExceedAlarmThreshold
		isOr = isOr || isExceedAlarmThreshold
		if isExceedAlarmThreshold {
			// 整合报警内容
			//thresholdValue := strconv.FormatFloat(resultItem.ThresholdValue, 'f', 2, 64)
			//actualValue := strconv.FormatFloat(resultItem.ActualValue, 'f', 2, 64)
			//alarmContent := "[预警指标]" + string(resultItem.TargetInd) + ", " + "[预警阈值]" + thresholdValue + ", " + "[实际值]" + actualValue
		}
	}
	if (op == "&&" && isAnd) || (op == "||" && isOr) {
		// 触发预警条件，添加报警记录
		global.WM_LOG.Info("[信息]预警定时任务", zap.Any("info", fmt.Sprintf("预警名称：%v，报警内容：%v", entity.Name, resultList)))
		err := saveAlarmRecordAndNotifyAllSubscribers(&entity, resultList)
		if err != nil {
			return err
		}
	}

	return nil
}

// 当category为0时，即选择的过滤条件为全部，此时需要对各个表的计算结果进行聚合分析
func setResultListByTempList(tempList *[]*validation.AlarmScheduleResult, resultList *[]*validation.AlarmScheduleResult) {
	for _, tempItem := range *tempList {
		var resultItem *validation.AlarmScheduleResult
		for _, tempRItem := range *resultList {
			if tempRItem.TargetInd == tempItem.TargetInd {
				resultItem = tempRItem
			}
		}
		if resultItem == nil {
			*resultList = append(*resultList, tempItem)
		} else {
			thresholdValue := resultItem.ThresholdValue
			oldActualValue := resultItem.ActualValue
			newActualValue := oldActualValue + tempItem.ActualValue
			resultItem.ActualValue = newActualValue
			resultItem.IsExceedAlarmThreshold = newActualValue > thresholdValue
		}
	}
}

// 保存报警记录，同时通知所有报警订阅方
func saveAlarmRecordAndNotifyAllSubscribers(alarm *model.AmsAlarm, resultList []*validation.AlarmScheduleResult) error {
	var err error
	subscriberEntityList := GetAllSubscriberByAlarmId(alarm.Id)

	// 保存报警记录
	alarmData, err := json.Marshal(resultList)
	if err != nil {
		return err
	}
	alarmRecord := model.AmsAlarmRecord{
		AlarmId:    alarm.Id,
		AlarmData:  string(alarmData),
		CreateTime: time.Now(),
	}
	err = AddAlarmRecord(&alarmRecord)
	if err != nil {
		return err
	}

	// 通知所有报警订阅方
	if len(resultList) > 0 && len(subscriberEntityList) > 0 {
		notifyAllSubscribers(alarm.Id, alarmRecord.AlarmId, resultList, subscriberEntityList)
	}

	return nil
}

// 通知所有报警订阅方
func notifyAllSubscribers(alarmId uint64, alarmRecordId uint64, resultList []*validation.AlarmScheduleResult, subscriberEntityList []*model.AmsSubscriber) {
	for _, subscriber := range subscriberEntityList {
		if subscriber.IsActive == 1 {
			subscriberNotifyRecord := model.AmsSubscriberNotifyRecord{
				AlarmRecordId: alarmRecordId,
				SubscriberId:  subscriber.Id,
			}

			// 设置通知内容
			content := strings.Builder{}
			_, projectName := GetProjectNameByAlarmId(alarmId)
			alarmTime := utils.GetCurrentTimeByDefaultLayout()
			content.WriteString("项目名：")
			content.WriteString(projectName)
			content.WriteString("\n")
			content.WriteString("报警时间：")
			content.WriteString(alarmTime)
			content.WriteString("\n")
			content.WriteString("报警内容：")
			content.WriteString("\n")
			for index, resultItem := range resultList {
				thresholdValue := strconv.FormatFloat(resultItem.ThresholdValue, 'f', 2, 64)
				actualValue := strconv.FormatFloat(resultItem.ActualValue, 'f', 2, 64)
				startTimeStr := utils.PassTimeToStrByDefaultLayout(resultItem.StartTime)

				// 简化endTimeStr，如果与startTime年月日相同，则省略
				preStart := utils.PassTimeToStrByLayout(resultItem.StartTime, "2006-01-02")
				preEnd := utils.PassTimeToStrByLayout(resultItem.EndTime, "2006-01-02")
				var endTimePattern string
				if preStart == preEnd {
					endTimePattern = "15:04:05"
				} else {
					endTimePattern = utils.DefaultLayout
				}
				endTimeStr := utils.PassTimeToStrByLayout(resultItem.EndTime, endTimePattern)

				content.WriteString(string(rune(index + 1)))
				content.WriteString(".")
				content.WriteString("[")
				content.WriteString(resultItem.TargetInd)
				content.WriteString("]")
				content.WriteString("区间：")
				content.WriteString(startTimeStr)
				content.WriteString("至")
				content.WriteString(endTimeStr)
				content.WriteString("，")
				content.WriteString("阈值：")
				content.WriteString(thresholdValue)
				content.WriteString("，")
				content.WriteString(actualValue)
				if index < len(resultList) {
					content.WriteString("\n")
				}
			}
			subscriberNotifyRecord.Content = content.String()

			if subscriber.Category == 1 {
				// 钉钉机器人
				// 若为钉钉机器人，则subscriberList为access_token列表，多个以逗号隔开
				subscriberList := strings.Split(subscriber.Subscriber, ",")
				for _, accessToken := range subscriberList {
					if accessToken != "" {

						// config.yaml配置文件中,获取钉钉推送关键词
						keyword := global.WM_CONFIG.Dingtalk.Keyword

						// 从application配置文件中，获取是否开启钉钉推送
						isEnableDingTalk := global.WM_CONFIG.Dingtalk.Enable
						if isEnableDingTalk {
							config := &dingtalk.Config{
								AccessToken: accessToken,
								KeyWord:     keyword,
							}
							err, client := dingtalk.NewClient(config)
							if err != nil {
								global.WM_LOG.Error("[失败]发送钉钉机器人通知", zap.Any("err", err))
							}
							params := dingtalk.GetTextParams()
							params.Text = dingtalk.TextContent{
								Content: keyword + "\n" + content.String(),
							}
							err = client.SendText(params)
							if err != nil {
								global.WM_LOG.Error("发送钉钉机器人通知失败", zap.Any("err", err))
								subscriberNotifyRecord.State = 0
							} else {
								subscriberNotifyRecord.State = 1
							}
						} else {
							subscriberNotifyRecord.State = 0
						}
						_ = AddSubscriberNotifyRecord(&subscriberNotifyRecord)
					}
				}
			} else if subscriber.Category == 2 {
				// TODO 发送邮件通知
				_ = AddSubscriberNotifyRecord(&subscriberNotifyRecord)
			}
		}
	}
}
