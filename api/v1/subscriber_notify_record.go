package v1

import (
	"bit.monitor.com/global"
	"bit.monitor.com/model/response"
	"bit.monitor.com/model/validation"
	"bit.monitor.com/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetSubscriberNotifyRecord(c *gin.Context) {
	var err error
	var r validation.GetSubscriberNotifyRecord
	err = c.ShouldBind(&r)
	if err != nil {
		global.WM_LOG.Error("[失败]查询报警记录", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	if err, data := service.GetSubscriberNotifyRecord(r); err != nil {
		global.WM_LOG.Error("[失败]查询报警记录", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		// global.WM_LOG.Info("[成功]查询报警记录", zap.Any("data", data))
		response.SuccessWithData(data, c)
	}
}

func GetWithRelatedInfoSubscriberNotifyRecord(c *gin.Context) {
	var err error
	var r validation.GetWithRelatedInfoSubscriberNotifyRecord
	err = c.ShouldBind(&r)
	if err != nil {
		global.WM_LOG.Error("[失败]查询报警记录-带关联信息", zap.Any("err", err))
		response.FailWithError(err, c)
		return
	}
	if err, data := service.GetWithRelatedInfoSubscriberNotifyRecord(r); err != nil {
		global.WM_LOG.Error("[失败]查询报警记录-带关联信息", zap.Any("err", err))
		response.FailWithError(err, c)
	} else {
		// global.WM_LOG.Info("[成功]查询报警记录-带关联信息", zap.Any("data", data))
		response.SuccessWithData(data, c)
	}
}
