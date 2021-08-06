package service

import (
	"bit.monitor.com/global"
	"bit.monitor.com/model"
	"bit.monitor.com/model/response"
	"bit.monitor.com/model/validation"
	"go.uber.org/zap"
	"math"
	"time"
)

func AddCustomErrorLog(r validation.AddCustomErrorLog) (err error, data interface{}) {
	db := global.WM_DB.Model(&model.LmsCustomErrorLog{})

	// 保存实体
	log := model.LmsCustomErrorLog{
		LogType:           r.LogType,
		ProjectIdentifier: r.ProjectIdentifier,
		CreateTime:        time.Now(),
		Cuuid:             r.Cuuid,
		Buid:              r.Buid,
		Buname:            r.Buname,
		PageUrl:           r.PageUrl,
		PageKey:           r.PageKey,
		DeviceName:        r.DeviceName,
		Os:                r.Os,
		OsVersion:         r.OsVersion,
		BrowserName:       r.BrowserName,
		BrowserVersion:    r.BrowserVersion,
		IpAddress:         r.IpAddress,
		Address:           r.Address,
		NetType:           r.NetType,
		ErrorType:         r.ErrorType,
		ErrorMessage:      r.ErrorMessage,
		ErrorStack:        r.ErrorStack,
	}
	err = db.Create(&log).Error
	if err != nil {
		global.WM_LOG.Error("保存custom异常日志失败", zap.Any("err", err))
		return err, false
	}
	return nil, true
}

func GetCustomErrorLog(r validation.GetCustomErrorLog) (err error, data interface{}) {
	limit := r.PageSize
	offset := limit * (r.PageNum - 1)
	db := global.WM_DB.Model(&model.LmsCustomErrorLog{})
	var totalNum int64
	var records []model.LmsCustomErrorLog

	// 日志类型
	if r.LogType != "" {
		db = db.Where("`log_type` like ?", "%"+r.LogType+"%")
	}
	// 项目标识
	if r.ProjectIdentifier != "" {
		db = db.Where("`project_identifier` = ?", r.ProjectIdentifier)
	}
	// 用户名
	if r.Buname != "" {
		db = db.Where("`b_uname` like ?", "%"+r.Buname+"%")
	}
	// 页面URL
	if r.PageUrl != "" {
		db = db.Where("`page_url` like ?", "%"+r.PageUrl+"%")
	}
	// JS错误类型
	if r.ErrorType != "" {
		db = db.Where("`error_type` like ?", "%"+r.ErrorType+"%")
	}
	// JS错误信息
	if r.ErrorMessage != "" {
		db = db.Where("`error_message` like ?", "%"+r.ErrorMessage+"%")
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
	err = db.Limit(limit).Offset(offset).Find(&records).Error
	data = map[string]interface{}{
		"totalNum":  totalNum,
		"totalPage": math.Ceil(float64(totalNum) / float64(r.PageSize)),
		"pageNum":   r.PageNum,
		"pageSize":  r.PageSize,
		"records":   records,
	}
	return err, data
}

func GetCustomErrorLogByGroup(r validation.GetCustomErrorLogByGroup) (err error, data interface{}) {
	limit := r.PageSize
	offset := limit * (r.PageNum - 1)
	db := global.WM_DB.Model(&model.LmsCustomErrorLog{})
	var totalNum int64
	var recordsTotal []interface{}
	var records []response.GetCustomErrorLogByGroup

	// 基础查询
	dbCount := global.WM_DB.Model(&model.LmsCustomErrorLog{}).Select("count(id) as count")
	dbData := db.Select("count(id) as count, max(create_time) as latest_record_time, count(distinct c_uuid) as affect_user_count, error_message")

	// 日志类型
	if r.LogType != "" {
		dbCount = dbCount.Where("`log_type` like ?", "%"+r.LogType+"%")
		dbData = dbData.Where("`log_type` like ?", "%"+r.LogType+"%")
	}
	// 项目标识
	if r.ProjectIdentifier != "" {
		dbCount = dbCount.Where("`project_identifier` = ?", r.ProjectIdentifier)
		dbData = dbData.Where("`project_identifier` = ?", r.ProjectIdentifier)
	}
	// 用户名
	if r.Buname != "" {
		dbCount = dbCount.Where("`b_uname` like ?", "%"+r.Buname+"%")
		dbData = dbData.Where("`b_uname` like ?", "%"+r.Buname+"%")
	}
	// 页面URL
	if r.PageUrl != "" {
		dbCount = dbCount.Where("`page_url` like ?", "%"+r.PageUrl+"%")
		dbData = dbData.Where("`page_url` like ?", "%"+r.PageUrl+"%")
	}
	// JS错误类型
	if r.ErrorType != "" {
		dbCount = dbCount.Where("`error_type` like ?", "%"+r.ErrorType+"%")
		dbData = dbData.Where("`error_type` like ?", "%"+r.ErrorType+"%")
	}
	// JS错误信息
	if r.ErrorMessage != "" {
		dbCount = dbCount.Where("`error_message` like ?", "%"+r.ErrorMessage+"%")
		dbData = dbData.Where("`error_message` like ?", "%"+r.ErrorMessage+"%")
	}

	// 开始时间、结束时间
	if r.StartTime != "" && r.EndTime != "" {
		dbCount = dbCount.Where("`create_time` BETWEEN ? AND ?", r.StartTime, r.EndTime)
		dbData = dbData.Where("`create_time` BETWEEN ? AND ?", r.StartTime, r.EndTime)
	} else if r.StartTime != "" {
		dbCount = dbCount.Where("`create_time` >= ?", r.StartTime)
		dbData = dbData.Where("`create_time` >= ?", r.StartTime)
	} else if r.EndTime != "" {
		dbCount = dbCount.Where("`create_time` <= ?", r.EndTime)
		dbData = dbData.Where("`create_time` <= ?", r.EndTime)
	}

	// 分组
	dbCount = dbCount.Group("error_message")
	dbData = dbData.Group("error_message").Order("count desc")

	// 计算总计
	err = dbCount.Find(&recordsTotal).Count(&totalNum).Error
	if err != nil {
		return err, nil
	}
	totalNum = int64(len(recordsTotal))

	// 分组
	err = dbData.Limit(limit).Offset(offset).Find(&records).Error

	data = map[string]interface{}{
		"totalNum":  totalNum,
		"totalPage": math.Ceil(float64(totalNum) / float64(r.PageSize)),
		"pageNum":   r.PageNum,
		"pageSize":  r.PageSize,
		"records":   records,
	}
	return err, data
}

// GetCusCountByIdBetweenStartTimeAndEndTime 查询某个时间段内的日志总数
func GetCusCountByIdBetweenStartTimeAndEndTime(projectIdentifier string, startTime string, endTime string) (error, int64) {
	var err error
	var count int64
	db := global.WM_DB.Model(&model.LmsCustomErrorLog{})
	db = db.Where("`project_identifier` = ?", projectIdentifier)
	db = db.Where("`create_time` BETWEEN ? AND ?", startTime, endTime)
	err = db.Count(&count).Error
	return err, count
}

// GetCusLogCountByHours 按小时间隔获取各小时内的日志数量
func GetCusLogCountByHours(projectIdentifier string, startTime time.Time, endTime time.Time) (error, []response.GetLogCountByHours) {
	var err error
	var results []response.GetLogCountByHours
	db := global.WM_DB.Model(&model.LmsCustomErrorLog{})
	db = db.Select("date_format(create_time, '%Y-%m-%d %H') as hour, count(id) as count")
	db = db.Where("`project_identifier` = ?", projectIdentifier)
	db = db.Where("`create_time` BETWEEN ? AND ?", startTime, endTime)
	db = db.Group("hour")
	err = db.Find(&results).Error
	return err, results
}

// GetCusLogCountByDays 按天间隔获取各天内的日志数量
func GetCusLogCountByDays(projectIdentifier string, startTime time.Time, endTime time.Time) (error, []response.GetLogCountByDays) {
	var err error
	var results []response.GetLogCountByDays
	db := global.WM_DB.Model(&model.LmsCustomErrorLog{})
	db = db.Select("date_format(create_time, '%Y-%m-%d') as day, count(id) as count")
	db = db.Where("`project_identifier` = ?", projectIdentifier)
	db = db.Where("`create_time` BETWEEN ? AND ?", startTime, endTime)
	db = db.Group("day")
	err = db.Find(&results).Error
	return err, results
}

// GetCusLogListByCreateTimeAndProjectIdentifier 获取时间间隔内的简易日志信息
func GetCusLogListByCreateTimeAndProjectIdentifier(projectIdentifier string, startTime time.Time, endTime time.Time) (error, []response.GetLogListByCreateTimeAndProjectIdentifier) {
	var err error
	var results []response.GetLogListByCreateTimeAndProjectIdentifier
	db := global.WM_DB.Model(&model.LmsCustomErrorLog{})
	db = db.Select("id, c_uuid, create_time")
	db = db.Where("`project_identifier` = ?", projectIdentifier)
	db = db.Where("`create_time` BETWEEN ? AND ?", startTime, endTime)
	err = db.Find(&results).Error
	return err, results
}

// GetCusAllLogsBetweenStartTimeAndEndTime 获取时间间隔内的所有日志
func GetCusAllLogsBetweenStartTimeAndEndTime(projectIdentifier string, startTime time.Time, endTime time.Time) (error, []response.GetAllLogsBetweenStartTimeAndEndTime) {
	var err error
	var results []response.GetAllLogsBetweenStartTimeAndEndTime
	db := global.WM_DB.Model(&model.LmsCustomErrorLog{})
	db = db.Select("net_type, device_name, os, browser_name")
	db = db.Where("`project_identifier` = ?", projectIdentifier)
	db = db.Where("`create_time` BETWEEN ? AND ?", startTime, endTime)
	err = db.Find(&results).Error
	return err, results
}

// CountCusDistinctCUuidByCreateTimeBetween 获取时间范围内的影响用户数
func CountCusDistinctCUuidByCreateTimeBetween(startTime time.Time, endTime time.Time) int64 {
	var count int64
	db := global.WM_DB.Model(&model.LmsCustomErrorLog{})
	db.Where("`create_time` BETWEEN ? AND ?", startTime, endTime).Distinct("c_uuid").Count(&count)
	return count
}

// CountCusByCreateTimeBetween 获取时间范围内的日志数
func CountCusByCreateTimeBetween(startTime time.Time, endTime time.Time) int64 {
	var count int64
	db := global.WM_DB.Model(&model.LmsCustomErrorLog{})
	db.Where("`create_time` BETWEEN ? AND ?", startTime, endTime).Count(&count)
	return count
}
