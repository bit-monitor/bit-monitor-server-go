package v1

import (
	"bit.monitor.com/global"
	"bit.monitor.com/model/response"
	"bit.monitor.com/model/validation"
	"bit.monitor.com/service"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

func AddAlarm(c *gin.Context) {
	var err error
	var r validation.AddAlarm
	err = c.ShouldBind(&r)
	if err != nil {
		global.WM_LOG.Error("[失败]新增预警", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}

	// 因gin暂时还不支持从content-type: application/x-www-form-urlencoded, 或content-type: multipart/form-data
	// 中解析出array等复杂结构，因此这里暂时改为单独解构
	subscriberListParam := c.PostForm("subscriberList")
	if subscriberListParam == "" {
		err = errors.New("subscriberList不能为空")
		global.WM_LOG.Error("[失败]新增预警", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	var subscriberList []validation.Subscriber
	err = json.Unmarshal([]byte(subscriberListParam), &subscriberList)
	if err != nil {
		global.WM_LOG.Error("[失败]新增预警", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}

	err, userId := service.GetUserIdByContext(c)
	if err != nil {
		global.WM_LOG.Error("[失败]新增预警", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	if err, data := service.AddAlarm(r, userId); err != nil {
		global.WM_LOG.Error("[失败]新增预警", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		// global.WM_LOG.Info("[成功]新增预警", zap.Any("data", data))
		response.SuccessWithData(data, c)
	}
}

func UpdateAlarm(c *gin.Context) {
	var err error
	var r validation.UpdateAlarm
	err = c.ShouldBind(&r)
	if err != nil {
		global.WM_LOG.Error("[失败]编辑预警", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}

	// 因gin暂时还不支持从content-type: application/x-www-form-urlencoded, 或content-type: multipart/form-data
	// 中解析出array等复杂结构，因此这里暂时改为单独解构
	subscriberListParam := c.PostForm("subscriberList")
	if subscriberListParam != "" {
		var subscriberList []validation.Subscriber
		err = json.Unmarshal([]byte(subscriberListParam), &subscriberList)
		if err != nil {
			global.WM_LOG.Error("[失败]编辑预警", zap.Any("err", err))
			response.FailWithError(err, c)
			return
		}
		r.SubscriberList = subscriberList
	}

	if err, data := service.UpdateAlarm(r); err != nil {
		global.WM_LOG.Error("[失败]编辑预警", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		// global.WM_LOG.Info("[成功]编辑预警", zap.Any("data", data))
		response.SuccessWithData(data, c)
	}
}

func GetAlarm(c *gin.Context) {
	var err error
	var r validation.GetAlarm
	err = c.ShouldBind(&r)
	if err != nil {
		global.WM_LOG.Error("[失败]查询预警", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	if err, data := service.GetAlarm(r); err != nil {
		global.WM_LOG.Error("[失败]查询预警", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		// global.WM_LOG.Info("[成功]查询预警", zap.Any("data", data))
		response.SuccessWithData(data, c)
	}
}

func DeleteAlarm(c *gin.Context) {
	var err error
	id := c.Param("id")
	if id == "" {
		global.WM_LOG.Error("[失败]删除预警，id不能为空", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	alarmId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		global.WM_LOG.Error("[失败]删除预警", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	err, data := service.DeleteAlarm(alarmId)
	if err != nil {
		global.WM_LOG.Error("[失败]删除预警", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		response.SuccessWithData(data, c)
	}
}

func GetAlarmRecord(c *gin.Context) {
	var err error
	var r validation.GetAlarmRecord
	err = c.ShouldBind(&r)
	if err != nil {
		global.WM_LOG.Error("[失败]查询预警记录", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	if err, data := service.GetAlarmRecord(r); err != nil {
		global.WM_LOG.Error("[失败]查询预警记录", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		// global.WM_LOG.Info("[成功]查询预警记录", zap.Any("data", data))
		response.SuccessWithData(data, c)
	}
}
